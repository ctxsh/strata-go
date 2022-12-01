// Copyright (C) 2022, Rob Lyon <rob@ctxswitch.com>
//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package apex

import (
	"strconv"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

// CollectAndCompare is a helper function for testing.  It creates prometheus
// strings and compares them with the collector using the CollectAndCompare
// test utility.
func CollectAndCompare(
	t *testing.T,
	vec MetricVec,
	name string,
	mtype string,
	labels map[string]string,
	value float64,
) {
	assert.Equal(t, name, vec.Name())
	assert.Equal(t, MetricType(mtype), vec.Type())

	expected := createPromString(vec.Name(), vec.Type(), labels, value)
	assert.NoError(t, testutil.CollectAndCompare(vec.Vec(), strings.NewReader(expected)))
}

func createPromString(name string, ctype MetricType, labels map[string]string, value float64) string {
	var builder strings.Builder
	builder.WriteString("# HELP ")
	builder.WriteString(name + " ")
	builder.WriteString(" created automagically by apex\n")
	builder.WriteString("# TYPE ")
	builder.WriteString(name + " ")
	builder.WriteString(string(ctype) + "\n")

	switch ctype {
	case SummaryType:
		writeQuantiles(&builder, name, value, labels)
	case HistogramType:
		writeBuckets(&builder, name, value, labels)
	default:
		writeMetric(&builder, name, value, labels)
	}

	return builder.String()
}

func writeMetric(builder *strings.Builder, name string, value float64, labels map[string]string) {
	builder.WriteString(name)

	if labels != nil {
		builder.WriteString("{")
		for k, v := range labels {
			builder.WriteString(k + "=\"")
			builder.WriteString(v + "\"}")
		}
	}

	builder.WriteString(" ")
	val := strconv.FormatFloat(value, 'E', -1, 64)
	builder.WriteString(val + "\n")
}

// Come back to this later to make it more configurable based on summary options.
func writeQuantiles(builder *strings.Builder, name string, value float64, labels map[string]string) { //nolint:unparam
	val := strconv.FormatFloat(value, 'E', -1, 64)

	for _, q := range []string{"0.5", "0.9", "0.99"} {
		builder.WriteString(name)
		builder.WriteString("{quantile=\"")
		builder.WriteString(q)
		builder.WriteString("\"} ")
		builder.WriteString(val + "\n")
	}

	builder.WriteString(name)
	builder.WriteString("_sum")
	builder.WriteString(" ")
	builder.WriteString(val + "\n")

	builder.WriteString(name)
	builder.WriteString("_count 1\n")
}

// Come back to this later to make it more configurable based on histogram options.
func writeBuckets(builder *strings.Builder, name string, value float64, labels map[string]string) { //nolint:unparam
	val := strconv.FormatFloat(value, 'E', -1, 64)

	le := []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}

	buckets := make(map[float64]float64)
	for _, b := range le {
		if _, ok := buckets[b]; !ok {
			buckets[b] = 0
		}
		if value <= b {
			buckets[b]++
		}
	}

	for _, b := range le {
		builder.WriteString(name)
		builder.WriteString("{le=\"")
		q := strconv.FormatFloat(b, 'E', -1, 64)
		builder.WriteString(q)
		builder.WriteString("\"} ")
		v := strconv.FormatFloat(buckets[b], 'E', -1, 64)
		builder.WriteString(v + "\n")
	}

	builder.WriteString(name)
	builder.WriteString("{le=\"+Inf\"} 1\n")

	builder.WriteString(name)
	builder.WriteString("_sum")
	builder.WriteString(" ")
	builder.WriteString(val + "\n")

	builder.WriteString(name)
	builder.WriteString("_count 1\n")
}
