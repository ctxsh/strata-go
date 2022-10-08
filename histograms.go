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

type HistogramOpts struct {
	Buckets []float64
}

type Histograms struct {
	metrics   map[string]*prometheus.HistogramVec
	namespace string
	subsystem interface{}
	separator rune
}

func NewHistograms(ns string, sub interface{}, sep rune) *Histograms {
	return &Histograms{
		metrics:   make(map[string]*prometheus.HistogramVec),
		namespace: ns,
		subsystem: sub,
		separator: sep,
	}
}

func (h *Histograms) Get(name string, labels Labels, opts HistogramOpts) (*prometheus.HistogramVec, error) {
	if metric, can := h.metrics[name]; can {
		return metric, nil
	}

	return h.Register(name, labels.Keys(), opts)
}

func (h *Histograms) Register(name string, labels []string, opts HistogramOpts) (*prometheus.HistogramVec, error) {
	n, err := NameBuilder(h.namespace, h.subsystem, name, h.separator)
	if err != nil {
		return nil, err
	}

	if opts.Buckets == nil {
		opts.Buckets = prometheus.DefBuckets
	}

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    n,
		Help:    "created automagically by apex",
		Buckets: opts.Buckets,
	}, labels)

	if err := Register(histogram); err != nil {
		return nil, err
	}

	h.metrics[name] = histogram
	return histogram, nil
}

func (h *Histograms) Observe(name string, value float64, labels Labels, opts HistogramOpts) error {
	histogram, err := h.Get(name, labels, opts)
	if err != nil {
		return err
	}
	histogram.With(prometheus.Labels(labels)).Observe(value)
	return nil
}

func (h *Histograms) Timer(name string, labels Labels, opts HistogramOpts) (*Timer, error) {
	histogram, err := h.Get(name, labels, opts)
	if err != nil {
		return nil, err
	}

	return NewTimer(histogram, labels), nil
}

func LinearBuckets(start, width float64, count int) []float64 {
	return prometheus.LinearBuckets(start, width, count)
}

func ExponentialBuckets(start, factor float64, count int) []float64 {
	return prometheus.ExponentialBuckets(start, factor, count)
}

func ExponentialBucketRange(min, max float64, count int) []float64 {
	return prometheus.ExponentialBucketsRange(min, max, count)
}
