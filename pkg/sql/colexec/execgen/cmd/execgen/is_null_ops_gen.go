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

package main

import (
	"io"
	"strings"
	"text/template"
)

const isNullOpsTmpl = "pkg/sql/colexec/is_null_ops_tmpl.go"

func genIsNullOps(inputFileContents string, wr io.Writer) error {
	r := strings.NewReplacer(
		"_IS_TUPLE", ".IsTuple",
		"_KIND", "{{.Kind}}",
	)
	s := r.Replace(inputFileContents)

	computeIsNullRe := makeFunctionRegex("_COMPUTE_IS_NULL", 5)
	s = computeIsNullRe.ReplaceAllString(s, `{{template "computeIsNull" buildDict "HasNulls" $4 "IsTuple" $5}}`)
	maybeSelectRe := makeFunctionRegex("_MAYBE_SELECT", 6)
	s = maybeSelectRe.ReplaceAllString(s, `{{template "maybeSelect" buildDict "IsTuple" $6}}`)

	tmpl, err := template.New("is_null_ops").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, []struct {
		IsTuple bool
		Kind    string
	}{
		{IsTuple: false, Kind: ""},
		{IsTuple: true, Kind: "Tuple"},
	})
}

func init() {
	registerGenerator(genIsNullOps, "is_null_ops.eg.go", isNullOpsTmpl)
}
