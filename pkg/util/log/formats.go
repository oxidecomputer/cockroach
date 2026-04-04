// Copyright 2020 The Cockroach Authors.
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

package log

type logFormatter interface {
	formatterName() string
	// doc is used to generate the formatter documentation.
	doc() string
	// formatEntry formats a logEntry into a newly allocated *buffer.
	// The caller is responsible for calling putBuffer() afterwards.
	formatEntry(entry logEntry) *buffer

	// contentType is the MIME content-type field to use on
	// transports which use this metadata.
	contentType() string
}

var formatParsers = map[string]string{
	"crdb-v1":             "v1",
	"crdb-v1-count":       "v1",
	"crdb-v1-tty":         "v1",
	"crdb-v1-tty-count":   "v1",
	"crdb-v2":             "v2",
	"crdb-v2-tty":         "v2",
	"json":                "json",
	"json-compact":        "json-compact",
	"json-fluent":         "json",
	"json-fluent-compact": "json-compact",
}

var formatters = func() map[string]logFormatter {
	m := make(map[string]logFormatter)
	r := func(f logFormatter) {
		m[f.formatterName()] = f
	}
	r(formatCrdbV1{})
	r(formatCrdbV1WithCounter{})
	r(formatCrdbV1TTY{})
	r(formatCrdbV1TTYWithCounter{})
	r(formatCrdbV2{})
	r(formatCrdbV2TTY{})
	r(formatFluentJSONCompact{})
	r(formatFluentJSONFull{})
	r(formatJSONCompact{})
	r(formatJSONFull{})
	return m
}()

// GetFormatterDocs returns the embedded documentation for all the
// supported formats.
func GetFormatterDocs() map[string]string {
	m := make(map[string]string)
	for fmtName, f := range formatters {
		m[fmtName] = f.doc()
	}
	return m
}
