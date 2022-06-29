package apex

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Timer struct {
	timer *prometheus.Timer
}

func NewTimer(collector prometheus.Collector, labels Labels) *Timer {
	t := &Timer{}
	switch metric := collector.(type) {
	case *prometheus.HistogramVec:
		t.timer = prometheus.NewTimer(metric.With(prometheus.Labels(labels)))
	case *prometheus.SummaryVec:
		t.timer = prometheus.NewTimer(metric.With(prometheus.Labels(labels)))
	default:
		t.timer = nil
	}

	return t
}

func (t *Timer) Func(name string, fn func(float64)) *Timer {
	return &Timer{
		timer: prometheus.NewTimer(prometheus.ObserverFunc(fn)),
	}
}

func (t *Timer) ObserveDuration() {
	if t.timer != nil {
		t.timer.ObserveDuration()
	}
}
