# Apex Metrics Package
Wrappers around the prometheus client


```go
m := New(MetricsOpts{
  Namespace: "apex",
  Subsystem: "httpserver",
  Port: 8080,
  MustRegister: true,
}).Start(wg)

m.Register(metrics.Counter, "http:response:count", []string{"path", "code"})

m.Inc("http:response:count", Labels{"path": "/", "code": 204})
```
