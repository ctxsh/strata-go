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

package strata

import "github.com/prometheus/client_golang/prometheus"

// HistogramOpts defines options that are available to the HistogramVec
// collectors.
type HistogramOpts struct {
	// Buckets defines the observation buckets for the histogram.  Each float
	// value is the upper inclusive bound of the bucket with +Inf added implicitly.
	// the default is
	Buckets []float64
}

// HistogramVec is a wrapper around the prometheus HistogramVec.
//
// It bundles a set of histograms used if you want to count the same thing
// partitioned by various dimensions.
type HistogramVec struct {
	name string
	vec  *prometheus.HistogramVec
}

// NewHistogramVec creates, registers, and returns a new HistogramVec.
func NewHistogramVec(registerer prometheus.Registerer, name string, buckets []float64, labels ...string) (*HistogramVec, error) {
	summary := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    name,
		Help:    DefaultHelpString,
		Buckets: buckets,
	}, labels)

	if err := Register(registerer, summary); err != nil {
		return nil, err
	}

	return &HistogramVec{
		name: name,
		vec:  summary,
	}, nil
}

func (h *HistogramVec) Observe(v float64, lv ...string) {
	h.vec.WithLabelValues(lv...).Observe(v)
}

func (h *HistogramVec) Timer(lv ...string) *Timer {
	return NewTimer(h.vec, lv...)
}

// Name returns the name of the HistogramVec.
func (g *HistogramVec) Name() string {
	return g.name
}

// Type returns the metric type.
func (g *HistogramVec) Type() MetricType {
	return HistogramType
}

// Vec returns the prometheus HistogramVec.
func (g *HistogramVec) Vec() prometheus.Collector {
	return g.vec
}

var _ MetricVec = &HistogramVec{}
