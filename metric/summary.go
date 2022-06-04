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

func (s *Summaries) Get(name string, labels prometheus.Labels) *prometheus.SummaryVec {
	if metric, can := s.metrics[name]; can {
		return metric
	}

	return s.Register(name, utils.LabelKeys(labels))
}

func (s *Summaries) Register(name string, labels []string) *prometheus.SummaryVec {
	n, err := utils.NameBuilder(s.namespace, s.subsystem, name, s.separator)
	if err != nil {
		panic(err)
	}

	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: n,
		Help: "created automagically by apex",
	}, labels)

	if err := utils.Register(summary); err != nil {
		panic(err)
	}

	s.metrics[name] = summary
	return summary
}

func (s *Summaries) Observe(name string, value float64, labels prometheus.Labels) {
	if summary := s.Get(name, labels); summary != nil {
		summary.With(prometheus.Labels(labels)).Observe(value)
	} else {
		panic("unanticipated error occured")
	}
}

func (s *Summaries) Timer(name string, labels prometheus.Labels) *prometheus.Timer {
	if summary := s.Get(name, labels); summary != nil {
		return prometheus.NewTimer(summary.With(
			prometheus.Labels(labels),
		))
	} else {
		panic("unanticipated error occured")
	}
}
