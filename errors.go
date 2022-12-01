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
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type ApexError string

const (
	// ErrInvalidMetricName is returned when a metric name contains other
	// characters other than [a-zA-Z_-].
	ErrInvalidMetricName = ApexError("Invalid metric name")
	// ErrRegistrationFailed is returned if prometheus is unable to register
	// the collector.
	ErrRegistrationFailed = ApexError("Unable to register collector")
	// ErrAlreadyRegistered is returned if prometheus has already registered
	// a collector.
	ErrAlreadyRegistered = ApexError("metric is already registered")
)

// Error implements the error interface for ApexError
func (e ApexError) Error() string {
	return string(e)
}

// ApexInternalErrorMetrics provides internal counters for recovered
// errors from the prometheus collector when PanicOnError is false.
type ApexInternalErrorMetrics struct {
	errPanicRecovery      *prometheus.CounterVec
	errInvalidMetricName  *prometheus.CounterVec
	errRegistrationFailed *prometheus.CounterVec
	errAlreadyRegistered  *prometheus.CounterVec
}

// NewApexInternalErrorMetrics defines and registers the internal collectors and
// returns a new ApexInternalErrorMetrics struct.
func NewApexInternalErrorMetrics(prefixes []string, sep rune) *ApexInternalErrorMetrics {
	prefix := strings.Join(prefixes, "_")

	errPanicRecovery := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prefix + "_panic_recovery",
	}, []string{"name", "type"})

	errInvalidMetricName := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prefix + "_invalid_metric_name",
	}, []string{"name", "type"})

	errRegistrationFailed := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prefix + "_registration_failed",
	}, []string{"name", "type"})

	errAlreadyRegistered := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prefix + "_already_registered",
	}, []string{"name", "type"})

	_ = register(errPanicRecovery)
	_ = register(errInvalidMetricName)
	_ = register(errRegistrationFailed)
	_ = register(errAlreadyRegistered)

	return &ApexInternalErrorMetrics{
		errPanicRecovery:      errPanicRecovery,
		errInvalidMetricName:  errInvalidMetricName,
		errRegistrationFailed: errRegistrationFailed,
		errAlreadyRegistered:  errAlreadyRegistered,
	}
}

// PanicRecovery provides a helper function for incrementing the errPanicRecovery
// collector.
func (a *ApexInternalErrorMetrics) PanicRecovery(name string, t string) {
	a.errPanicRecovery.With(prometheus.Labels{
		"name": name,
		"type": t,
	}).Inc()
}

// InvalidMeticName provides a helper function for incrementing the errInvalidMetricName
// collector.
func (a *ApexInternalErrorMetrics) InvalidMetricName(name string, t string) {
	a.errInvalidMetricName.With(prometheus.Labels{
		"name": name,
		"type": t,
	}).Inc()
}

// RegistrationFailed provides a helper function for incrementing the errRegistrationFailed
// collector.
func (a *ApexInternalErrorMetrics) RegistrationFailed(name string, t string) {
	a.errRegistrationFailed.With(prometheus.Labels{
		"name": name,
		"type": t,
	}).Inc()
}

// AlreadyRegistered provides a helper function for incrementing the errAlreadyRegistered
// collector.
func (a *ApexInternalErrorMetrics) AlreadyRegistered(name string, t string) {
	a.errAlreadyRegistered.With(prometheus.Labels{
		"name": name,
		"type": t,
	}).Inc()
}

func register(metric prometheus.Collector) error {
	if err := prometheus.Register(metric); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}
	return nil
}
