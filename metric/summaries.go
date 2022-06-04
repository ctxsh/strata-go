package metric

import (
	"github.com/ctxswitch/apex/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type Summaries struct {
	metrics   map[string]*prometheus.SummaryVec
	namespace string
	subsystem string
	separator rune
}

func NewSummaries(ns string, sub string, sep rune) *Summaries {
	return &Summaries{
		metrics:   make(map[string]*prometheus.SummaryVec),
		namespace: ns,
		subsystem: sub,
		separator: sep,
	}
}

func (s *Summaries) Get(name string, labels prometheus.Labels) (*prometheus.SummaryVec, error) {
	if metric, can := s.metrics[name]; can {
		return metric, nil
	}

	return s.Register(name, utils.LabelKeys(labels))
}

func (s *Summaries) Register(name string, labels []string) (*prometheus.SummaryVec, error) {
	n, err := utils.NameBuilder(s.namespace, s.subsystem, name, s.separator)
	if err != nil {
		return nil, err
	}

	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: n,
		Help: "created automagically by apex",
	}, labels)

	if err := utils.Register(summary); err != nil {
		return nil, err
	}

	s.metrics[name] = summary
	return summary, nil
}

func (s *Summaries) Observe(name string, value float64, labels prometheus.Labels) error {
	summary, err := s.Get(name, labels)
	if err != nil {
		return err
	}
	summary.With(prometheus.Labels(labels)).Observe(value)
	return nil
}

func (s *Summaries) Timer(name string, labels prometheus.Labels) (*Timer, error) {
	summary, err := s.Get(name, labels)
	if err != nil {
		return nil, err
	}

	return NewTimer(summary, labels), nil
}
