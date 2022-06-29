package apex

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	name := "test_counter"
	help := "created automagically by apex"
	labels := Labels{"this": "one"}

	c := NewCounters("", "", ':')
	metric, _ := c.Get(name, labels)

	err := c.Add(name, 5.0, labels)
	expected := BuildProm(name, help, "counter", labels, 5)
	assert.NoError(t, err)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)))

	err = c.Inc(name, labels)
	expected = BuildProm(name, help, "counter", labels, 6)
	assert.NoError(t, err)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)))

	err = c.Add(name, 6.0, labels)
	expected = BuildProm(name, help, "counter", labels, 12)
	assert.NoError(t, err)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)))

	err = c.Inc(name, labels)
	expected = BuildProm(name, help, "counter", labels, 13)
	assert.NoError(t, err)
	assert.NoError(t, testutil.CollectAndCompare(metric, strings.NewReader(expected)))
}
