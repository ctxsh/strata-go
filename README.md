# Apex Metrics Package

Wrappers around the prometheus client.

TODO: More documentation.

```go
metrics := apex.New(apex.MetricsOpts{
		Namespace:    "apex",
		Subsystem:    "example",
		MustRegister: true,
		Separator:    ':',
	})

	metrics.NewCounter("my_counter_metric", []string{"env"})
	metrics.NewGauge("my_gauge_metric", []string{"region"})
  metrics.NewHistogram("my_latency_metric", []string{"path"}, []float64{0.5, 0.9, 0.99})

  m.CounterInc("my_counter_metric", apex.Labels{"env": "production"})
	m.GaugeSet("my_gauge_metric", 100, apex.Labels{"region": "us-east-1"})

  timer := metrics.NewTimer("my_latency_metric", apex.Labels{"path": "/blog"})
  defer timer.ObserveDuration()

```
