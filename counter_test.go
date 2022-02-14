package apex

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCounterInc(t *testing.T) {
	name := "test_counter"
	help := "created automagically by apex"
	labels := Labels{"this": "one"}
	m := New(MetricsOpts{MustRegister: true})
	m.NewCounter(name, []string{"this"})
	c := m.getCounter(name)

	m.CounterInc(name, labels)
	expected := buildProm(t, name, help, "counter", labels, 1)
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(expected)), "name")
	m.CounterInc(name, labels)
	expected = buildProm(t, name, help, "counter", labels, 2)
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(expected)), "name")
}

func TestCounterAdd(t *testing.T) {
	name := "test_counter"
	help := "created automagically by apex"
	labels := Labels{"this": "one"}
	m := New(MetricsOpts{MustRegister: true})
	m.NewCounter(name, []string{"this"})
	c := m.getCounter(name)

	m.CounterAdd(name, 5, labels)
	expected := buildProm(t, name, help, "counter", labels, 5)
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(expected)), "name")
	m.CounterAdd(name, 6, labels)
	expected = buildProm(t, name, help, "counter", labels, 11)
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(expected)), "name")
}

func TestCounterIncMismatchedLabels(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	m.NewCounter("test_counter_mismatch_inc", []string{"this", "that"})

	assert.Equal(t, 0, testutil.CollectAndCount(m.mPanicRecovery))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		// Uncaught:
		// Panic value:	inconsistent label cardinality: expected 2 label values but got 1 in prometheus.Labels{"name":"test.metrics"}
		m.CounterInc("test_counter_mismatch", Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mPanicRecovery))
}

func TestCounterAddMismatchedLabels(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	m.NewCounter("test_counter_mismatched_add", []string{"this", "that"})

	assert.Equal(t, 0, testutil.CollectAndCount(m.mPanicRecovery))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		// Uncaught:
		// Panic value:	inconsistent label cardinality: expected 2 label values but got 1 in prometheus.Labels{"name":"test.metrics"}
		m.CounterAdd("test_counter_mismatched_add", 1.0, Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mPanicRecovery))
}

func TestCounterIncInvalidCounter(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	assert.Equal(t, 0, testutil.CollectAndCount(m.mInvalidCounter))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.CounterInc("were_is_waldo", Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mInvalidCounter))
}

func TestCounterAddInvalidCounter(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	assert.Equal(t, 0, testutil.CollectAndCount(m.mInvalidCounter))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.CounterAdd("were_is_waldo", 1.0, Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mInvalidCounter))
}
