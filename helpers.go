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
)

func BuildProm(name string, help string, ctype string, labels map[string]string, value float64) string {
	var builder strings.Builder
	builder.WriteString("# HELP ")
	builder.WriteString(name + " ")
	builder.WriteString(help + "\n")
	builder.WriteString("# TYPE ")
	builder.WriteString(name + " ")
	builder.WriteString(ctype + "\n")
	builder.WriteString(name + "{")
	for k, v := range labels {
		builder.WriteString(k + "=\"")
		builder.WriteString(v + "\"} ")
	}

	val := strconv.FormatFloat(value, 'E', -1, 64)
	builder.WriteString(val + "\n")

	return builder.String()
}
