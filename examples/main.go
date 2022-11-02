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
package main

import (
	"math/rand"
	"sync"
	"time"

	"ctx.sh/apex"
)

func random(min int, max int) float64 {
	return float64(min) + rand.Float64()*(float64(max-min))
}

func runOnce(m *apex.Metrics) {
	// Histogram timer
	timer := m.HistogramTimer("latency")
	defer timer.ObserveDuration()

	// Counter functions
	n := m.WithLabels("label")
	n.CounterInc("test_counter", "value1")
	n.CounterAdd("test_counter", 5.0, "value1")

	// Gauge functions
	n.GaugeInc("test_gauge", "value2")
	n.GaugeSet("test_gauge", random(1, 100), "value2")
	n.GaugeAdd("test_gauge", 2.0, "value2")
	n.GaugeSub("test_gauge", 1.0, "value2")

	// Summary observation
	n.SummaryObserve("test_summary", random(0, 10), "value3")

	delay := time.Duration(random(1, 1500)) * time.Millisecond
	time.Sleep(delay)
}

func main() {
	var wg sync.WaitGroup
	metrics := apex.New(apex.MetricsOpts{
		Separator:      ':',
		PanicOnError:   true,
		Port:           9090,
		ConstantLabels: []string{"role", "server"},
		SummaryOpts: &apex.SummaryOpts{
			MaxAge:     10 * time.Minute,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			AgeBuckets: 5,
		},
		HistogramBuckets: []float64{.01, .025, .05, .1, .25, .5, 1, 2.5},
	}).WithPrefix("apex", "example")

	wg.Add(1)
	go func() {
		err := metrics.Start()
		if err != nil {
			panic(err)
		}
	}()

	wg.Add(1)
	go func() {
		for {
			metrics.CounterInc("loop_total")
			runOnce(metrics.WithPrefix("runonce"))
		}
	}()
	wg.Wait()
}
