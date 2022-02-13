package apex

import "github.com/prometheus/client_golang/prometheus"

func (m *Metrics) HistogramObserve(name string, value float64, labels Labels) {
	defer m.recover(name, "histogram_observe")
	if histogram := m.getHistogram(name); histogram != nil {
		// Do something different
		histogram.With(prometheus.Labels(labels)).Observe(value)
	}
}

func (m *Metrics) getHistogram(name string) *prometheus.HistogramVec {
	if metric, can := m.metrics[name]; can {
		return metric.(*prometheus.HistogramVec)
	}
	return nil
}
