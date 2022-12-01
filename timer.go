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

// Timer is a helper type to time functions.
type Timer struct {
	timer *prometheus.Timer
}

// NewTimer creates a new Timer.
func NewTimer(collector prometheus.Collector, lv ...string) *Timer {
	t := &Timer{}
	switch metric := collector.(type) {
	case *prometheus.HistogramVec:
		t.timer = prometheus.NewTimer(metric.WithLabelValues(lv...))
	case *prometheus.SummaryVec:
		t.timer = prometheus.NewTimer(metric.WithLabelValues(lv...))
	default:
		t.timer = nil
	}

	return t
}

// Func allows the use of ordinary functions as Observers.
func (t *Timer) Func(name string, fn func(float64)) *Timer {
	return &Timer{
		timer: prometheus.NewTimer(prometheus.ObserverFunc(fn)),
	}
}

// ObserveDuration records the duration that has passed between the time that
// the Timer was created.
func (t *Timer) ObserveDuration() {
	if t.timer != nil {
		t.timer.ObserveDuration()
	}
}
