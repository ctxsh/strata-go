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

import (
	"github.com/prometheus/client_golang/prometheus"
)

// CounterVec is a wrapper around the prometheus CounterVec.
//
// It bundles a set of Counters that all share the same Desc, but have different
// values for their variable labels. This is used if you want to count the same
// thing partitioned by various dimensions (e.g. number of HTTP requests,
// partitioned by response code and method).
type CounterVec struct {
	vec  *prometheus.CounterVec
	name string
}

// NewCounterVec creates, registers, and returns a new CounterVec.
func NewCounterVec(registerer prometheus.Registerer, name string, labels ...string) (*CounterVec, error) {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: DefaultHelpString,
	}, labels)

	if err := Register(registerer, counter); err != nil {
		return nil, err
	}

	return &CounterVec{
		name: name,
		vec:  counter,
	}, nil
}

// Inc increments the counter by 1 with the label values in the order that
// the labels were defined in NewCounterVec.
func (c *CounterVec) Inc(lv ...string) {
	c.vec.WithLabelValues(lv...).Inc()
}

// Add increases the counter by the given float value with the label values
// in the order that the labels were defined in NewCounterVec.
func (c *CounterVec) Add(v float64, lv ...string) {
	c.vec.WithLabelValues(lv...).Add(v)
}

// Name returns the name of the CounterVec.
func (c *CounterVec) Name() string {
	return c.name
}

// Type returns the metric type.
func (c *CounterVec) Type() MetricType {
	return CounterType
}

// Vec returns the prometheus CounterVec.
func (c *CounterVec) Vec() prometheus.Collector {
	return c.vec
}

var _ MetricVec = &CounterVec{}
