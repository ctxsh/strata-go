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

func NameBuilder(ns string, sub interface{}, name string, sep rune) (string, error) {
	var builder strings.Builder

	if ns != "" {
		builder.WriteString(ns)
		builder.WriteRune(sep)
	}

	ss := subSystemToString(sub, sep)
	if ss != "" {
		builder.WriteString(ss)
		builder.WriteRune(sep)
	}

	if name == "" {
		return "", ErrInvalidMetricName
	}
	builder.WriteString(name)
	return builder.String(), nil
}

func Register(metric prometheus.Collector) error {
	if err := prometheus.Register(metric); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return ErrAlreadyRegistered
		} else {
			return ErrRegistrationFailed
		}
	}
	return nil
}

func subSystemToString(sub interface{}, sep rune) string {
	switch s := sub.(type) {
	case string:
		return s
	case []string:
		return strings.Join(s, string(sep))
	}

	panic("unsupported subsystem type")
}
