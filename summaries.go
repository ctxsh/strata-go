package apex

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	DefObjectives               = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
	DefMaxAge     time.Duration = 10 * time.Minute
	DefAgeBuckets uint32        = 5
)

type SummaryOpts struct {
	Objectives map[float64]float64
	MaxAge     time.Duration
	AgeBuckets uint32
}

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

func (s *Summaries) Get(name string, labels Labels, opts SummaryOpts) (*prometheus.SummaryVec, error) {
	if metric, can := s.metrics[name]; can {
		return metric, nil
	}

	return s.Register(name, labels.Keys(), opts)
}

func (s *Summaries) Register(name string, labels []string, opts SummaryOpts) (*prometheus.SummaryVec, error) {
	n, err := NameBuilder(s.namespace, s.subsystem, name, s.separator)
	if err != nil {
		return nil, err
	}

	if opts.AgeBuckets < 1 {
		opts.AgeBuckets = DefAgeBuckets
	}

	if opts.MaxAge < 1 {
		opts.MaxAge = DefMaxAge
	}

	if opts.Objectives == nil {
		opts.Objectives = DefObjectives
	}

	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       n,
		Help:       "created automagically by apex",
		Objectives: opts.Objectives,
		MaxAge:     opts.MaxAge,
		AgeBuckets: opts.AgeBuckets,
	}, labels)

	if err := Register(summary); err != nil {
		return nil, err
	}

	s.metrics[name] = summary
	return summary, nil
}

func (s *Summaries) Observe(name string, value float64, labels Labels, opts SummaryOpts) error {
	summary, err := s.Get(name, labels, opts)
	if err != nil {
		return err
	}
	summary.With(prometheus.Labels(labels)).Observe(value)
	return nil
}

func (s *Summaries) Timer(name string, labels Labels, opts SummaryOpts) (*Timer, error) {
	summary, err := s.Get(name, labels, opts)
	if err != nil {
		return nil, err
	}

	return NewTimer(summary, labels), nil
}
