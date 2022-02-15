package apex

import "github.com/prometheus/client_golang/prometheus"

func (m *Metrics) NewTimer(name string, labels Labels) *prometheus.Timer {
	if metric, can := m.metrics[name]; can {
		switch metric := metric.(type) {
		case *prometheus.HistogramVec:
			return prometheus.NewTimer(metric.With(
				prometheus.Labels(labels),
			))
		case *prometheus.SummaryVec:
			return prometheus.NewTimer(metric.With(
				prometheus.Labels(labels),
			))
		}
	}
	return prometheus.NewTimer(m.mInvalidTimer.WithLabelValues(name))
}

func (m *Metrics) NewTimerFunc(name string, fn func(float64)) *prometheus.Timer {
	return prometheus.NewTimer(prometheus.ObserverFunc(fn))
}
