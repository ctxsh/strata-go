# Apex Metrics Package

Wrappers around the prometheus client.

TODO: More documentation.

```go
import "github.com/ctxswitch/go-apex"

var metrics apex.Metrics

func init() {
  metrics = apex.New(apex.MetricsOpts{
    Namespace:    "apex",
    Subsystem:    "example",
    MustRegister: true,
    Separator:    ':',
  })
}

func main() {
  metrics.NewCounter("my_counter_metric", []string{"env"})
  metrics.NewGauge("my_gauge_metric", []string{"region"})
  metrics.NewHistogram(
    "my_latency_metric",
    []string{"what"},
    []float64{0.5, 0.9, 0.99}
  )
  m.Start()

  run()
}

func run() {
  timer := metrics.NewTimer("my_latency_metric", apex.Labels{"what": "something"})
  defer timer.ObserveDuration()

  // Do stuff
  metrics.CounterInc("my_counter_metric", apex.Labels{"env": "production"})
  metrics.GaugeSet("my_gauge_metric", 100, apex.Labels{"region": "us-east-1"})
}
```
