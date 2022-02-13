package apex

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRegisterInvalidName(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})

	assert.Equal(t, 0, testutil.CollectAndCount(m.mPanicRecovery), "panic_recovery")
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		// Uncaught:
		// Panic value:	panic: descriptor Desc{...} is invalid: "test.metric" is not a valid metric name [recovered]
		m.RegisterCounter("test.metric", []string{"this", "that"})
	}))
	assert.Equal(t, 1, testutil.CollectAndCount(m.mPanicRecovery))
}
