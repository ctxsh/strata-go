package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/ctxswitch/apex"
)

func random(min int, max int) float64 {
	return float64(min) + rand.Float64()*(float64(max-min))
}

func runOnce(m *apex.Metrics) {
	timer := m.HistogramTimer("latency", apex.Labels{
		"func":   "runOnce",
		"region": "us-east-1",
	}, 0.5, 0.9, 0.99, 0.999, 1.0)
	defer timer.ObserveDuration()

	m.CounterInc("inc_counter", apex.Labels{"region": "us-east-1"})
	m.CounterAdd("add_counter", 5.0, apex.Labels{"region": "us-east-1"})
	m.GaugeInc("test_gauge", apex.Labels{"region": "us-east-1"})
	m.GaugeSet("test_gauge", random(1, 100), apex.Labels{"region": "us-east-1"})
	delay := time.Duration(random(500, 1500)) * time.Millisecond
	time.Sleep(delay)
}

func main() {
	var wg sync.WaitGroup

	metrics := apex.New(apex.MetricsOpts{
		Namespace:    "apex",
		Subsystem:    "example",
		MustRegister: true,
		Separator:    ':',
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
