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
	ErrInvalidMetricName  = ApexError("Invalid metric name")
	ErrRegistrationFailed = ApexError("Unable to register collector")
	ErrAlreadyRegistered  = ApexError("metric is already registered")
)

func (e ApexError) Error() string {
	return string(e)
}

type ApexInternalErrorMetrics struct {
	errPanicRecovery      *prometheus.CounterVec
	errInvalidMetricName  *prometheus.CounterVec
	errRegistrationFailed *prometheus.CounterVec
	errAlreadyRegistered  *prometheus.CounterVec
}

func NewApexInternalErrorMetrics(ns string, sub string, sep rune) *ApexInternalErrorMetrics {
	var builder strings.Builder

	if ns != "" {
		builder.WriteString(ns)
		builder.WriteRune(sep)
	}

	if sub != "" {
		builder.WriteString(sub)
		builder.WriteRune(sep)
	}

	builder.WriteString("apex")
	builder.WriteRune(sep)
	builder.WriteString("error")
	builder.WriteRune(sep)

	prefix := builder.String()

	errPanicRecovery := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prefix + "panic_recovery",
	}, []string{"name", "type"})

	errInvalidMetricName := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prefix + "invalid_metric_name",
	}, []string{"name", "type"})

	errRegistrationFailed := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prefix + "registration_failed",
	}, []string{"name", "type"})

	errAlreadyRegistered := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: prefix + "already_registered",
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

func (a *ApexInternalErrorMetrics) PanicRecovery(name string, t string) {
	a.errPanicRecovery.With(prometheus.Labels{
		"name": name,
		"type": t,
	}).Inc()
}

func (a *ApexInternalErrorMetrics) InvalidMetricName(name string, t string) {
	a.errInvalidMetricName.With(prometheus.Labels{
		"name": name,
		"type": t,
	}).Inc()
}

func (a *ApexInternalErrorMetrics) RegistrationFailed(name string, t string) {
	a.errRegistrationFailed.With(prometheus.Labels{
		"name": name,
		"type": t,
	}).Inc()
}

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
