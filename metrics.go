package apex

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metric int
type Labels prometheus.Labels

const (
	Counter Metric = iota
	Gauge
)

type MetricsOpts struct {
	Namespace    string
	Subsystem    string
	MustRegister bool
	Path         string
	Port         int
	Separator    rune
}

type Metrics struct {
	counters map[string]*prometheus.CounterVec
	opts     MetricsOpts

	mErrorUnknown  *prometheus.CounterVec
	mErrorInvalid  *prometheus.CounterVec
	mPanicRecovery *prometheus.CounterVec
}

func New(opts MetricsOpts) *Metrics {
	m := &Metrics{
		counters: make(map[string]*prometheus.CounterVec),
	}
	m.opts = defaults(opts)
	m.init()
	return m
}

func (m *Metrics) Register(ptype Metric, name string, labels []string) {
	switch ptype {
	case Counter:
		metric, err := m.newCounter(name, labels)
		if err != nil {
			m.ErrorInvalidInc(Labels{
				"name": name,
				"kind": "counter",
			})
		}
		if err := m.register(metric); err == nil {
			m.counters[name] = metric
		}
	}

}

func (m *Metrics) Start(wg sync.WaitGroup) {
	mux := http.NewServeMux()
	mux.Handle(m.opts.Path, promhttp.Handler())

	wg.Add(1)
	go func() {
		addr := fmt.Sprintf("localhost:%d", m.opts.Port)
		_ = http.ListenAndServe(addr, mux)
	}()
}

func (m *Metrics) init() {
	var name string

	name, _ = m.nameBuilder("error_unknown")
	m.mErrorUnknown = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: m.opts.Namespace,
		Subsystem: m.opts.Subsystem,
		Name:      name,
	}, []string{"name", "kind"})

	name, _ = m.nameBuilder("error_invalid")
	m.mErrorInvalid = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: m.opts.Namespace,
		Subsystem: m.opts.Subsystem,
		Name:      name,
	}, []string{"name", "kind"})

	name, _ = m.nameBuilder("panic_recovery")
	m.mPanicRecovery = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: m.opts.Namespace,
		Subsystem: m.opts.Subsystem,
		Name:      name,
	}, []string{"name", "method"})

	// This is the only place we want to panic as if these are
	// broken, we don't have any insights into what is going on
	// with the collectors.
	prometheus.MustRegister(m.mErrorUnknown)
	prometheus.MustRegister(m.mErrorInvalid)
	prometheus.MustRegister(m.mPanicRecovery)
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

func (m *Metrics) register(metric prometheus.Collector) error {
	defer m.recover("notset", "recover")
	err := prometheus.Register(metric)
	if err != nil && m.opts.MustRegister {
		panic(err)
	} else if err != nil {
		return RegistrationFailed
	}
	return nil
}

func (m *Metrics) recover(name string, method string) {
	// Add logging through the logging interface later
	if r := recover(); r != nil {
		m.PanicRecoveryInc(Labels{
			"name":   name,
			"method": method,
		})
	}
}

func (m *Metrics) ErrorUnknownInc(labels Labels) {
	m.mErrorUnknown.With(prometheus.Labels(labels)).Inc()
}

func (m *Metrics) ErrorInvalidInc(labels Labels) {
	m.mErrorInvalid.With(prometheus.Labels(labels)).Inc()
}

func (m *Metrics) PanicRecoveryInc(labels Labels) {
	m.mPanicRecovery.With(prometheus.Labels(labels)).Inc()
}
