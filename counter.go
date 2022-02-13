package apex

import (
	"github.com/prometheus/client_golang/prometheus"
)

func (m *Metrics) Inc(name string, labels Labels) {
	defer m.recover(name, "inc")

	if counter := m.getCounter(name); counter != nil {
		counter.With(prometheus.Labels(labels)).Inc()
	} else {
		m.ErrorUnknownInc(Labels{
			"name": name,
			"kind": "counter",
		})
	}
}

func (m *Metrics) Incv(name string, value float64, labels Labels) {
	defer m.recover(name, "add")
	if counter := m.getCounter(name); counter != nil {
		counter.With(prometheus.Labels(labels)).Add(value)
	} else {
		m.ErrorUnknownInc(Labels{
			"name": name,
			"kind": "counter",
		})
	}
}

func (m *Metrics) getCounter(name string) *prometheus.CounterVec {
	if counter, has := m.counters[name]; has {
		return counter
	}
	return nil
}

func (m *Metrics) newCounter(name string, labels []string) (*prometheus.CounterVec, error) {
	n, err := m.nameBuilder(name)
	if err != nil {
		return nil, err
	}
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: n,
		Help: "created automagically by apex",
	}, labels)

	return counter, nil
}
