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
	timer := m.NewTimer("test_histogram", apex.Labels{"what": "something"})
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
	})

	metrics.NewCounter("inc_counter", []string{"region"})
	metrics.NewCounter("add_counter", []string{"region"})
	metrics.NewGauge("test_gauge", []string{"region"})
	metrics.NewHistogram(
		"test_histogram",
		[]string{"what"},
		[]float64{0.5, 0.9, 0.99},
	)

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
