package apex

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInc(t *testing.T) {
	const metadata = `
		# HELP test_counter_inc created automagically by apex
		# TYPE test_counter_inc counter
	`
	m := New(MetricsOpts{MustRegister: true})
	m.Register(Counter, "test_counter_inc", []string{"this"})
	c := m.counters["test_counter_inc"]

	m.Inc("test_counter_inc", Labels{"this": "one"})
	expected := `
		test_counter_inc{this="one"} 1
	`
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(metadata+expected)), "test_counter_inc")

	m.Inc("test_counter_inc", Labels{"this": "one"})
	expected = `
		test_counter_inc{this="one"} 2
	`
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(metadata+expected)), "test_counter_inc")

}

func TestIncMismatchedLabels(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	m.Register(Counter, "test_counter_mismatch_inc", []string{"this", "that"})

	assert.Equal(t, 0, testutil.CollectAndCount(m.mPanicRecovery))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		// Uncaught:
		// Panic value:	inconsistent label cardinality: expected 2 label values but got 1 in prometheus.Labels{"name":"test.metrics"}
		m.Inc("test_counter_mismatch_inc", Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mPanicRecovery))
}

func TestIncv(t *testing.T) {
	const metadata = `
		# HELP test_counter_add created automagically by apex
		# TYPE test_counter_add counter
	`
	m := New(MetricsOpts{MustRegister: true})
	m.Register(Counter, "test_counter_add", []string{"this"})
	c := m.counters["test_counter_add"]

	m.Incv("test_counter_add", 2.0, Labels{"this": "one"})
	expected := `
		test_counter_add{this="one"} 2
	`
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(metadata+expected)), "test_counter_add")

	m.Incv("test_counter_add", 2.0, Labels{"this": "one"})
	expected = `
		test_counter_inc{this="one"} 2
	`
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(metadata+expected)), "test_counter_add")

}

func TestIncvMismatchedLabels(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	m.Register(Counter, "test_counter_mismatched_add", []string{"this", "that"})

	assert.Equal(t, 0, testutil.CollectAndCount(m.mPanicRecovery))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		// Uncaught:
		// Panic value:	inconsistent label cardinality: expected 2 label values but got 1 in prometheus.Labels{"name":"test.metrics"}
		m.Incv("test_counter_mismatched_add", 1.0, Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mPanicRecovery))
}
