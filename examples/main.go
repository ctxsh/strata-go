package main

import (
	"math/rand"
	"sync"
	"time"

	"ctx.sh/apex"
)

var (
	histogramOpts = apex.HistogramOpts{
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}
	summaryOpts = apex.SummaryOpts{
		MaxAge:     5 * time.Minute,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		AgeBuckets: 5,
	}
	labels = apex.Labels{"region": "us-east-1"}
)

func random(min int, max int) float64 {
	return float64(min) + rand.Float64()*(float64(max-min))
}

func runOnce(m *apex.Metrics) {
	// Histogram timer
	timer := m.HistogramTimer("latency", apex.Labels{
		"func":   "runOnce",
		"region": "us-east-1",
	}, apex.HistogramOpts{})
	defer timer.ObserveDuration()

	// Counter functions
	m.CounterInc("test_counter", labels)
	m.CounterAdd("test_counter", 5.0, labels)

	// Gauge functions
	m.GaugeInc("test_gauge", labels)
	m.GaugeSet("test_gauge", random(1, 100), labels)
	m.GaugeAdd("test_gauge", 2.0, labels)
	m.GaugeSub("test_gauge", 1.0, labels)

	// Summary observation
	m.SummaryObserve("test_summary", random(0, 10), labels, summaryOpts)

	delay := time.Duration(random(10, 1500)) * time.Millisecond
	time.Sleep(delay)
}

func main() {
	var wg sync.WaitGroup

	metrics := apex.New(apex.MetricsOpts{
		Namespace:    "apex",
		Subsystem:    "example",
		Separator:    '_',
		PanicOnError: true,
	})

	wg.Add(1)
	go func() {
		_ = metrics.Start()
	}()

	wg.Add(1)
	go func() {
		for {
			runOnce(metrics)
		}
	}()
	wg.Wait()
}
