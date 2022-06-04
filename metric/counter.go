package metric

import (
	"github.com/ctxswitch/apex/utils"
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

func (c *Counters) Get(name string, labels prometheus.Labels) *prometheus.CounterVec {
	if metric, can := c.metrics[name]; can {
		return metric
	}

	return c.Register(name, utils.LabelKeys(labels))
}

func (c *Counters) Register(name string, labels []string) *prometheus.CounterVec {
	n, err := utils.NameBuilder(c.namespace, c.subsystem, name, c.separator)
	if err != nil {
		panic(err)
	}

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: n,
		Help: "created automagically by apex",
	}, labels)

	if err := utils.Register(counter); err != nil {
		panic(err)
	}

	c.metrics[name] = counter
	return counter
}

func (c *Counters) Inc(name string, labels prometheus.Labels) {
	if counter := c.Get(name, labels); counter != nil {
		counter.With(prometheus.Labels(labels)).Inc()
	} else {
		panic("unanticipated error occured")
	}
}

func (c *Counters) Add(name string, value float64, labels prometheus.Labels) {
	if counter := c.Get(name, labels); counter != nil {
		counter.With(prometheus.Labels(labels)).Add(value)
	} else {
		panic("unanticipated error occured")
	}
}
