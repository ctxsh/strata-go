package main

import (
	"sync"
	"time"

	"github.com/ctxswitch/apex"
)

func run(m *apex.Metrics) {
	for {
		m.Inc("inc_counter", apex.Labels{"region": "us-east-1"})
		m.Incv("incv_counter", 5.0, apex.Labels{"region": "us-east-1"})
		time.Sleep(1 * time.Second)
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

	metrics.Register(apex.Counter, "inc_counter", []string{"region"})
	metrics.Register(apex.Counter, "incv_counter", []string{"region"})
	metrics.Start(wg)

	wg.Add(1)
	go run(metrics)
	wg.Wait()
}
