package metric

import (
	"ctx.sh/apex/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type Histograms struct {
	metrics   map[string]*prometheus.HistogramVec
	namespace string
	subsystem string
	separator rune
}

func NewHistograms(ns string, sub string, sep rune) *Histograms {
	return &Histograms{
		metrics:   make(map[string]*prometheus.HistogramVec),
		namespace: ns,
		subsystem: sub,
		separator: sep,
	}
}

func (h *Histograms) Get(name string, labels prometheus.Labels, buckets ...float64) (*prometheus.HistogramVec, error) {
	if metric, can := h.metrics[name]; can {
		return metric, nil
	}

	return h.Register(name, utils.LabelKeys(labels), buckets...)
}

func (h *Histograms) Register(name string, labels []string, buckets ...float64) (*prometheus.HistogramVec, error) {
	n, err := utils.NameBuilder(h.namespace, h.subsystem, name, h.separator)
	if err != nil {
		return nil, err
	}

	if buckets == nil {
		buckets = prometheus.DefBuckets
	}

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    n,
		Help:    "created automagically by apex",
		Buckets: buckets,
	}, labels)

	if err := utils.Register(histogram); err != nil {
		return nil, err
	}

	h.metrics[name] = histogram
	return histogram, nil
}

func (h *Histograms) Observe(name string, value float64, labels prometheus.Labels, buckets ...float64) error {
	histogram, err := h.Get(name, labels, buckets...)
	if err != nil {
		return err
	}
	histogram.With(prometheus.Labels(labels)).Observe(value)
	return nil
}

func (h *Histograms) Timer(name string, labels prometheus.Labels, buckets ...float64) (*Timer, error) {
	histogram, err := h.Get(name, labels, buckets...)
	if err != nil {
		return nil, err
	}

	return NewTimer(histogram, labels), nil
}
