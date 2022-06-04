package metric

import "github.com/prometheus/client_golang/prometheus"

type Timer struct{}

func (t *Timer) Func(name string, fn func(float64)) *prometheus.Timer {
	return prometheus.NewTimer(prometheus.ObserverFunc(fn))
}
