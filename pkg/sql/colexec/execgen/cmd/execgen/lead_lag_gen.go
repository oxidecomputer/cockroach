// Copyright 2021 The Cockroach Authors.
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

const leadLagTmpl = "pkg/sql/colexec/colexecwindow/lead_lag_tmpl.go"

func replaceLeadLagTmplVariables(tmpl string) string {
	r := strings.NewReplacer(
		"_CANONICAL_TYPE_FAMILY", "{{.CanonicalTypeFamilyStr}}",
		"_TYPE_WIDTH", typeWidthReplacement,
		"_GOTYPESLICE", "{{.GoTypeSliceName}}",
		"_TYPE", "{{.VecMethod}}",
		"TemplateType", "{{.VecMethod}}",
	)
	s := r.Replace(tmpl)

	processBatch := makeFunctionRegex("_PROCESS_BATCH", 2)
	s = processBatch.ReplaceAllString(s, `{{template "processBatchTmpl" buildDict "VecMethod" .VecMethod "IsBytesLike" .IsBytesLike "OffsetHasNulls" $1 "DefaultHasNulls" $2}}`)

	return replaceManipulationFuncs(s)
}

func genLeadOp(inputFileContents string, wr io.Writer) error {
	s := replaceLeadLagTmplVariables(inputFileContents)
	r := strings.NewReplacer(
		"_OP_NAME", "lead",
		"_UPPERCASE_NAME", "Lead",
	)
	s = r.Replace(s)

	// Now, generate the op, from the template.
	tmpl, err := template.New("lead_ops").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, sameTypeComparisonOpToOverloads[treecmp.EQ])
}

func genLagOp(inputFileContents string, wr io.Writer) error {
	s := replaceLeadLagTmplVariables(inputFileContents)
	r := strings.NewReplacer(
		"_OP_NAME", "lag",
		"_UPPERCASE_NAME", "Lag",
	)
	s = r.Replace(s)

	// Now, generate the op, from the template.
	tmpl, err := template.New("lag_ops").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, sameTypeComparisonOpToOverloads[treecmp.EQ])
}

func init() {
	registerGenerator(genLeadOp, "lead.eg.go", leadLagTmpl)
	registerGenerator(genLagOp, "lag.eg.go", leadLagTmpl)
}
