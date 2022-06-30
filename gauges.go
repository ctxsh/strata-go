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

type Gauges struct {
	metrics   map[string]*prometheus.GaugeVec
	namespace string
	subsystem string
	separator rune
}

func NewGauges(ns string, sub string, sep rune) *Gauges {
	return &Gauges{
		metrics:   make(map[string]*prometheus.GaugeVec),
		namespace: ns,
		subsystem: sub,
		separator: sep,
	}
}

func (g *Gauges) Get(name string, labels Labels) (*prometheus.GaugeVec, error) {
	if metric, can := g.metrics[name]; can {
		return metric, nil
	}

	return g.Register(name, labels.Keys())
}

func (g *Gauges) Register(name string, labels []string) (*prometheus.GaugeVec, error) {
	n, err := NameBuilder(g.namespace, g.subsystem, name, g.separator)
	if err != nil {
		return nil, err
	}

	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: n,
		Help: "created automagically by apex",
	}, labels)

	if err := Register(gauge); err != nil {
		return nil, err
	}

	g.metrics[name] = gauge
	return gauge, nil
}

func (g *Gauges) Set(name string, value float64, labels Labels) error {
	gauge, err := g.Get(name, labels)
	if err != nil {
		return err
	}

	gauge.With(prometheus.Labels(labels)).Set(value)
	return nil
}

func (g *Gauges) Inc(name string, labels Labels) error {
	gauge, err := g.Get(name, labels)
	if err != nil {
		return err
	}

	gauge.With(prometheus.Labels(labels)).Inc()
	return nil
}

func (g *Gauges) Dec(name string, labels Labels) error {
	gauge, err := g.Get(name, labels)
	if err != nil {
		return err
	}

	gauge.With(prometheus.Labels(labels)).Dec()
	return nil
}

func (g *Gauges) Add(name string, value float64, labels Labels) error {
	gauge, err := g.Get(name, labels)
	if err != nil {
		return err
	}

	gauge.With(prometheus.Labels(labels)).Add(value)
	return nil
}

func (g *Gauges) Sub(name string, value float64, labels Labels) error {
	gauge, err := g.Get(name, labels)
	if err != nil {
		return err
	}

	gauge.With(prometheus.Labels(labels)).Sub(value)
	return nil
}
