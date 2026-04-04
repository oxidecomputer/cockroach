// Copyright 2018 The Cockroach Authors.
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

	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree/treecmp"
)

const minMaxAggTmpl = "pkg/sql/colexec/colexecagg/min_max_agg_tmpl.go"

func genMinMaxAgg(inputFileContents string, wr io.Writer) error {
	r := strings.NewReplacer(
		"_CANONICAL_TYPE_FAMILY", "{{.CanonicalTypeFamilyStr}}",
		"_TYPE_WIDTH", typeWidthReplacement,
		"_AGG_TITLE", "{{.AggTitle}}",
		"_AGG", "{{$agg}}",
		"_GOTYPESLICE", "{{.GoTypeSliceName}}",
		"_GOTYPE", "{{.GoType}}",
		"_TYPE", "{{.VecMethod}}",
		"TemplateType", "{{.VecMethod}}",
	)
	s := r.Replace(inputFileContents)

	assignCmpRe := makeFunctionRegex("_ASSIGN_CMP", 6)
	s = assignCmpRe.ReplaceAllString(s, makeTemplateFunctionCall("Assign", 6))

	accumulateMinMax := makeFunctionRegex("_ACCUMULATE_MINMAX", 5)
	s = accumulateMinMax.ReplaceAllString(s, `{{template "accumulateMinMax" buildDict "Global" . "HasNulls" $4 "HasSel" $5}}`)

	s = replaceManipulationFuncs(s)

	tmpl, err := template.New("min_max_agg").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}
	return tmpl.Execute(wr, []struct {
		Agg       string
		AggTitle  string
		Overloads []*oneArgOverload
	}{
		{
			Agg:       "min",
			AggTitle:  "Min",
			Overloads: sameTypeComparisonOpToOverloads[treecmp.LT],
		},
		{
			Agg:       "max",
			AggTitle:  "Max",
			Overloads: sameTypeComparisonOpToOverloads[treecmp.GT],
		},
	})
}

func init() {
	registerAggGenerator(
		genMinMaxAgg, "min_max_agg.eg.go", minMaxAggTmpl, true /* genWindowVariant */)
}
