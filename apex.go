package apex

import (
	"fmt"
	"net/http"

	"ctx.sh/apex/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsOpts struct {
	BindAddr     string
	Namespace    string
	Subsystem    string
	Path         string
	Port         int
	Separator    rune
	PanicOnError bool
}

type Metrics struct {
	opts MetricsOpts

	counters   *Counters
	gauges     *Gauges
	histograms *Histograms
	summaries  *Summaries
	errors     *errors.ApexInternalErrorMetrics
}

func New(opts MetricsOpts) *Metrics {
	m := &Metrics{
		counters:   NewCounters(opts.Namespace, opts.Subsystem, opts.Separator),
		gauges:     NewGauges(opts.Namespace, opts.Subsystem, opts.Separator),
		histograms: NewHistograms(opts.Namespace, opts.Subsystem, opts.Separator),
		summaries:  NewSummaries(opts.Namespace, opts.Subsystem, opts.Separator),
	}
	m.opts = defaults(opts)
	m.errors = errors.NewApexInternalErrorMetrics(opts.Namespace, opts.Subsystem, opts.Separator)
	return m
}

func (m *Metrics) Start() error {
	mux := http.NewServeMux()
	mux.Handle(m.opts.Path, promhttp.Handler())
	addr := fmt.Sprintf("%s:%d", m.opts.BindAddr, m.opts.Port)
	return http.ListenAndServe(addr, mux)
}

func (m *Metrics) CounterInc(name string, labels Labels) {
	defer m.recover(name, "CounterInc")

	if err := m.counters.Inc(name, labels); err != nil {
		m.emitError(err, name, "CounterInc")
	}
}

func (m *Metrics) CounterAdd(name string, value float64, labels Labels) {
	defer m.recover(name, "CounterAdd")

	if err := m.counters.Add(name, value, labels); err != nil {
		m.emitError(err, name, "CounterAdd")
	}
}

func (m *Metrics) GaugeSet(name string, value float64, labels Labels) {
	defer m.recover(name, "GaugeSet")

	if err := m.gauges.Set(name, value, labels); err != nil {
		m.emitError(err, name, "GaugeSet")
	}
}

func (m *Metrics) GaugeInc(name string, labels Labels) {
	defer m.recover(name, "GaugeInc")

	if err := m.gauges.Inc(name, labels); err != nil {
		m.emitError(err, name, "GaugeInc")
	}
}

func (m *Metrics) GaugeDec(name string, labels Labels) {
	defer m.recover(name, "GaugeDec")

	if err := m.gauges.Dec(name, labels); err != nil {
		m.emitError(err, name, "GaugeDec")
	}
}

func (m *Metrics) GaugeAdd(name string, value float64, labels Labels) {
	defer m.recover(name, "GaugeAdd")

	if err := m.gauges.Add(name, value, labels); err != nil {
		m.emitError(err, name, "GaugeAdd")
	}
}

func (m *Metrics) GaugeSub(name string, value float64, labels Labels) {
	defer m.recover(name, "GaugeSub")

	if err := m.gauges.Sub(name, value, labels); err != nil {
		m.emitError(err, name, "GaugeSub")
	}
}

func (m *Metrics) HistogramObserve(name string, value float64, labels Labels, opts HistogramOpts) {
	defer m.recover(name, "HistogramObserve")

	if err := m.histograms.Observe(name, value, labels, opts); err != nil {
		m.emitError(err, name, "HistogramObserve")
	}
}

func (m *Metrics) HistogramTimer(name string, labels Labels, opts HistogramOpts) *Timer {
	defer m.recover(name, "HistogramTimers")

	timer, err := m.histograms.Timer(name, labels, opts)
	if err != nil {
		m.emitError(err, name, "HistogramTimer")
	}
	return timer
}

func (m *Metrics) SummaryObserve(name string, value float64, labels Labels, opts SummaryOpts) {
	defer m.recover(name, "SummaryObserve")

	if err := m.summaries.Observe(name, value, labels, opts); err != nil {
		m.emitError(err, name, "SummaryObserve")
	}
}

func (m *Metrics) SummaryTimer(name string, labels Labels, opts SummaryOpts) *Timer {
	defer m.recover(name, "SummaryTimers")

	timer, err := m.summaries.Timer(name, labels, opts)
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
	if opts.BindAddr == "" {
		opts.BindAddr = "0.0.0.0"
	}

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
