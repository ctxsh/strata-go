package apex

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestGauge(t *testing.T) {
	reg := prometheus.NewPedanticRegistry()
	vec, err := NewGaugeVec(reg, "test_g")
	assert.NoError(t, err)

	vec.Inc()
	CollectAndCompare(t, vec, "test_g", "gauge", nil, 1.0)

	vec.Add(5.0)
	CollectAndCompare(t, vec, "test_g", "gauge", nil, 6.0)

	vec.Set(10.0)
	CollectAndCompare(t, vec, "test_g", "gauge", nil, 10.0)

	vec.Dec()
	CollectAndCompare(t, vec, "test_g", "gauge", nil, 9.0)

	vec.Sub(9.0)
	CollectAndCompare(t, vec, "test_g", "gauge", nil, 0.0)
}

func TestGaugeWithLabels(t *testing.T) {
	labels := map[string]string{
		"label": "one",
	}

	reg := prometheus.NewPedanticRegistry()
	vec, err := NewGaugeVec(reg, "test_g", "label")
	assert.NoError(t, err)

	vec.Inc("one")
	CollectAndCompare(t, vec, "test_g", "gauge", labels, 1.0)

	vec.Add(5.0, "one")
	CollectAndCompare(t, vec, "test_g", "gauge", labels, 6.0)

	vec.Set(10.0, "one")
	CollectAndCompare(t, vec, "test_g", "gauge", labels, 10.0)

	vec.Dec("one")
	CollectAndCompare(t, vec, "test_g", "gauge", labels, 9.0)

	vec.Sub(9.0, "one")
	CollectAndCompare(t, vec, "test_g", "gauge", labels, 0.0)
}
