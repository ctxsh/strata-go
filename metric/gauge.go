package metric

import (
	"github.com/ctxswitch/apex/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type Gauges struct {
	metrics   map[string]*prometheus.GaugeVec
	namespace string
	subsystem string
	separator rune
}

func NewGauges(ns string, sub string, sep rune) *Gauges {
	return &Gauges{
		metrics:   make(map[string]*prometheus.GaugeVec),
		namespace: ns,
		subsystem: sub,
		separator: sep,
	}
}

func (g *Gauges) Get(name string, labels prometheus.Labels) *prometheus.GaugeVec {
	if metric, can := g.metrics[name]; can {
		return metric
	}

	return g.Register(name, utils.LabelKeys(labels))
}

func (g *Gauges) Register(name string, labels []string) *prometheus.GaugeVec {
	n, err := utils.NameBuilder(g.namespace, g.subsystem, name, g.separator)
	if err != nil {
		panic(err)
	}

	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: n,
		Help: "created automagically by apex",
	}, labels)

	if err := utils.Register(gauge); err != nil {
		panic(err)
	}

	g.metrics[name] = gauge
	return gauge
}

func (g *Gauges) Set(name string, value float64, labels prometheus.Labels) {
	if gauge := g.Get(name, labels); gauge != nil {
		gauge.With(prometheus.Labels(labels)).Set(value)
	} else {
		panic("unanticipated error occured")
	}
}

func (g *Gauges) Inc(name string, labels prometheus.Labels) {
	if gauge := g.Get(name, labels); gauge != nil {
		gauge.With(prometheus.Labels(labels)).Inc()
	} else {
		panic("unanticipated error occured")
	}
}

func (g *Gauges) Dec(name string, labels prometheus.Labels) {
	if gauge := g.Get(name, labels); gauge != nil {
		gauge.With(prometheus.Labels(labels)).Dec()
	} else {
		panic("unanticipated error occured")
	}
}

func (g *Gauges) Add(name string, value float64, labels prometheus.Labels) {
	if gauge := g.Get(name, labels); gauge != nil {
		gauge.With(prometheus.Labels(labels)).Add(value)
	} else {
		panic("unanticipated error occured")
	}
}

func (g *Gauges) Sub(name string, value float64, labels prometheus.Labels) {
	if gauge := g.Get(name, labels); gauge != nil {
		gauge.With(prometheus.Labels(labels)).Sub(value)
	} else {
		panic("unanticipated error occured")
	}
}
