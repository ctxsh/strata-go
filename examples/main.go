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

func run(m *apex.Metrics) {
	for {
		m.CounterInc("inc_counter", apex.Labels{"region": "us-east-1"})
		m.CounterAdd("add_counter", 5.0, apex.Labels{"region": "us-east-1"})
		time.Sleep(1 * time.Second)
		m.GaugeInc("test_gauge", apex.Labels{"region": "us-east-1"})
		m.GaugeSet("test_gauge", random(1, 100), apex.Labels{"region": "us-east-1"})
	}
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
	metrics.Start(wg)

	wg.Add(1)
	go run(metrics)
	wg.Wait()
}
