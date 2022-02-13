package apex

import "github.com/prometheus/client_golang/prometheus"

func (m *Metrics) SummaryObserve(name string, value float64, labels Labels) {
	defer m.recover(name, "summary_observe")
	if summary := m.getSummary(name); summary != nil {
		// Do something different
		summary.With(prometheus.Labels(labels)).Observe(value)
	}
}

func (m *Metrics) getSummary(name string) *prometheus.SummaryVec {
	if metric, can := m.metrics[name]; can {
		return metric.(*prometheus.SummaryVec)
	}
	return nil
}
