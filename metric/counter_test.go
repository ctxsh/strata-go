package metric

import (
	"strings"
	"testing"

	"github.com/ctxswitch/apex/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	name := "test_counter"
	help := "created automagically by apex"
	labels := prometheus.Labels{"this": "one"}

	c := NewCounters("", "", ':')

	c.Add(name, 5.0, labels)
	expected := utils.BuildProm(name, help, "counter", labels, 5)
	assert.NoError(t, testutil.CollectAndCompare(c.Get(name, labels), strings.NewReader(expected)), "name")

	c.Inc(name, labels)
	expected = utils.BuildProm(name, help, "counter", labels, 6)
	assert.NoError(t, testutil.CollectAndCompare(c.Get(name, labels), strings.NewReader(expected)), "name")

	c.Add(name, 6.0, labels)
	expected = utils.BuildProm(name, help, "counter", labels, 12)
	assert.NoError(t, testutil.CollectAndCompare(c.Get(name, labels), strings.NewReader(expected)), "name")

	c.Inc(name, labels)
	expected = utils.BuildProm(name, help, "counter", labels, 13)
	assert.NoError(t, testutil.CollectAndCompare(c.Get(name, labels), strings.NewReader(expected)), "name")
}
