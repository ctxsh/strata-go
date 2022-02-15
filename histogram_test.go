package apex

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestHistogramSetMismatchedLabels(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	m.NewHistogram("test_metric_observe", []string{"this", "that"}, []float64{0.5, 0.9, 0.99})

	assert.Equal(t, 0, testutil.CollectAndCount(m.mPanicRecovery))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.HistogramObserve("test_metric_observe", 100, Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mPanicRecovery))
}

func TestHistogramAddInvalidCounter(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	assert.Equal(t, 0, testutil.CollectAndCount(m.mInvalidGauge))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.GaugeAdd("were_is_waldo", 100, Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mInvalidGauge))
}
