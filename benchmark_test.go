package apex

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func BenchmarkMain(b *testing.B) {
	registry := prometheus.NewPedanticRegistry()

	metrics := New(MetricsOpts{
		Separator:    '_',
		Registry:     registry,
		PanicOnError: true,
	}).WithPrefix("apex", "example").WithLabels("role")

	for i := 0; i < 1000000; i++ {
		metrics.CounterInc("foo", "example")
	}

	for i := 0; i < 1000000; i++ {
		metrics.CounterAdd("foo", 5.0, "example")
	}
}
