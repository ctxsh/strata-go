package apex

import (
	"fmt"
	"net/http"

	"github.com/ctxswitch/apex/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsOpts struct {
	Namespace    string
	Subsystem    string
	MustRegister bool
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
}

func New(opts MetricsOpts) *Metrics {
	m := &Metrics{
		counters:   metric.NewCounters(opts.Namespace, opts.Subsystem, opts.Separator),
		gauges:     metric.NewGauges(opts.Namespace, opts.Subsystem, opts.Separator),
		histograms: metric.NewHistograms(opts.Namespace, opts.Subsystem, opts.Separator),
		summaries:  metric.NewSummaries(opts.Namespace, opts.Subsystem, opts.Separator),
	}
	m.opts = defaults(opts)
	m.init()
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

	m.counters.Inc(name, prometheus.Labels(labels))
}

func (m *Metrics) CounterAdd(name string, value float64, labels Labels) {
	defer m.recover(name, "CounterAdd")

	m.counters.Add(name, value, prometheus.Labels(labels))
}

func (m *Metrics) GaugeSet(name string, value float64, labels Labels) {
	defer m.recover(name, "GaugeSet")

	m.gauges.Set(name, value, prometheus.Labels(labels))
}

func (m *Metrics) GaugeInc(name string, labels Labels) {
	defer m.recover(name, "GaugeInc")

	m.gauges.Inc(name, prometheus.Labels(labels))
}

func (m *Metrics) GaugeDec(name string, labels Labels) {
	defer m.recover(name, "GaugeDec")

	m.gauges.Dec(name, prometheus.Labels(labels))
}

func (m *Metrics) GaugeAdd(name string, value float64, labels Labels) {
	defer m.recover(name, "GaugeAdd")

	m.gauges.Add(name, value, prometheus.Labels(labels))
}

func (m *Metrics) GaugeSub(name string, value float64, labels Labels) {
	defer m.recover(name, "GaugeSub")

	m.gauges.Sub(name, value, prometheus.Labels(labels))
}

func (m *Metrics) HistogramObserve(name string, value float64, labels Labels, buckets ...float64) {
	defer m.recover(name, "HistogramObserve")

	m.histograms.Observe(name, value, prometheus.Labels(labels), buckets...)
}

func (m *Metrics) HistogramTimer(name string, labels Labels, buckets ...float64) *prometheus.Timer {
	defer m.recover(name, "HistogramTimers")

	return m.histograms.Timer(name, prometheus.Labels(labels), buckets...)
}

func (m *Metrics) SummaryObserve(name string, value float64, labels Labels) {
	defer m.recover(name, "SummaryObserve")

	m.summaries.Observe(name, value, prometheus.Labels(labels))
}

func (m *Metrics) SummaryTimer(name string, labels Labels) *prometheus.Timer {
	defer m.recover(name, "SummaryTimers")

	return m.summaries.Timer(name, prometheus.Labels(labels))
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

func (m *Metrics) init() {
	// Register and allow the internal metrics to panic regardless of the
	// settings.  Not sure that I like this as it is dependent on the same
	// code that it wraps.  Will look at something else.
	m.counters.Register("apex_panic_recovery", []string{"name", "method"})
	// m.counters.Register("apex_invalid_counter", []string{"name"})
	// m.counters.Register("apex_invalid_gauge", []string{"name"})
	// m.counters.Register("apex_invalid_histogram", []string{"name"})
	// m.counters.Register("apex_invalid_summary", []string{"name"})
	// m.counters.Register("apex_invalid_timer", []string{"name"})
}

func (m *Metrics) recover(name string, method string) {
	if !m.opts.PanicOnError {
		return
	}

	if r := recover(); r != nil {
		m.counters.Inc("apex_panic_recovery", prometheus.Labels{"name": name, "method": method})
	}
}
