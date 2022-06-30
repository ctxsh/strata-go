# Apex Prometheus Client [![unit tests](https://github.com/ctxswitch/apex-go/actions/workflows/test.yml/badge.svg)](https://github.com/ctxswitch/apex-go/actions/workflows/test.yml)

The Apex Go package provides a wrapper around the prometheus client to automatically register and collect metrics.

## Install

```
go get ctx.sh/apex
```

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
| BindAddr | `0.0.0.0` | The address the promethus collector will listen on for connections |
| Namespace | empty | The prefix for a metric |
| Subsystem | empty | A string that represents the subsystem.  This value is joined to the namespace with the defined seperator |
| Separator | `_` | The seperator that will be used to join the metric name components. |
| Path | `/metrics` | The path used by the HTTP server |
| Port | `9000` | The port used by the HTTP server |
| PanicOnError | `false` | Maintain the default behavior of prometheus to panic on errors.  If this value is set to false, the library attempts to recover from any panics and emits an internally managed metric `apex_errors_panic_recovery` to inform the operator that visibility is degraded.  If set to true the original behavior is maintained and all errors are treated as panics. |  

### Counter

A counter is a cumulative metric whose value can only increase or be reset to zero on restart. Counters are often used to represent the number of requests served, tasks completed, or errors.

#### `CounterInc`

Increment a counter metric by one.

```go
metrics.CounterInc("my_counter", apex.Labels{"app": "my_app"})
```

#### `CounterAdd`

Add a float64 value to the counter metric. 

```go
metrics.CounterAdd("my_counter", 2.0, apex.Labels{"app": "my_app"})
```

### Gauge

A gauge represents a numerical value that can be arbitrarily increased or decreased.  Gauges are typically used for measured values like temperatures or current memory usage, but also "counts" that can go up and down.  Gauges are often used to represent things like disk and memory usage and concurrent requests. 

#### `GaugeSet`

Set a gauge to the value that is passed.

```go
metrics.GaugeSet("my_gauge", 2.0, apex.Labels{"app": "my_app"})
```

#### `GaugeInc`

Increment a gauge by one.

```go
metrics.GaugeInc("my_gauge", apex.Labels{"app": "my_app"})
```

#### `GaugeDec`

Decrement a gauge by one.

```go
metrics.GaugeDec("my_gauge", apex.Labels{"app": "my_app"})
```

#### `GaugeAdd`

Adds the value that is passed to a gauge.

```go
metrics.GaugeAdd("my_gauge", 2.0, apex.Labels{"app": "my_app"})
```

#### `GaugeSub`

Subtract the value that is passed to a gauge.

```go
metrics.GaugeSub("my_gauge", 2.0, apex.Labels{"app": "my_app"})
```

### Histogram

A histogram samples observations and counts them in configurable buckets. Most often histograms are used to measure durations or sizes.  Histograms expose multiple measurements during a scrape.  These include bucket measurements in the format `<name>_bucket{le="<upper_bound>"}`, the total sum of observed values as `<name>_sum`, and the number of observered events in the format of `<name>_count`.  Histograms are configurable through the `HistogramOpts` struct which allow overrides of the following attributes:

* `Buckets`: The time buckets into which observations are counted.  Values must be sorted in increasing order and .  The `+inf` bucket is automatically added to catch values .

#### `HistogramObserve`

Add a single observation to the histogram.

```go
metrics.HistogramObserve("my_histogram", response_time, apex.Labels{
	"app": "my_app",
}, apex.HistogramOpts{
	Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
})
```

#### `Histogram Timer`

Create a histogram timer. 

```go
timer := m.HistogramTimer("response", apex.Labels{
		"path":   "/api/comments",
		"region": "us-east-1",
	}, apex.HistogramOpts{})
defer timer.ObserveDuration()
```

### Summary

A summary samples observations and calculates quantiles over a sliding time windo.  Like histograms, they are used to measure durations or sizes.  Summaries expose multiple measurements during a scrape.  Thiese include quantiles in the form of `<name>{quantile="φ"}`, , the total sum of observed values as `<name>_sum`, and the number of observered events in the format of `<name>_count`.  Summaries are configurable through the SummaryOpts struct which allow overrides of the following attributes:

* `Objectives`: The quantile rank estimates with their respective absolute error defined as a `map[float64]float64`.
* `MaxAge`: The duration that observations stay relevant as `time.Duration`.
* `AgeBuckets`: Number of buckets used to calculate the age of observations as a `uint32`.

#### `SummaryObserve`

Add a single observations to the summary

```go
metrics.SummaryObserve("test_summary", response, apex.Labels{
	"site": "api.example.com",
}, apex.SummaryOpts{
	MaxAge:     5 * time.Minute,
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	AgeBuckets: 5,
})
```

#### `SummaryTimer`

Create a summary timer. 

```go
timer := m.SummaryTimer("response", apex.Labels{
		"path":   "/api/comments",
		"region": "us-east-1",
	}, apex.SummaryOpts{})
defer timer.ObserveDuration()
```
