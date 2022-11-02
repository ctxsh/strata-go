# Apex - Prometheus Client [![unit tests](https://github.com/ctxswitch/apex-go/actions/workflows/test.yml/badge.svg)](https://github.com/ctxswitch/apex-go/actions/workflows/test.yml)

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
	ConstantLabels: []string{"role", "server"},
	Separator:    ':',
	PanicOnError: false,
})
```

#### MetricOpts

| Option | Default | Description |
|--------|---------|-------------|
| BindAddr | `0.0.0.0` | The address the promethus collector will listen on for connections |
| ConstantLabels | empty | An array of label/value pairs that will be constant across all metrics. |
| HistogramBuckets | []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10} | Buckets used for histogram observation counts |
| PanicOnError | `false` | Maintain the default behavior of prometheus to panic on errors.  If this value is set to false, the library attempts to recover from any panics and emits an internally managed metric `apex_errors_panic_recovery` to inform the operator that visibility is degraded.  If set to true the original behavior is maintained and all errors are treated as panics. |
| Path | `/metrics` | The path used by the HTTP server. |
| Port | `9090` | The port used by the HTTP server. |
| Prefix | empty | An array of strings that represent the base prefix for the metric. |
| Separator | `_` | The seperator that will be used to join the metric name components. |
| SummaryOpts | defaults | Options used for configuring summary metrics |

#### SummaryOpts

| Option | Default | Description |
|--------|---------|-------------|
| AgeBuckets | 5 | AgeBuckets is the number of buckets used to exclude observations that are older than MaxAge from the summary. |
| MaxAge | 10 minutes | MaxAge defines the duration for which an observation stays relevant for the summary. |
| Objectives | map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001} | Objectives defines the quantile rank estimates with their respective absolute error. |


### Prefixes and Labels

#### `WithLabels(...string)`

The `WithLabels` function adds labels to the metrics.  If labels are added to metrics, the subsequent metrics must include the label values.  Each metric function includes a variadic parameter that is used to pass in the values in the order that the labels were previously passed.

```go
m := apex.New(apex.MetricsOpts{})
n := m.WithLabels("label1", "label2")
n.CounterInc("a_total", "value1", "value2")
```

#### `WithPrefix(...string)`

The `WithPrefix` function appends additional prefix values to the metric name.  By default metrics are created without prefixes unless added in `MetricOpts`.  For example:

```go
m := apex.New(apex.MetricsOpts{})
// prefix: ""
m.WithPrefix("apex", "example")
// prefix: "apex_example"
m.CounterInc("a_total")
// metric: "apex_example_a_total"
n := m.WithPrefix("component")
// prefix: "apex_example_component"
n.CounterInc("b_total")
// metric: "apex_example_component_b_total"
m.CounterInc("c_total")
// metric: "apex_example_c_total"
```

### Counter

A counter is a cumulative metric whose value can only increase or be reset to zero on restart. Counters are often used to represent the number of requests served, tasks completed, or errors.

#### `CounterInc(string, ...string)`

Increment a counter metric by one.

```go
// Without labels
metrics.CounterInc("my_counter")

// With labels
metrics = metrics.WithValues("label1", "label2")
metrics.CounterInc("my_counter", "value1", "value2")
```

#### `CounterAdd(string, float64, ...string)`

Add a float64 value to the counter metric. 

```go
// Without labels
metrics.CounterAdd("my_counter", 2.0)

// With labels
metrics = metrics.WithValues("label1", "label2")
metrics.CounterAdd("my_counter", 2.0, "value1", "value2")
```

### Gauge

A gauge represents a numerical value that can be arbitrarily increased or decreased.  Gauges are typically used for measured values like temperatures or current memory usage, but also "counts" that can go up and down.  Gauges are often used to represent things like disk and memory usage and concurrent requests.

#### `GaugeSet(string, float64, ...string)`

Set a gauge to the value that is passed.

```go
// Without labels
metrics.GaugeSet("my_gauge", 2.0)

// With labels
metrics = metrics.WithValues("label1", "label2")
metrics.GaugeSet("gauge_with_values", 2.0, "value1", "value2")
```

#### `GaugeInc(string, ...string)`

Increment a gauge by one.

```go
// Without labels
metrics.GaugeInc("my_gauge")

// With labels
metrics = metrics.WithValues("label1", "label2")
metrics.GaugeInc("gauge_with_values", "value1", "value2")
```

#### `GaugeDec(string, ...string)`

Decrement a gauge by one.

```go
// Without labels
metrics.GaugeDec("my_gauge")

// With labels
metrics = metrics.WithValues("label1", "label2")
metrics.GaugeDec("gauge_with_values", "value1", "value2")
```

#### `GaugeAdd(string, float64, ...string)`

Adds the value that is passed to a gauge.

```go
// Without labels
metrics.GaugeAdd("my_gauge", 2.0)

// With labels
metrics = metrics.WithValues("label1", "label2")
metrics.GaugeAdd("gauge_with_values", 2.0, "value1", "value2")
```

#### `GaugeSub(string, float64, ...string)`

Subtract the value that is passed to a gauge.

```go
// Without labels
metrics.GaugeSub("my_gauge", 2.0)

// With labels
metrics = metrics.WithValues("label1", "label2")
metrics.GaugeSub("gauge_with_values", 2.0, "value1", "value2")
```

### Histogram

A histogram samples observations and counts them in configurable buckets. Most often histograms are used to measure durations or sizes.  Histograms expose multiple measurements during a scrape.  These include bucket measurements in the format `<name>_bucket{le="<upper_bound>"}`, the total sum of observed values as `<name>_sum`, and the number of observered events in the format of `<name>_count`.  Histograms buckets are configurable through the `Buckets` struct in `MetricsOpts` which allow overrides of the following attributes:

* `Buckets`: The time buckets into which observations are counted.  Values must be sorted in increasing order and .  The `+inf` bucket is automatically added to catch values .

#### `HistogramObserve(string, float64, ...string)`

Add a single observation to the histogram.

```go
m := apex.New(apex.MetricsOpts{
	Buckets: []float{0.01, 0.5, 0.1, 1, 5, 10}
})

// Without labels
metrics.HistogramObserve("my_histogram", response_time)

// With labels
metrics = metrics.WithValues("label1", "label2")
metrics.HistogramObserve("my_histogram", response_time, "value1", "value2")
```

#### `Histogram Timer(string, ...string)`

Create a histogram timer. 

```go
m := apex.New(apex.MetricsOpts{
	Buckets: apex.ExponentialBuckets(100, 1.2, 3)
})

// Without labels
timer := m.HistogramTimer("response")
defer timer.ObserveDuration()

// With labels
metrics = metrics.WithValues("label1", "label2")
timer := m.HistogramTimer("response", "value1", "value2")
defer timer.ObserveDuration()
```

### Summary

A summary samples observations and calculates quantiles over a sliding time windo.  Like histograms, they are used to measure durations or sizes.  Summaries expose multiple measurements during a scrape.  Thiese include quantiles in the form of `<name>{quantile="φ"}`, , the total sum of observed values as `<name>_sum`, and the number of observered events in the format of `<name>_count`.  Summaries are configurable through the SummaryOpts struct which allow overrides of the following attributes:

* `Objectives`: The quantile rank estimates with their respective absolute error defined as a `map[float64]float64`.
* `MaxAge`: The duration that observations stay relevant as `time.Duration`.
* `AgeBuckets`: Number of buckets used to calculate the age of observations as a `uint32`.

#### `SummaryObserve(string, float64, ...string)`

Add a single observations to the summary

```go
m := apex.New(apex.MetricsOpts{
	SummaryOpts: &apex.SummaryOpts{
		MaxAge:     10 * time.Minute,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		AgeBuckets: 5,
	}
})

// Without labels
metrics.SummaryObserve("test_summary", response, apex.SummaryOpts{
	MaxAge:     5 * time.Minute,
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	AgeBuckets: 5,
})

// With labels
metrics = metrics.WithValues("label1", "label2")
metrics.SummaryObserve("test_summary", response, apex.SummaryOpts{
	MaxAge:     5 * time.Minute,
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	AgeBuckets: 5,
, "value1", "value2")
```

#### `SummaryTimer(string, ...string)`

Create a summary timer. 

```go
m := apex.New(apex.MetricsOpts{
	SummaryOpts: &apex.SummaryOpts{
		MaxAge:     10 * time.Minute,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		AgeBuckets: 5,
	}
})

// Without labels
timer := m.SummaryTimer("response")
defer timer.ObserveDuration()

// With labels
metrics = metrics.WithValues("label1", "label2")
timer := m.SummaryTimer("response", apex.SummaryOpts{}, "value1", "value2")
defer timer.ObserveDuration()
```
