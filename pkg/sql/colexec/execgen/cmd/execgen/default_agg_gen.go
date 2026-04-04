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
	"text/template"
)

const defaultAggTmpl = "pkg/sql/colexec/colexecagg/default_agg_tmpl.go"

func genDefaultAgg(inputFileContents string, wr io.Writer) error {
	addTuple := makeFunctionRegex("_ADD_TUPLE", 5)
	s := addTuple.ReplaceAllString(inputFileContents, `{{template "addTuple" buildDict "HasSel" $5}}`)

	setResult := makeFunctionRegex("_SET_RESULT", 2)
	s = setResult.ReplaceAllString(s, `{{template "setResult"}}`)

	tmpl, err := template.New("default_agg").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, struct{}{})

}

func init() {
	registerAggGenerator(
		genDefaultAgg, "default_agg.eg.go", defaultAggTmpl, false /* genWindowVariant */)
}
