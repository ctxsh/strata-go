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

type CounterVec struct {
	vec  *prometheus.CounterVec
	name string
}

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

func (c *CounterVec) Inc(lv ...string) {
	c.vec.WithLabelValues(lv...).Inc()
}

func (c *CounterVec) Add(v float64, lv ...string) {
	c.vec.WithLabelValues(lv...).Add(v)
}

func (c *CounterVec) Name() string {
	return c.name
}

func (c *CounterVec) Type() MetricType {
	return CounterType
}

func (c *CounterVec) Vec() prometheus.Collector {
	return c.vec
}

var _ MetricVec = &CounterVec{}
