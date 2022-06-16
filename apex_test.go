package apex

import (
	"strings"
	"testing"

	"ctx.sh/apex/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMetricsCounter(t *testing.T) {
	name := "test_counter"
	help := "created automagically by apex"
	labels := Labels{"this": "one"}

	m := New(MetricsOpts{
		Namespace: "apex",
		Subsystem: "example",
		Separator: ':',
	})

	fullName, _ := utils.NameBuilder("apex", "example", name, ':')

	elem := m.counters
	metric, _ := elem.Get(name, prometheus.Labels(labels))

	m.CounterInc(name, labels)
	expected := utils.BuildProm(fullName, help, "counter", labels, 1)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)))

	m.CounterAdd(name, 5.0, labels)
	expected = utils.BuildProm(fullName, help, "counter", labels, 6)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)))

	m.CounterInc(name, labels)
	expected = utils.BuildProm(fullName, help, "counter", labels, 7)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)))

	m.CounterAdd(name, 6.0, labels)
	expected = utils.BuildProm(fullName, help, "counter", labels, 13)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)))
}

func TestMetricsGauge(t *testing.T) {
	name := "test_gauge"
	help := "created automagically by apex"
	labels := Labels{"this": "one"}

	m := New(MetricsOpts{
		Namespace: "apex",
		Subsystem: "example",
		Separator: ':',
	})

	fullName, _ := utils.NameBuilder("apex", "example", name, ':')

	elem := m.gauges
	metric, _ := elem.Get(name, prometheus.Labels(labels))

	m.GaugeSet(name, 100, labels)
	expected := utils.BuildProm(fullName, help, "gauge", labels, 100)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)), "name")

	m.GaugeInc(name, labels)
	expected = utils.BuildProm(fullName, help, "gauge", labels, 101)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)), "name")

	m.GaugeInc(name, labels)
	expected = utils.BuildProm(fullName, help, "gauge", labels, 102)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)), "name")

	m.GaugeDec(name, labels)
	expected = utils.BuildProm(fullName, help, "gauge", labels, 101)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)), "name")

	m.GaugeAdd(name, 9, labels)
	expected = utils.BuildProm(fullName, help, "gauge", labels, 110)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)), "name")

	m.GaugeSub(name, 109, labels)
	expected = utils.BuildProm(fullName, help, "gauge", labels, 1)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)), "name")
}
