package apex

import (
	"github.com/prometheus/client_golang/prometheus"
)

// TODO: Can we handle panics around labels better than a catchall?
func (m *Metrics) CounterInc(name string, labels Labels) {
	defer m.recover(name, "inc")
	if counter := m.getCounter(name); counter != nil {
		counter.With(prometheus.Labels(labels)).Inc()
	} else {
		m.mInvalidCounter.WithLabelValues(name).Inc()
	}
}

func (m *Metrics) CounterAdd(name string, value float64, labels Labels) {
	defer m.recover(name, "add")
	if counter := m.getCounter(name); counter != nil {
		counter.With(prometheus.Labels(labels)).Add(value)
	} else {
		m.mInvalidCounter.WithLabelValues(name).Inc()
	}
}

func (m *Metrics) getCounter(name string) *prometheus.CounterVec {
	if metric, can := m.metrics[name]; can {
		return metric.(*prometheus.CounterVec)
	}
	return nil
}
