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

package apex

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	DefaultTimeout                  = 5 * time.Second
	DefaultMaxAge     time.Duration = 10 * time.Minute
	DefaultAgeBuckets uint32        = 5
)

var (
	DefaultObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
)

type SummaryOpts struct {
	Objectives map[float64]float64
	MaxAge     time.Duration
	AgeBuckets uint32
}

type MetricsOpts struct {
	ConstantLabels   []string
	HistogramBuckets []float64
	SummaryOpts      *SummaryOpts
	Registry         *prometheus.Registry
	Separator        rune
	BindAddr         string
	Path             string
	Port             int
	PanicOnError     bool
	Prefix           []string
}

type Metrics struct {
	separator        rune
	prefix           string
	path             string
	port             int
	bindAddr         string
	histogramBuckets []float64
	summaryOpts      *SummaryOpts
	store            *Store
	labels           []string
	errors           *ApexInternalErrorMetrics
	panicOnError     bool
	registry         *prometheus.Registry
	registerer       prometheus.Registerer
}

func New(opts MetricsOpts) *Metrics {
	opts = defaulted(opts)
	prefix := strings.Join(opts.Prefix, string(opts.Separator))
	labels := SlicePairsToMap(opts.ConstantLabels)
	return &Metrics{
		prefix:           prefix,
		separator:        opts.Separator,
		port:             opts.Port,
		path:             opts.Path,
		bindAddr:         opts.BindAddr,
		histogramBuckets: opts.HistogramBuckets,
		summaryOpts:      opts.SummaryOpts,
		store:            newStore(),
		labels:           []string{},
		panicOnError:     opts.PanicOnError,
		errors:           NewApexInternalErrorMetrics(opts.Prefix, opts.Separator),
		registry:         opts.Registry,
		registerer:       prometheus.WrapRegistererWith(prometheus.Labels(labels), opts.Registry),
	}
}

func (m *Metrics) Start() error {
	mux := http.NewServeMux()
	mux.Handle(m.path, promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{
		Timeout: DefaultTimeout,
	}))
	addr := fmt.Sprintf("%s:%d", m.bindAddr, m.port)
	return http.ListenAndServe(addr, mux)
}

func (m *Metrics) WithPrefix(prefix ...string) *Metrics {
	p := strings.Join(prefix, string(m.separator))
	newPrefix := prefixedName(m.prefix, p, m.separator)
	metrics := m.copy()
	metrics.prefix = newPrefix
	// Labels are reset when a new prefix is added
	metrics.labels = []string{}
	return metrics
}

func (m *Metrics) WithLabels(labels ...string) *Metrics {
	metrics := m.copy()
	metrics.labels = labels
	return metrics
}

func (m *Metrics) CounterInc(name string, lv ...string) {
	defer m.recover(name, "counter_inc")
	vec, err := m.store.getCounter(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "counter_inc")
		return
	}
	vec.Inc(lv...)
}

func (m *Metrics) CounterAdd(name string, v float64, lv ...string) {
	defer m.recover(name, "counter_add")
	vec, err := m.store.getCounter(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "counter_add")
		return
	}
	vec.Add(v, lv...)
}

func (m *Metrics) GaugeSet(name string, v float64, lv ...string) {
	defer m.recover(name, "gauge_set")
	vec, err := m.store.getGauge(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "gauge_set")
		return
	}
	vec.Set(v, lv...)
}

func (m *Metrics) GaugeInc(name string, lv ...string) {
	defer m.recover(name, "gauge_inc")
	vec, err := m.store.getGauge(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "gauge_inc")
		return
	}
	vec.Inc(lv...)
}

func (m *Metrics) GaugeDec(name string, lv ...string) {
	defer m.recover(name, "gauge_dec")
	vec, err := m.store.getGauge(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "gauge_dec")
		return
	}
	vec.Dec(lv...)
}

func (m *Metrics) GaugeAdd(name string, v float64, lv ...string) {
	defer m.recover(name, "gauge_add")
	vec, err := m.store.getGauge(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "gauge_add")
		return
	}
	vec.Add(v, lv...)
}

func (m *Metrics) GaugeSub(name string, v float64, lv ...string) {
	defer m.recover(name, "gauge_sub")
	vec, err := m.store.getGauge(m.registerer, prefixedName(m.prefix, name, m.separator), m.labels...)
	if err != nil {
		m.emitError(err, name, "gauge_sub")
		return
	}
	vec.Sub(v, lv...)
}

func (m *Metrics) SummaryObserve(name string, v float64, lv ...string) {
	defer m.recover(name, "summary_observe")
	vec, err := m.store.getSummary(m.registerer, prefixedName(m.prefix, name, m.separator), *m.summaryOpts, m.labels...)
	if err != nil {
		m.emitError(err, name, "summary_timer")
		return
	}
	vec.Observe(v, lv...)
}

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

func (m *Metrics) HistogramObserve(name string, v float64, lv ...string) {
	defer m.recover(name, "histogram_observe")
	vec, err := m.store.getHistogram(m.registerer, prefixedName(m.prefix, name, m.separator), m.histogramBuckets, m.labels...)
	if err != nil {
		m.emitError(err, name, "histogram_observe")
		return
	}
	vec.Observe(v, lv...)
}

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

func (m *Metrics) copy() *Metrics {
	return &Metrics{
		prefix:           m.prefix,
		separator:        m.separator,
		port:             m.port,
		path:             m.path,
		bindAddr:         m.bindAddr,
		histogramBuckets: m.histogramBuckets,
		summaryOpts:      m.summaryOpts,
		store:            m.store,
		labels:           m.labels,
		panicOnError:     m.panicOnError,
		errors:           m.errors,
		registry:         m.registry,
		registerer:       m.registerer,
	}
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

func defaulted(opts MetricsOpts) MetricsOpts {
	// opts.Prefix default is empty
	// opts.MustRegister default is false
	// opts.PanicOnError default is false
	if opts.Registry == nil {
		opts.Registry = prometheus.NewRegistry()
	}

	if opts.HistogramBuckets == nil {
		opts.HistogramBuckets = prometheus.DefBuckets
	}

	if opts.BindAddr == "" {
		opts.BindAddr = "0.0.0.0"
	}

	if opts.Path == "" {
		opts.Path = "/metrics"
	}

	if opts.Port == 0 {
		opts.Port = 9090
	}

	if opts.Separator == 0 {
		opts.Separator = '_'
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
			m.errors.PanicRecovery(name, fn)
		}
	}
}
