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

func (s *SummaryVec) Observe(v float64, lv ...string) {
	s.vec.WithLabelValues(lv...).Observe(v)
}

func (s *SummaryVec) Timer(lv ...string) *Timer {
	return NewTimer(s.vec, lv...)
}

func (s *SummaryVec) Name() string {
	return s.name
}

func (s *SummaryVec) Type() MetricType {
	return SummaryType
}

func (s *SummaryVec) Vec() prometheus.Collector {
	return s.vec
}

var _ MetricVec = &SummaryVec{}
