// Copyright (C) 2022, Rob Lyon <rob@ctxswitch.com>
//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package apex

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DefaultObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001} //nolint:gochecknoglobals
)

// SummaryVec is a wrapper around the prometheus SummaryVec.
//
// It bundles a set of summaries used if you want to count the same thing
// partitioned by various dimensions.
type SummaryVec struct {
	name string
	vec  *prometheus.SummaryVec
}

func NewSummaryVec(registerer prometheus.Registerer, name string, opts SummaryOpts, labels ...string) (*SummaryVec, error) {
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       name,
		Help:       DefaultHelpString,
		Objectives: opts.Objectives,
		MaxAge:     opts.MaxAge,
		AgeBuckets: opts.AgeBuckets,
	}, labels)

	if err := Register(registerer, summary); err != nil {
		return nil, err
	}

	return &SummaryVec{
		name: name,
		vec:  summary,
	}, nil
}

// Observe adds a single observation to the summary.
func (s *SummaryVec) Observe(v float64, lv ...string) {
	s.vec.WithLabelValues(lv...).Observe(v)
}

// Timer returns a new summary timer.
func (s *SummaryVec) Timer(lv ...string) *Timer {
	return NewTimer(s.vec, lv...)
}

// Name returns the name of the SummaryVec.
func (s *SummaryVec) Name() string {
	return s.name
}

// Type returns the metric type.
func (s *SummaryVec) Type() MetricType {
	return SummaryType
}

// Vec returns the prometheus SummaryVec.
func (s *SummaryVec) Vec() prometheus.Collector {
	return s.vec
}

var _ MetricVec = &SummaryVec{}
