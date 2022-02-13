package apex

import "github.com/prometheus/client_golang/prometheus"

func (m *Metrics) RegisterCounter(name string, labels []string) {
	n, err := m.nameBuilder(name)
	if err != nil {
		m.metrics[name] = m.mErrorInvalid
		return
	}

	collector := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: n,
		Help: "created automagically by apex",
	}, labels)

	if err := m.register(collector); err == nil {
		m.metrics[name] = collector
	}
}

func (m *Metrics) RegisterGauge(name string, labels []string) {
	n, err := m.nameBuilder(name)
	if err != nil {
		m.metrics[name] = m.mErrorInvalid
		return
	}

	collector := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: n,
		Help: "created automagically by apex",
	}, labels)

	if err := m.register(collector); err == nil {
		m.metrics[name] = collector
	}
}

func (m *Metrics) RegisterSummary(name string, labels []string) {
	n, err := m.nameBuilder(name)
	if err != nil {
		m.metrics[name] = m.mErrorInvalid
		return
	}

	collector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: n,
		Help: "created automagically by apex",
	}, labels)

	if err := m.register(collector); err == nil {
		m.metrics[name] = collector
	}
}

func (m *Metrics) RegisterHistogram(name string, labels []string, buckets []float64) {
	n, err := m.nameBuilder(name)
	if err != nil {
		m.metrics[name] = m.mErrorInvalid
		return
	}

	collector := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    n,
		Help:    "created automagically by apex",
		Buckets: buckets,
	}, labels)

	if err := m.register(collector); err == nil {
		m.metrics[name] = collector
	}
}

func (m *Metrics) register(metric prometheus.Collector) error {
	defer m.recover("notset", "recover")
	if err := prometheus.Register(metric); err != nil {
		_, ok := err.(prometheus.AlreadyRegisteredError)
		if ok {
			// Metric is already registered, so ignore. I may update this
			// later to also respect MustRegister. The only reason why I did
			// this was to fix tests.  I think I may be able to resolve this
			// with configurable registries later.
			return nil
		} else if m.opts.MustRegister {
			panic(err)
		} else {
			return RegistrationFailed
		}
	}
	return nil
}
