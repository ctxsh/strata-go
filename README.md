# Strata - Prometheus Client [![unit tests](https://github.com/ctxsh/strata-go/actions/workflows/test.yml/badge.svg)](https://github.com/ctxswitch/strata-go/actions/workflows/test.yml)

The Strata Go package provides a wrapper around the prometheus client to automatically register and collect metrics.

## Install

```
go get ctx.sh/strata
```

## Usage

### Initialize

Initialize the metrics collectors using the `strata.New` function.  There are several options available.  As an example:
```golang
metrics := strata.New(strata.MetricsOpts{
	ConstantLabels: []string{"role", "server"},
	Separator:    ':',
	PanicOnError: false,
})
```

#### MetricOpts

| Option | Default | Description |
|--------|---------|-------------|
| ConstantLabels | empty | An array of label/value pairs that will be constant across all metrics. |
| HistogramBuckets | `[]float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}` | Buckets used for histogram observation counts |
| Logger | nil | Provide a logger that implements the `Logger` interface.  A valid logger must have the following methods defined: `Info(msg string, keysAndValues ...any)` and `Error(err error, msg string, keysAndValues ...any)` | 
| PanicOnError | `false` | Maintain the default behavior of prometheus to panic on errors.  If this value is set to false, the library attempts to recover from any panics and emits an internally managed metric `strata_errors_panic_recovery` to inform the operator that visibility is degraded.  If set to true the original behavior is maintained and all errors are treated as panics. |
| Prefix | empty | An array of strings that represent the base prefix for the metric. |
| Separator | `_` | The seperator that will be used to join the metric name components. |
| SummaryOpts | see below | Options used for configuring summary metrics |

#### SummaryOpts

| Option | Default | Description |
|--------|---------|-------------|
| AgeBuckets | `5` | AgeBuckets is the number of buckets used to exclude observations that are older than MaxAge from the summary. |
| MaxAge | `10 minutes` | MaxAge defines the duration for which an observation stays relevant for the summary. |
| Objectives | `map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}` | Objectives defines the quantile rank estimates with their respective absolute error. |

## Starting and Stopping the collection endpoints

### Starting

There are two options for starting the collection endpoint.  You can start the built in HTTP(S) server or retrieve the handler to register the metrics route in an existing multiplexer/request router.

To start a standard http server:

```golang
err := metrics.Start(ctx, strata.ServerOpts{
	Port: 9090,
	TerminationGracePeriod: 10 * time.Second
})
```

To start an http server with TLS support, at minimum you must provide the key and the certificate:

```golang
err := metrics.Start(ctx, strata.ServerOpts{
	Port: 9090,
	TLS: &strata.TLSOpts{
		CertFile: *certFile,
		KeyFile:  *keyFile,
	},
})
```

To retrieve the handler for use in an existing router:

```golang
mux := http.NewServeMux()
mux.Handle("/metrics", metrics.Handler())
```


#### ServerOpts

| Option | Default | Description |
|--------|---------|-------------|
| BindAddr | `0.0.0.0` | The address the promethus collector will listen on for connections |
| TerminationGracePeriod | `0` |  |
| Path | `/metrics` | The path used by the HTTP server. |
| Port | `9090` | The port used by the HTTP server. |
| TLS | see below | Options used to configure TLS for the collection endpoint |

#### TLS

| Option | Default | Description |
|--------|---------|-------------|
| CertFile | - | The path to the file containing the certificate or the certificate bundle. |
| InsecureSkipVerify | false | controls whether a client verifies the server's certificate chain and host name. |
| KeyFile | - | The path to the private key file. |
| MinVersion | TLS 1.3 | The minimum TLS version that the server will accept. |

### Shutdown the collection endpoint

The metrics http collection endpoint will shutdown automatically when the context is closed.  You can control the shutdown time by setting a grace period for the collection endpoint to remain active before shutting down to ensure that the final metrics are scraped.

```golang
metrics := strata.New(strata.MetricsOpts{})

var obs sync.WaitGroup
obs.Add(1)
go func() {
	defer obs.Done()
	_ = metrics.Start(ctx, strata.ServerOpts{})
}

var wg sync.WaitGroup
wg.Add(1)
go func() {
	defer wg.Done()
	myApp.Start()
}
wg.Wait()
obs.Wait()
```

## API

### Prefixes and Labels

#### `WithLabels(...string)`

The `WithLabels` function adds labels to the metrics.  If labels are added to metrics, the subsequent metrics must include the label values.  Each metric function includes a variadic parameter that is used to pass in the values in the order that the labels were previously passed.

```go
m := strata.New(strata.MetricsOpts{})
n := m.WithLabels("label1", "label2")
n.CounterInc("a_total", "value1", "value2")
```

#### `WithPrefix(...string)`

The `WithPrefix` function appends additional prefix values to the metric name.  By default metrics are created without prefixes unless added in `MetricOpts`.  For example:

```go
m := strata.New(strata.MetricsOpts{})
// prefix: ""
m.WithPrefix("strata", "example")
// prefix: "strata_example"
m.CounterInc("a_total")
// metric: "strata_example_a_total"
n := m.WithPrefix("component")
// prefix: "strata_example_component"
n.CounterInc("b_total")
// metric: "strata_example_component_b_total"
m.CounterInc("c_total")
// metric: "strata_example_c_total"
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

A histogram samples observations and counts them in configurable buckets. Most often histograms are used to measure durations or sizes.  Histograms expose multiple measurements during a scrape.  These include bucket measurements in the format `<name>_bucket{le="<upper_bound>"}`, the total sum of observed values as `<name>_sum`, and the number of observered events in the format of `<name>_count`.  Histograms buckets are configurable through `HistogramBuckets` in `MetricsOpts` which allow overrides the time buckets into which observations are counted.  Values must be sorted in increasing order.  The `+inf` bucket is automatically added to catch values.

#### `HistogramObserve(string, float64, ...string)`

Add a single observation to the histogram.

```go
m := strata.New(strata.MetricsOpts{
	HistogramBuckets: []float{0.01, 0.5, 0.1, 1, 5, 10}
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
m := strata.New(strata.MetricsOpts{
	HistogramBuckets: strata.ExponentialBuckets(100, 1.2, 3)
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
m := strata.New(strata.MetricsOpts{
	SummaryOpts: &strata.SummaryOpts{
		MaxAge:     10 * time.Minute,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		AgeBuckets: 5,
	}
})

// Without labels
metrics.SummaryObserve("test_summary", response)

// With labels
metrics = metrics.WithValues("label1", "label2")
metrics.SummaryObserve("test_summary", response, "value1", "value2")
```

#### `SummaryTimer(string, ...string)`

Create a summary timer. 

```go
m := strata.New(strata.MetricsOpts{
	SummaryOpts: &strata.SummaryOpts{
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
timer := m.SummaryTimer("response", "value1", "value2")
defer timer.ObserveDuration()
```
