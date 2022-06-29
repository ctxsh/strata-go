package apex

import (
	"github.com/prometheus/client_golang/prometheus"
)

type HistogramOpts struct {
	Buckets []float64
}

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

func (h *Histograms) Get(name string, labels Labels, opts HistogramOpts) (*prometheus.HistogramVec, error) {
	if metric, can := h.metrics[name]; can {
		return metric, nil
	}

	return h.Register(name, labels.Keys(), opts)
}

func (h *Histograms) Register(name string, labels []string, opts HistogramOpts) (*prometheus.HistogramVec, error) {
	n, err := NameBuilder(h.namespace, h.subsystem, name, h.separator)
	if err != nil {
		return nil, err
	}

	if opts.Buckets == nil {
		opts.Buckets = prometheus.DefBuckets
	}

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    n,
		Help:    "created automagically by apex",
		Buckets: opts.Buckets,
	}, labels)

	if err := Register(histogram); err != nil {
		return nil, err
	}

	h.metrics[name] = histogram
	return histogram, nil
}

func (h *Histograms) Observe(name string, value float64, labels Labels, opts HistogramOpts) error {
	histogram, err := h.Get(name, labels, opts)
	if err != nil {
		return err
	}
	histogram.With(prometheus.Labels(labels)).Observe(value)
	return nil
}

func (h *Histograms) Timer(name string, labels Labels, opts HistogramOpts) (*Timer, error) {
	histogram, err := h.Get(name, labels, opts)
	if err != nil {
		return nil, err
	}

	return NewTimer(histogram, labels), nil
}

func LinearBuckets(start, width float64, count int) []float64 {
	return prometheus.LinearBuckets(start, width, count)
}

func ExponentialBuckets(start, factor float64, count int) []float64 {
	return prometheus.ExponentialBuckets(start, factor, count)
}

func ExponentialBucketRange(min, max float64, count int) []float64 {
	return prometheus.ExponentialBucketsRange(min, max, count)
}
