package apex

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Labels prometheus.Labels

type MetricsOpts struct {
	Namespace    string
	Subsystem    string
	MustRegister bool
	Path         string
	Port         int
	Separator    rune
}

type Metrics struct {
	metrics map[string]prometheus.Collector
	opts    MetricsOpts

	mErrorInvalid  *prometheus.CounterVec
	mPanicRecovery *prometheus.CounterVec
	// New style of handling (try not to panic)
	mInvalidCounter *prometheus.CounterVec
	mInvalidGauge   *prometheus.GaugeVec
	mInvalidTimer   *prometheus.HistogramVec
}

func New(opts MetricsOpts) *Metrics {
	m := &Metrics{
		metrics: make(map[string]prometheus.Collector),
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

func (m *Metrics) nameBuilder(name string) (string, error) {
	var builder strings.Builder

	if m.opts.Namespace != "" {
		builder.WriteString(m.opts.Namespace)
		builder.WriteRune(m.opts.Separator)
	}

	if m.opts.Subsystem != "" {
		builder.WriteString(m.opts.Subsystem)
		builder.WriteRune(m.opts.Separator)
	}

	if name == "" {
		return "", InvalidMetricName
	}
	builder.WriteString(name)
	return builder.String(), nil
}

func (m *Metrics) recover(name string, method string) {
	if r := recover(); r != nil {
		m.mPanicRecovery.With(prometheus.Labels{
			"name":   name,
			"method": method,
		}).Inc()
	}
}

func (m *Metrics) init() {
	var name string

	name, _ = m.nameBuilder("panic_recovery")
	m.mPanicRecovery = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
	}, []string{"name", "method"})

	name, _ = m.nameBuilder("panic_recovery")
	m.mErrorInvalid = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
	}, []string{"name", "method"})

	name, _ = m.nameBuilder("invalid_counter")
	m.mInvalidCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
	}, []string{"name"})

	name, _ = m.nameBuilder("invalid_gauge")
	m.mInvalidGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
	}, []string{"name"})

	name, _ = m.nameBuilder("invalid_timer")
	m.mInvalidTimer = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: m.opts.Namespace,
		Subsystem: m.opts.Subsystem,
		Name:      name,
	}, []string{"name"})

	_ = m.register(m.mErrorInvalid)
	_ = m.register(m.mPanicRecovery)
	_ = m.register(m.mInvalidCounter)
	_ = m.register(m.mInvalidGauge)
	_ = m.register(m.mInvalidTimer)
}
