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

import "github.com/prometheus/client_golang/prometheus"

// GaugeVec is a wrapper around the prometheus GaugeVec.
//
// It bundles a set of Gauges that all share the same Desc, but have different
// values for their variable labels. This is used if you want to count the same
// thing partitioned by various dimensions (e.g. number of operations queued,
// partitioned by user and operation type). Create instances with NewGaugeVec.
//
// A gauge represents a numerical value that can be arbitrarily increased or
// decreased.  Gauges are typically used for measured values like temperatures
// or current memory usage, but also "counts" that can go up and down.  Gauges
// are often used to represent things like disk and memory usage and concurrent
// requests.
type GaugeVec struct {
	name string
	vec  *prometheus.GaugeVec
}

// NewGaugeVec creates, registers, and returns a new GaugeVec.
func NewGaugeVec(registerer prometheus.Registerer, name string, labels ...string) (*GaugeVec, error) {
	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: DefaultHelpString,
	}, labels)

	if err := Register(registerer, gauge); err != nil {
		return nil, err
	}

	return &GaugeVec{
		name: name,
		vec:  gauge,
	}, nil
}

// Set sets the Gauge to an arbitrary value using the label values in the order that
// the labels were defined in NewGaugeVec.
func (g *GaugeVec) Set(v float64, lv ...string) {
	g.vec.WithLabelValues(lv...).Set(v)
}

// Inc increments the Gauge by 1 using the label values in the order that the labels
// were defined in NewGaugeVec.
func (g *GaugeVec) Inc(lv ...string) {
	g.vec.WithLabelValues(lv...).Inc()
}

// Dec decrements the Gauge by 1 using the label values in the order that the labels
// were defined in NewGaugeVec.
func (g *GaugeVec) Dec(lv ...string) {
	g.vec.WithLabelValues(lv...).Dec()
}

// Add increases the counter by the given float value with the label values in the
// order that the labels were defined in NewGaugeVec.
func (g *GaugeVec) Add(v float64, lv ...string) {
	g.vec.WithLabelValues(lv...).Add(v)
}

// Add subtracts the counter by the given float value with the label values in the
// order that the labels were defined in NewGaugeVec.
func (g *GaugeVec) Sub(v float64, lv ...string) {
	g.vec.WithLabelValues(lv...).Sub(v)
}

// Name returns the name of the GaugeVec.
func (g *GaugeVec) Name() string {
	return g.name
}

// Type returns the metric type.
func (g *GaugeVec) Type() MetricType {
	return GaugeType
}

// Vec returns the prometheus GaugeVec.
func (g *GaugeVec) Vec() prometheus.Collector {
	return g.vec
}

var _ MetricVec = &GaugeVec{}
