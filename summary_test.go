package apex

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestSummary(t *testing.T) {
	reg := prometheus.NewPedanticRegistry()
	vec, err := NewSummaryVec(reg, "test_smy", SummaryOpts{})
	assert.NoError(t, err)

	vec.Observe(10.0)
	CollectAndCompare(t, vec, "test_smy", "summary", nil, 10.0)
}
