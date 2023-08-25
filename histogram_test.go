package strata

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestHistogram(t *testing.T) {
	reg := prometheus.NewPedanticRegistry()
	vec, err := NewHistogramVec(reg, "test_hst", prometheus.DefBuckets)
	assert.NoError(t, err)

	vec.Observe(10.0)
	CollectAndCompare(t, vec, "test_hst", "histogram", nil, 10.0)
}
