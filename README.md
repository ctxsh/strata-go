# Apex Metrics

The Apex Go package provides a wrapper around the prometheus client to automatically register and collect metrics.

## Usage

### Initialize

Initialize the metrics collectors using the `apex.New` function.  There are several options available.  As an example:
```golang
metrics := apex.New(apex.MetricsOpts{
  Namespace:    "apex",
  Subsystem:    "example",
  Separator:    ':',
  PanicOnError: false,
})
```

| Option | Default | Description |
|--------|---------|-------------|
| Namespace | empty | The prefix for a metric |
| Subsystem | empty | A string that represents the subsystem.  This value is joined to the namespace with the defined seperator |
| Separator | `_` | The seperator that will be used to join the metric name components. |
| Path | `/metrics` | The path used by the HTTP server |
| Port | `9000` | The port used by the HTTP server |
| PanicOnError | `false` | Maintain the default behavior of prometheus to panic on errors.  If this value is set to false, the library attempts to recover from any panics and emits an internally managed metric `apex:errors:panic_recovery` to inform the operator that visibility is degraded.  If set to true the original behavior is maintained and all errors are treated as panics. |  

### Example
```golang
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
		PanicOnError: false,
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
```
