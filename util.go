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

// Register registers a collector with prometheus.
func Register(reg prometheus.Registerer, metric prometheus.Collector) error {
	if err := reg.Register(metric); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return ErrAlreadyRegistered
		} else {
			return err
		}
	}
	return nil
}

// SlicePairsToMap copies key value pairs to a map.
func SlicePairsToMap(pairs []string) map[string]string {
	// TODO: this is too brittle. Fix me.
	m := make(map[string]string)
	for i := 0; i < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func prefixedName(prefix, name string, sep rune) string {
	if prefix == "" {
		return name
	}
	return prefix + string(sep) + name
}
