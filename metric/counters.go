package metric

import (
	"ctx.sh/apex/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type Counters struct {
	metrics   map[string]*prometheus.CounterVec
	namespace string
	subsystem string
	separator rune
}

func NewCounters(ns string, sub string, sep rune) *Counters {
	return &Counters{
		metrics:   make(map[string]*prometheus.CounterVec),
		namespace: ns,
		subsystem: sub,
		separator: sep,
	}
}

func (c *Counters) Get(name string, labels prometheus.Labels) (*prometheus.CounterVec, error) {
	if metric, can := c.metrics[name]; can {
		return metric, nil
	}

	return c.Register(name, utils.LabelKeys(labels))
}

func (c *Counters) Register(name string, labels []string) (*prometheus.CounterVec, error) {
	n, err := utils.NameBuilder(c.namespace, c.subsystem, name, c.separator)
	if err != nil {
		return nil, err
	}

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: n,
		Help: "created automagically by apex",
	}, labels)

	if err := utils.Register(counter); err != nil {
		return nil, err
	}

	c.metrics[name] = counter
	return counter, nil
}

func (c *Counters) Inc(name string, labels prometheus.Labels) error {
	counter, err := c.Get(name, labels)
	if err != nil {
		return err
	}

	counter.With(prometheus.Labels(labels)).Inc()
	return nil
}

func (c *Counters) Add(name string, value float64, labels prometheus.Labels) error {
	counter, err := c.Get(name, labels)
	if err != nil {
		return err
	}

	counter.With(prometheus.Labels(labels)).Add(value)
	return nil
}
