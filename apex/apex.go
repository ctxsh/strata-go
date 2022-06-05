package apex

import (
	"fmt"
	"net/http"

	"github.com/ctxswitch/apex-go/errors"
	"github.com/ctxswitch/apex-go/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsOpts struct {
	Namespace    string
	Subsystem    string
	Path         string
	Port         int
	Separator    rune
	PanicOnError bool
}

type Metrics struct {
	opts MetricsOpts

	counters   *metric.Counters
	gauges     *metric.Gauges
	histograms *metric.Histograms
	summaries  *metric.Summaries
	errors     *errors.ApexInternalErrorMetrics
}

func New(opts MetricsOpts) *Metrics {
	m := &Metrics{
		counters:   metric.NewCounters(opts.Namespace, opts.Subsystem, opts.Separator),
		gauges:     metric.NewGauges(opts.Namespace, opts.Subsystem, opts.Separator),
		histograms: metric.NewHistograms(opts.Namespace, opts.Subsystem, opts.Separator),
		summaries:  metric.NewSummaries(opts.Namespace, opts.Subsystem, opts.Separator),
	}
	m.opts = defaults(opts)
	m.errors = errors.NewApexInternalErrorMetrics(opts.Namespace, opts.Subsystem, opts.Separator)
	return m
}

func (m *Metrics) Start() error {
	mux := http.NewServeMux()
	mux.Handle(m.opts.Path, promhttp.Handler())
	addr := fmt.Sprintf("localhost:%d", m.opts.Port)
	return http.ListenAndServe(addr, mux)
}

func (m *Metrics) CounterInc(name string, labels Labels) {
	defer m.recover(name, "CounterInc")

	if err := m.counters.Inc(name, prometheus.Labels(labels)); err != nil {
		m.emitError(err, name, "CounterInc")
	}
}

func (m *Metrics) CounterAdd(name string, value float64, labels Labels) {
	defer m.recover(name, "CounterAdd")

	if err := m.counters.Add(name, value, prometheus.Labels(labels)); err != nil {
		m.emitError(err, name, "CounterAdd")
	}
}

func (m *Metrics) GaugeSet(name string, value float64, labels Labels) {
	defer m.recover(name, "GaugeSet")

	if err := m.gauges.Set(name, value, prometheus.Labels(labels)); err != nil {
		m.emitError(err, name, "GaugeSet")
	}
}

func (m *Metrics) GaugeInc(name string, labels Labels) {
	defer m.recover(name, "GaugeInc")

	if err := m.gauges.Inc(name, prometheus.Labels(labels)); err != nil {
		m.emitError(err, name, "GaugeInc")
	}
}

func (m *Metrics) GaugeDec(name string, labels Labels) {
	defer m.recover(name, "GaugeDec")

	if err := m.gauges.Dec(name, prometheus.Labels(labels)); err != nil {
		m.emitError(err, name, "GaugeDec")
	}
}

func (m *Metrics) GaugeAdd(name string, value float64, labels Labels) {
	defer m.recover(name, "GaugeAdd")

	if err := m.gauges.Add(name, value, prometheus.Labels(labels)); err != nil {
		m.emitError(err, name, "GaugeAdd")
	}
}

func (m *Metrics) GaugeSub(name string, value float64, labels Labels) {
	defer m.recover(name, "GaugeSub")

	if err := m.gauges.Sub(name, value, prometheus.Labels(labels)); err != nil {
		m.emitError(err, name, "GaugeSub")
	}
}

func (m *Metrics) HistogramObserve(name string, value float64, labels Labels, buckets ...float64) {
	defer m.recover(name, "HistogramObserve")

	if err := m.histograms.Observe(name, value, prometheus.Labels(labels), buckets...); err != nil {
		m.emitError(err, name, "HistogramObserve")
	}
}

func (m *Metrics) HistogramTimer(name string, labels Labels, buckets ...float64) *metric.Timer {
	defer m.recover(name, "HistogramTimers")

	timer, err := m.histograms.Timer(name, prometheus.Labels(labels), buckets...)
	if err != nil {
		m.emitError(err, name, "HistogramTimer")
	}
	return timer
}

func (m *Metrics) SummaryObserve(name string, value float64, labels Labels) {
	defer m.recover(name, "SummaryObserve")

	if err := m.summaries.Observe(name, value, prometheus.Labels(labels)); err != nil {
		m.emitError(err, name, "SummaryObserve")
	}
}

func (m *Metrics) SummaryTimer(name string, labels Labels) *metric.Timer {
	defer m.recover(name, "SummaryTimers")

	timer, err := m.summaries.Timer(name, prometheus.Labels(labels))
	if err != nil {
		m.emitError(err, name, "counter_inc")
	}
	return timer
}

func defaults(opts MetricsOpts) MetricsOpts {
	// opts.Namespace default is empty
	// opts.Subsystem default is empty
	// opts.MustRegister default is false
	// opts.PanicOnError default is false
	if opts.Path == "" {
		opts.Path = "/metrics"
	}

	if opts.Port == 0 {
		opts.Port = 9000
	}

	if opts.Separator == 0 {
		opts.Separator = '_'
	}

	return opts
}

func (m *Metrics) emitError(err error, name string, fn string) {
	if m.opts.PanicOnError {
		panic(err)
	}

	switch err {
	case errors.ErrInvalidMetricName:
		m.errors.InvalidMetricName(name, fn)
	case errors.ErrRegistrationFailed:
		m.errors.RegistrationFailed(name, fn)
	case errors.ErrAlreadyRegistered:
		m.errors.AlreadyRegistered(name, fn)
	}
}

func (m *Metrics) recover(name string, fn string) {
	if r := recover(); r != nil && !m.opts.PanicOnError {
		m.errors.PanicRecovery(name, fn)
	}
}
