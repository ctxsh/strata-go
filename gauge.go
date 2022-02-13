package apex

import "github.com/prometheus/client_golang/prometheus"

func (m *Metrics) GaugeSet(name string, value float64, labels Labels) {
	defer m.recover(name, "gauge_inc")

	if gauge := m.getGauge(name); gauge != nil {
		gauge.With(prometheus.Labels(labels)).Set(value)
	} else {
		m.mInvalidGauge.WithLabelValues(name).Inc()
	}
}

func (m *Metrics) GaugeInc(name string, labels Labels) {
	defer m.recover(name, "gauge_inc")

	if gauge := m.getGauge(name); gauge != nil {
		gauge.With(prometheus.Labels(labels)).Inc()
	} else {
		m.mInvalidGauge.WithLabelValues(name).Inc()
	}
}

func (m *Metrics) GaugeDec(name string, labels Labels) {
	defer m.recover(name, "gauge_dec")

	if gauge := m.getGauge(name); gauge != nil {
		gauge.With(prometheus.Labels(labels)).Inc()
	} else {
		m.mInvalidGauge.WithLabelValues(name).Inc()
	}
}

func (m *Metrics) GaugeAdd(name string, value float64, labels Labels) {
	defer m.recover(name, "gauge_add")
	if gauge := m.getGauge(name); gauge != nil {
		gauge.With(prometheus.Labels(labels)).Add(value)
	} else {
		m.mInvalidGauge.WithLabelValues(name).Inc()
	}
}

func (m *Metrics) GaugeSub(name string, value float64, labels Labels) {
	defer m.recover(name, "gauge_add")
	if gauge := m.getGauge(name); gauge != nil {
		gauge.With(prometheus.Labels(labels)).Sub(value)
	} else {
		m.mInvalidGauge.WithLabelValues(name).Inc()
	}
}

func (m *Metrics) getGauge(name string) *prometheus.GaugeVec {
	if metric, can := m.metrics[name]; can {
		return metric.(*prometheus.GaugeVec)
	}
	return nil
}
