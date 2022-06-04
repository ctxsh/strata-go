package metric

import (
	"github.com/ctxswitch/apex/utils"
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

func (h *Histograms) Get(name string, labels prometheus.Labels, buckets ...float64) *prometheus.HistogramVec {
	if metric, can := h.metrics[name]; can {
		return metric
	}

	return h.Register(name, utils.LabelKeys(labels), buckets...)
}

func (h *Histograms) Register(name string, labels []string, buckets ...float64) *prometheus.HistogramVec {
	n, err := utils.NameBuilder(h.namespace, h.subsystem, name, h.separator)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	h.metrics[name] = histogram
	return histogram
}

func (h *Histograms) Observe(name string, value float64, labels prometheus.Labels, buckets ...float64) {
	if histogram := h.Get(name, labels, buckets...); histogram != nil {
		histogram.With(prometheus.Labels(labels)).Observe(value)
	} else {
		panic("unanticipated error occured")
	}
}

func (h *Histograms) Timer(name string, labels prometheus.Labels, buckets ...float64) *prometheus.Timer {
	if histogram := h.Get(name, labels, buckets...); histogram != nil {
		return prometheus.NewTimer(histogram.With(
			prometheus.Labels(labels),
		))
	} else {
		panic("unanticipated error occured")
	}
}
