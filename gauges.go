package apex

import (
	"ctx.sh/apex/utils"
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

func (g *Gauges) Get(name string, labels prometheus.Labels) (*prometheus.GaugeVec, error) {
	if metric, can := g.metrics[name]; can {
		return metric, nil
	}

	return g.Register(name, utils.LabelKeys(labels))
}

func (g *Gauges) Register(name string, labels []string) (*prometheus.GaugeVec, error) {
	n, err := utils.NameBuilder(g.namespace, g.subsystem, name, g.separator)
	if err != nil {
		return nil, err
	}

	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: n,
		Help: "created automagically by apex",
	}, labels)

	if err := utils.Register(gauge); err != nil {
		return nil, err
	}

	g.metrics[name] = gauge
	return gauge, nil
}

func (g *Gauges) Set(name string, value float64, labels prometheus.Labels) error {
	gauge, err := g.Get(name, labels)
	if err != nil {
		return err
	}

	gauge.With(prometheus.Labels(labels)).Set(value)
	return nil
}

func (g *Gauges) Inc(name string, labels prometheus.Labels) error {
	gauge, err := g.Get(name, labels)
	if err != nil {
		return err
	}

	gauge.With(prometheus.Labels(labels)).Inc()
	return nil
}

func (g *Gauges) Dec(name string, labels prometheus.Labels) error {
	gauge, err := g.Get(name, labels)
	if err != nil {
		return err
	}

	gauge.With(prometheus.Labels(labels)).Dec()
	return nil
}

func (g *Gauges) Add(name string, value float64, labels prometheus.Labels) error {
	gauge, err := g.Get(name, labels)
	if err != nil {
		return err
	}

	gauge.With(prometheus.Labels(labels)).Add(value)
	return nil
}

func (g *Gauges) Sub(name string, value float64, labels prometheus.Labels) error {
	gauge, err := g.Get(name, labels)
	if err != nil {
		return err
	}

	gauge.With(prometheus.Labels(labels)).Sub(value)
	return nil
}
