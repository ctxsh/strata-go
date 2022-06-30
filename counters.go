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

type Counters struct {
	metrics   map[string]*prometheus.CounterVec
	namespace string
	subsystem string
	separator rune
}

func NewCounters(ns string, sub string, sep rune) *Counters {
	return &Counters{
		metrics:   make(map[string]*prometheus.CounterVec),
		namespace: ns,
		subsystem: sub,
		separator: sep,
	}
}

func (c *Counters) Get(name string, labels Labels) (*prometheus.CounterVec, error) {
	if metric, can := c.metrics[name]; can {
		return metric, nil
	}

	return c.Register(name, labels.Keys())
}

func (c *Counters) Register(name string, labels []string) (*prometheus.CounterVec, error) {
	n, err := NameBuilder(c.namespace, c.subsystem, name, c.separator)
	if err != nil {
		return nil, err
	}

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: n,
		Help: "created automagically by apex",
	}, labels)

	if err := Register(counter); err != nil {
		return nil, err
	}

	c.metrics[name] = counter
	return counter, nil
}

func (c *Counters) Inc(name string, labels Labels) error {
	counter, err := c.Get(name, labels)
	if err != nil {
		return err
	}

	counter.With(prometheus.Labels(labels)).Inc()
	return nil
}

func (c *Counters) Add(name string, value float64, labels Labels) error {
	counter, err := c.Get(name, labels)
	if err != nil {
		return err
	}

	counter.With(prometheus.Labels(labels)).Add(value)
	return nil
}
