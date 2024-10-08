// Copyright (C) 2022, Rob Lyon <rob@ctxswitch.com>
//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package strata

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type SummaryOpts struct {
	// Objectives defines the quantile rank estimates with their respective
	// absolute error.
	Objectives map[float64]float64
	// MaxAge defines the duration for which an observation stays relevant
	// for the summary.
	MaxAge time.Duration
	// AgeBuckets is the number of buckets used to exclude observations that
	// are older than MaxAge from the summary.
	AgeBuckets uint32
}

// MetricsOpts defines options that are available for the metrics wrapper.
type MetricsOpts struct {
	// ConstantLabels is an array of label/value pairs that will be constant
	// across all metrics.
	ConstantLabels []string
	// HistogramBuckets are buckets used for histogram observation counts.
	HistogramBuckets []float64
	// SummaryOpts defines the options available to summary collectors.
	SummaryOpts *SummaryOpts
	// Registry is the prometheus registry that will be used to register
	// collectors.
	Registry *prometheus.Registry
	// Separator is the separator that will be used to join the metric name
	// components.
	Separator rune
	// PanicOnError maintains the default behavior of prometheus to panic on
	// errors. If this value is set to false, the library attempts to recover
	// from any panics and emits an internally managed metric
	// strata_errors_panic_recovery to inform the operator that visibility is
	// degraded. If set to true the original behavior is maintained and all
	// errors are treated as panics.
	PanicOnError bool
	// Prefix is an array of prefixes that will be appended to the metric name.
	Prefix []string
	// Logger takes a value that matches the Logger interface and is used for
	// log output of errors and other debug information.
	Logger Logger
}

// Metrics provides a wrapper around the prometheus client to automatically
// register and collect metrics.
type Metrics struct {
	separator        rune
	prefix           string
	histogramBuckets []float64
	summaryOpts      *SummaryOpts
	store            *Store
	labels           []string
	errors           *ApexInternalErrorMetrics
	panicOnError     bool
	registry         *prometheus.Registry
	registerer       prometheus.Registerer
	server           *Server
	logger           Logger
}

// New creates a new Apex metrics store using the options that have
// been provided.
func New(opts MetricsOpts) *Metrics {
	opts = defaultedMetrics(opts)
	prefix := strings.Join(opts.Prefix, string(opts.Separator))
	labels := SlicePairsToMap(opts.ConstantLabels)

	_ = opts.Registry.Register(collectors.NewGoCollector())
	_ = opts.Registry.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return &Metrics{
		prefix:           prefix,
		separator:        opts.Separator,
		histogramBuckets: opts.HistogramBuckets,
		summaryOpts:      opts.SummaryOpts,
		store:            newStore(),
		labels:           []string{},
		panicOnError:     opts.PanicOnError,
		errors:           NewApexInternalErrorMetrics(opts.Prefix, opts.Separator),
		registry:         opts.Registry,
		registerer:       prometheus.WrapRegistererWith(prometheus.Labels(labels), opts.Registry),
		logger:           opts.Logger,
	}
}

// Start starts the HTTP server.  It blocks until Stop is called.
func (m *Metrics) Start(ctx context.Context, opts ServerOpts) error {
	m.server = newServer(opts).WithLogger(m.logger)
	err := m.server.Start(ctx, m.registry)
	if !errors.Is(err, http.ErrServerClosed) {
		m.logger.Error(err, "prometheus collector endpoint error")
	}

	return nil
}

// Stop shuts down the HTTP server gracefully.
func (m *Metrics) Stop() {
	// Stop has been deprecated in favor of using the callers context.
}

// WithPrefix appends additional values to the metric name to prefix any new
// metric names that are added. By default metrics are created without prefixes
// unless added in MetricOpts. For example:
//
//	m := strata.New(strata.MetricsOpts{})
//	// prefix: ""
//	m.WithPrefix("strata", "example")
//	// prefix: "strata_example"
//	m.CounterInc("a_total")
//	// metric: "strata_example_a_total"
//	n := m.WithPrefix("component")
//	// prefix: "strata_example_component"
//	n.CounterInc("b_total")
//	// metric: "strata_example_component_b_total"
//	m.CounterInc("c_total")
//	// metric: "strata_example_c_total"
func (m *Metrics) WithPrefix(prefix ...string) *Metrics {
	p := strings.Join(prefix, string(m.separator))
	newPrefix := prefixedName(m.prefix, p, m.separator)
	metrics := m.clone()
	metrics.prefix = newPrefix
	// Labels are reset when a new prefix is added
	metrics.labels = []string{}
	return metrics
}

// WithLabels creates a new metric with the provided labels.  Example:
//
//	metrics = metrics.WithValues("label1", "label2")
//	metrics.GaugeAdd("gauge_with_values", 2.0, "value1", "value2")
func (m *Metrics) WithLabels(labels ...string) *Metrics {
	metrics := m.clone()
	metrics.labels = labels
	return metrics
}

// CounterInc increments a counter by 1.
func (m *Metrics) CounterInc(name string, lv ...string) {
	defer m.recover(name, "counter_inc")
	vec, err := m.store.getCounter(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "counter_inc")
		return
	}
	vec.Inc(lv...)
}

// CounterAdd increments a counter by the provided value.
func (m *Metrics) CounterAdd(name string, v float64, lv ...string) {
	defer m.recover(name, "counter_add")
	vec, err := m.store.getCounter(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "counter_add")
		return
	}
	vec.Add(v, lv...)
}

// GaugeSet sets a gauge to an arbitrary value.
func (m *Metrics) GaugeSet(name string, v float64, lv ...string) {
	defer m.recover(name, "gauge_set")
	vec, err := m.store.getGauge(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "gauge_set")
		return
	}
	vec.Set(v, lv...)
}

// GaugeInc increments a gauge by 1.
func (m *Metrics) GaugeInc(name string, lv ...string) {
	defer m.recover(name, "gauge_inc")
	vec, err := m.store.getGauge(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "gauge_inc")
		return
	}
	vec.Inc(lv...)
}

// GaugeDec decrements a gauge by 1.
func (m *Metrics) GaugeDec(name string, lv ...string) {
	defer m.recover(name, "gauge_dec")
	vec, err := m.store.getGauge(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "gauge_dec")
		return
	}
	vec.Dec(lv...)
}

// GaugeAdd adds an arbitrary value to the gauge.
func (m *Metrics) GaugeAdd(name string, v float64, lv ...string) {
	defer m.recover(name, "gauge_add")
	vec, err := m.store.getGauge(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "gauge_add")
		return
	}
	vec.Add(v, lv...)
}

// GaugeSub subtracts an arbitrary value to the gauge.
func (m *Metrics) GaugeSub(name string, v float64, lv ...string) {
	defer m.recover(name, "gauge_sub")
	vec, err := m.store.getGauge(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "gauge_sub")
		return
	}
	vec.Sub(v, lv...)
}

// SummaryObserve adds a single observation to the summary.
func (m *Metrics) SummaryObserve(name string, v float64, lv ...string) {
	defer m.recover(name, "summary_observe")
	vec, err := m.store.getSummary(m.registerer, prefixedName(m.prefix, name, m.separator), *m.summaryOpts, m.labels...)
	if err != nil {
		m.emitError(err, name, "summary_timer")
		return
	}
	vec.Observe(v, lv...)
}

// SummaryTimer returns a Timer helper to measure duration.  ObserveDuration is
// used to measure the time. Example:
//
//	timer := m.SummaryTimer("response")
//	defer timer.ObserveDuration()
func (m *Metrics) SummaryTimer(name string, lv ...string) *Timer {
	defer m.recover(name, "summary_timer")
	vec, err := m.store.getSummary(m.registerer, prefixedName(m.prefix, name, m.separator), *m.summaryOpts, m.labels...)
	if err != nil {
		m.emitError(err, name, "summary_timer")
		// TODO: this is dangerous, fix me
		return nil
	}
	return vec.Timer(lv...)
}

// HistogramObserve adds a single observation to the histogram.
func (m *Metrics) HistogramObserve(name string, v float64, lv ...string) {
	defer m.recover(name, "histogram_observe")
	vec, err := m.store.getHistogram(m.registerer, prefixedName(m.prefix, name, m.separator), m.histogramBuckets, m.labels...)
	if err != nil {
		m.emitError(err, name, "histogram_observe")
		return
	}
	vec.Observe(v, lv...)
}

// HistogramTimer returns a Timer helper to measure duration.  ObserveDuration is
// used to measure the time. Example:
//
//	timer := m.HistogramTimer("response")
//	defer timer.ObserveDuration()
func (m *Metrics) HistogramTimer(name string, lv ...string) *Timer {
	defer m.recover(name, "histogram_timer")
	vec, err := m.store.getHistogram(m.registerer, prefixedName(m.prefix, name, m.separator), m.histogramBuckets, m.labels...)
	if err != nil {
		m.emitError(err, name, "histogram_timer")
		// TODO: this is dangerous, fix me
		return nil
	}
	return vec.Timer(lv...)
}

func (m *Metrics) clone() *Metrics {
	n := *m
	return &n
}

func (m *Metrics) emitError(err error, name string, fn string) {
	if m.panicOnError {
		panic(err)
	}

	switch err {
	case ErrInvalidMetricName:
		m.errors.InvalidMetricName(name, fn)
	case ErrRegistrationFailed:
		m.errors.RegistrationFailed(name, fn)
	case ErrAlreadyRegistered:
		m.errors.AlreadyRegistered(name, fn)
	}
}

func defaultedMetrics(opts MetricsOpts) MetricsOpts {
	// opts.Prefix default is empty
	// opts.MustRegister default is false
	// opts.PanicOnError default is false
	if opts.Registry == nil {
		opts.Registry = prometheus.NewRegistry()
	}

	if opts.HistogramBuckets == nil {
		opts.HistogramBuckets = DefBuckets
	}

	if opts.Separator == 0 {
		opts.Separator = '_'
	}

	if opts.Logger == nil {
		// nil logger that discards all logs.
		opts.Logger = logr.New(nil)
	}

	opts.SummaryOpts = defaultedSummaryOpts(opts.SummaryOpts)

	return opts
}

func defaultedSummaryOpts(opts *SummaryOpts) *SummaryOpts {
	if opts == nil {
		opts = &SummaryOpts{}
	}

	if opts.AgeBuckets < 1 {
		opts.AgeBuckets = DefaultAgeBuckets
	}

	if opts.MaxAge < 1 {
		opts.MaxAge = DefaultMaxAge
	}

	if opts.Objectives == nil {
		opts.Objectives = DefaultObjectives
	}

	return opts
}

func (m *Metrics) recover(name string, fn string) {
	if !m.panicOnError {
		if r := recover(); r != nil {
			var err error
			switch e := r.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("unknown error")
			}

			m.logger.Error(err, "panic recovery", "name", name, "func", fn)
			m.errors.PanicRecovery(name, fn)
		}
	}
}
