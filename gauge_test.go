package apex

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGauge(t *testing.T) {
	name := "test_metric"
	help := "created automagically by apex"
	labels := Labels{"this": "one"}
	m := New(MetricsOpts{MustRegister: true})
	m.NewGauge(name, []string{"this"})
	c := m.getGauge(name)

	m.GaugeSet(name, 100, labels)
	expected := buildProm(t, name, help, "gauge", labels, 100)
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(expected)), "name")
	m.GaugeInc(name, labels)
	expected = buildProm(t, name, help, "gauge", labels, 101)
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(expected)), "name")
	m.GaugeInc(name, labels)
	expected = buildProm(t, name, help, "gauge", labels, 102)
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(expected)), "name")
	m.GaugeDec(name, labels)
	expected = buildProm(t, name, help, "gauge", labels, 101)
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(expected)), "name")
	m.GaugeAdd(name, 9, labels)
	expected = buildProm(t, name, help, "gauge", labels, 110)
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(expected)), "name")
	m.GaugeSub(name, 109, labels)
	expected = buildProm(t, name, help, "gauge", labels, 1)
	assert.NoError(t, testutil.CollectAndCompare(c, strings.NewReader(expected)), "name")
}

func TestGaugeSetMismatchedLabels(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	m.NewGauge("test_metric_set", []string{"this", "that"})

	assert.Equal(t, 0, testutil.CollectAndCount(m.mPanicRecovery))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.GaugeSet("test_metric_set", 100, Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mPanicRecovery))
}

func TestGaugeIncMismatchedLabels(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	m.NewGauge("test_metric_inc", []string{"this", "that"})

	assert.Equal(t, 0, testutil.CollectAndCount(m.mPanicRecovery))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.GaugeInc("test_metric_inc", Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mPanicRecovery))
}

func TestGaugeDecMismatchedLabels(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	m.NewGauge("test_metric_dec", []string{"this", "that"})

	assert.Equal(t, 0, testutil.CollectAndCount(m.mPanicRecovery))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.GaugeDec("test_metric_dec", Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mPanicRecovery))
}

func TestGaugeAddMismatchedLabels(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	m.NewGauge("test_metric_add", []string{"this", "that"})

	assert.Equal(t, 0, testutil.CollectAndCount(m.mPanicRecovery))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.GaugeAdd("test_metric_add", 100, Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mPanicRecovery))
}

func TestGaugeSubMismatchedLabels(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	m.NewGauge("test_metric_sub", []string{"this", "that"})

	assert.Equal(t, 0, testutil.CollectAndCount(m.mPanicRecovery))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.GaugeSub("test_metric_sub", 100, Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mPanicRecovery))
}

func TestGaugeSetInvalidCounter(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	assert.Equal(t, 0, testutil.CollectAndCount(m.mInvalidCounter))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.GaugeSet("were_is_waldo", 100, Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mInvalidGauge))
}

func TestGaugeIncInvalidCounter(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	assert.Equal(t, 0, testutil.CollectAndCount(m.mInvalidGauge))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.GaugeInc("were_is_waldo", Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mInvalidGauge))
}

func TestGaugeDecInvalidCounter(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	assert.Equal(t, 0, testutil.CollectAndCount(m.mInvalidGauge))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.GaugeDec("were_is_waldo", Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mInvalidGauge))
}

func TestGaugeAddInvalidCounter(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	assert.Equal(t, 0, testutil.CollectAndCount(m.mInvalidGauge))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.GaugeAdd("were_is_waldo", 100, Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mInvalidGauge))
}

func TestGaugeSubInvalidCounter(t *testing.T) {
	m := New(MetricsOpts{MustRegister: true})
	assert.Equal(t, 0, testutil.CollectAndCount(m.mInvalidGauge))
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		m.GaugeSub("were_is_waldo", 100, Labels{"this": "one"})
	}))

	assert.Equal(t, 1, testutil.CollectAndCount(m.mInvalidGauge))
}
