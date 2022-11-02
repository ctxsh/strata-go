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

type GaugeVec struct {
	name string
	vec  *prometheus.GaugeVec
}

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

func (g *GaugeVec) Set(v float64, lv ...string) {
	g.vec.WithLabelValues(lv...).Set(v)
}

func (g *GaugeVec) Inc(lv ...string) {
	g.vec.WithLabelValues(lv...).Inc()
}

func (g *GaugeVec) Dec(lv ...string) {
	g.vec.WithLabelValues(lv...).Dec()
}

func (g *GaugeVec) Add(v float64, lv ...string) {
	g.vec.WithLabelValues(lv...).Add(v)
}

func (g *GaugeVec) Sub(v float64, lv ...string) {
	g.vec.WithLabelValues(lv...).Sub(v)
}

func (g *GaugeVec) Name() string {
	return g.name
}

func (g *GaugeVec) Type() MetricType {
	return GaugeType
}

func (g *GaugeVec) Vec() prometheus.Collector {
	return g.vec
}

var _ MetricVec = &GaugeVec{}
