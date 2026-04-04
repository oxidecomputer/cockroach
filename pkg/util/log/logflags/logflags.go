// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package logflags

import "flag"

// LogToStderrName and others are flag names.
const (
	VModuleName       = "vmodule"
	ShowLogsName      = "show-logs"
	TestLogConfigName = "test-log-config"
)

// InitFlags creates logging flags which update the given variables. The passed mutex is
// locked while the boolean variables are accessed during flag updates.
func InitFlags(showLogs *bool, testLogConfig *string, vmodule flag.Value) {
	flag.Var(vmodule, VModuleName, "comma-separated list of pattern=N settings for file-filtered logging (significantly hurts performance)")
	flag.StringVar(testLogConfig, TestLogConfigName, "", "YAML log configuration for tests")
	flag.BoolVar(showLogs, ShowLogsName, *showLogs, "print logs instead of saving them in files")
}
