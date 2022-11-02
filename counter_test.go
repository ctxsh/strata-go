package apex

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	reg := prometheus.NewPedanticRegistry()
	vec, err := NewCounterVec(reg, "test_total")
	assert.NoError(t, err)

	vec.Inc()
	CollectAndCompare(t, vec, "test_total", "counter", nil, 1.0)

	vec.Add(5.0)
	CollectAndCompare(t, vec, "test_total", "counter", nil, 6.0)
}

func TestCounterWithLabels(t *testing.T) {
	labels := map[string]string{
		"label": "one",
	}

	reg := prometheus.NewPedanticRegistry()
	vec, err := NewCounterVec(reg, "test_total", "label")
	assert.NoError(t, err)

	vec.Inc("one")
	CollectAndCompare(t, vec, "test_total", "counter", labels, 1.0)

	vec.Add(5.0, "one")
	CollectAndCompare(t, vec, "test_total", "counter", labels, 6.0)
}
