package apex

import (
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

var labels = map[string]string{
	"region": "us-east-1",
}

func TestMetricsCounter(t *testing.T) {
	name := "test_total"
	m := testMetrics().WithLabels("region")

	m.CounterInc(name, "us-east-1")
	vec, err := getCounter(m, prefixedName(m.prefix, name, m.separator))
	assert.NoError(t, err)
	CollectAndCompare(t, vec, "apex_example_test_total", "counter", labels, 1.0)
	m.CounterAdd(name, 5.0, "us-east-1")
	CollectAndCompare(t, vec, "apex_example_test_total", "counter", labels, 6.0)

	m1 := m.WithPrefix("next")
	m1.CounterInc(name)
	vec, err = getCounter(m1, prefixedName(m1.prefix, name, m1.separator))
	assert.NoError(t, err)
	CollectAndCompare(t, vec, "apex_example_next_test_total", "counter", nil, 1.0)
	m1.CounterAdd(name, 5.0)
	CollectAndCompare(t, vec, "apex_example_next_test_total", "counter", nil, 6.0)
}

func TestMetricsGauge(t *testing.T) {
	name := "test_g"
	m := testMetrics().WithLabels("region")

	m.GaugeInc(name, "us-east-1")
	vec, err := getGauge(m, prefixedName(m.prefix, name, m.separator))
	assert.NoError(t, err)
	CollectAndCompare(t, vec, "apex_example_test_g", "gauge", labels, 1.0)
	m.GaugeAdd(name, 5.0, "us-east-1")
	CollectAndCompare(t, vec, "apex_example_test_g", "gauge", labels, 6.0)
	m.GaugeSet(name, 10.0, "us-east-1")
	CollectAndCompare(t, vec, "apex_example_test_g", "gauge", labels, 10.0)
	m.GaugeDec(name, "us-east-1")
	CollectAndCompare(t, vec, "apex_example_test_g", "gauge", labels, 9.0)
	m.GaugeSub(name, 9.0, "us-east-1")
	CollectAndCompare(t, vec, "apex_example_test_g", "gauge", labels, 0.0)

	m1 := m.WithPrefix("next")
	m1.GaugeInc(name)
	vec, err = getGauge(m, prefixedName(m1.prefix, name, m1.separator))
	assert.NoError(t, err)
	CollectAndCompare(t, vec, "apex_example_next_test_g", "gauge", nil, 1.0)
	m1.GaugeAdd(name, 5.0)
	CollectAndCompare(t, vec, "apex_example_next_test_g", "gauge", nil, 6.0)
	m1.GaugeSet(name, 10.0)
	CollectAndCompare(t, vec, "apex_example_next_test_g", "gauge", nil, 10.0)
	m1.GaugeDec(name)
	CollectAndCompare(t, vec, "apex_example_next_test_g", "gauge", nil, 9.0)
	m1.GaugeSub(name, 9.0)
	CollectAndCompare(t, vec, "apex_example_next_test_g", "gauge", nil, 0.0)
}

func getCounter(metrics *Metrics, n string) (MetricVec, error) {
	if v, ok := metrics.store.counters[n]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("missing counter")
}

func getGauge(metrics *Metrics, n string) (MetricVec, error) {
	if v, ok := metrics.store.gauges[n]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("missing gauge")
}

func testMetrics() *Metrics {
	registry := prometheus.NewPedanticRegistry()
	metrics := New(MetricsOpts{
		Separator:    '_',
		Registry:     registry,
		PanicOnError: true,
	}).WithPrefix("apex", "example")

	return metrics
}
