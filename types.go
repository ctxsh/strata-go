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

package strata

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricType string

const (
	// CounterType represents an strata wrapper around the prometheus CounterVec
	// type.
	CounterType MetricType = "counter"
	// GaugeType represents an strata wrapper around the prometheus GaugeVec
	// type.
	GaugeType MetricType = "gauge"
	// SummaryType represents an strata wrapper around the prometheus SummaryVec
	// type.
	SummaryType MetricType = "summary"
	// HistogramType represents an strata wrapper around the prometheus HistogramVec
	// type.
	HistogramType MetricType = "histogram"
	// Defines the metrics help string.  This is currently not settable.
	DefaultHelpString string = "created automagically by strata"
)

// MetricVec defines the interface for strata metrics collector wrappers.
type MetricVec interface {
	Name() string
	Type() MetricType
	Vec() prometheus.Collector
}
