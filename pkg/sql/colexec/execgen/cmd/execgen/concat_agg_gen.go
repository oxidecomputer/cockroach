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

	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

const concatAggTmpl = "pkg/sql/colexec/colexecagg/concat_agg_tmpl.go"

func genConcatAgg(inputFileContents string, wr io.Writer) error {
	accumulateConcatRe := makeFunctionRegex("_ACCUMULATE_CONCAT", 5)
	s := accumulateConcatRe.ReplaceAllString(inputFileContents, `{{template "accumulateConcat" buildDict "HasNulls" $4 "HasSel" $5}}`)

	s = replaceManipulationFuncs(s)

	tmpl, err := template.New("concat_agg").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}
	return tmpl.Execute(wr, aggTmplInfoBase{canonicalTypeFamily: types.BytesFamily})
}

func init() {
	registerAggGenerator(
		genConcatAgg, "concat_agg.eg.go", concatAggTmpl, true /* genWindowVariant */)
}
