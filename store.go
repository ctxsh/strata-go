package strata

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Store manages all of the prometheus collectors.
type Store struct {
	counters   map[string]*CounterVec
	gauges     map[string]*GaugeVec
	summaries  map[string]*SummaryVec
	histograms map[string]*HistogramVec
	// TODO: part of the issue with the race condition was that we were
	// setting the metric store value to nil and not revisiting.  This will
	// pretty much address the double register race that caused the nil, but
	// I need to come back through this and make the check/get more resilient
	// so I can shrink the footprint of the lock.
	sync.Mutex
}

func newStore() *Store {
	return &Store{
		counters:   make(map[string]*CounterVec),
		gauges:     make(map[string]*GaugeVec),
		summaries:  make(map[string]*SummaryVec),
		histograms: make(map[string]*HistogramVec),
	}
}

func (s *Store) getCounter(reg prometheus.Registerer, name string, labels ...string) (*CounterVec, error) {
	s.Lock()
	defer s.Unlock()

	if vec, ok := s.counters[name]; ok {
		return vec, nil
	}

	vec, err := NewCounterVec(reg, name, labels...)
	s.counters[name] = vec
	return vec, err
}

func (s *Store) getGauge(reg prometheus.Registerer, name string, labels ...string) (*GaugeVec, error) {
	s.Lock()
	defer s.Unlock()

	if vec, ok := s.gauges[name]; ok {
		return vec, nil
	}

	vec, err := NewGaugeVec(reg, name, labels...)
	s.gauges[name] = vec
	return vec, err
}

func (s *Store) getSummary(reg prometheus.Registerer, name string, opts SummaryOpts, labels ...string) (*SummaryVec, error) {
	s.Lock()
	defer s.Unlock()

	if vec, ok := s.summaries[name]; ok {
		return vec, nil
	}

	vec, err := NewSummaryVec(reg, name, opts, labels...)
	s.summaries[name] = vec
	return vec, err
}

func (s *Store) getHistogram(reg prometheus.Registerer, name string, buckets []float64, labels ...string) (*HistogramVec, error) {
	s.Lock()
	defer s.Unlock()

	if vec, ok := s.histograms[name]; ok {
		return vec, nil
	}

	vec, err := NewHistogramVec(reg, name, buckets, labels...)
	s.histograms[name] = vec
	return vec, err
}
