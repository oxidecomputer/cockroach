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
)

const selectionOpsTmpl = "pkg/sql/colexec/colexecsel/selection_ops_tmpl.go"

func getSelectionOpsTmpl(inputFileContents string) (*template.Template, error) {
	r := strings.NewReplacer(
		"_LEFT_CANONICAL_TYPE_FAMILY", "{{.LeftCanonicalFamilyStr}}",
		"_LEFT_TYPE_WIDTH", typeWidthReplacement,
		"_RIGHT_CANONICAL_TYPE_FAMILY", "{{.RightCanonicalFamilyStr}}",
		"_RIGHT_TYPE_WIDTH", typeWidthReplacement,

		"_OP_CONST_NAME", "sel{{.Name}}{{.Left.VecMethod}}{{.Right.VecMethod}}ConstOp",
		"_OP_NAME", "sel{{.Name}}{{.Left.VecMethod}}{{.Right.VecMethod}}Op",
		"_NAME", "{{.Name}}",
		"_R_GO_TYPE", "{{.Right.GoType}}",
		"_L_TYP", "{{.Left.VecMethod}}",
		"_R_TYP", "{{.Right.VecMethod}}",
	)
	s := r.Replace(inputFileContents)

	assignCmpRe := makeFunctionRegex("_ASSIGN_CMP", 6)
	s = assignCmpRe.ReplaceAllString(s, makeTemplateFunctionCall("Right.Assign", 6))

	s = strings.ReplaceAll(s, "_HAS_NULLS", "$hasNulls")
	selConstLoop := makeFunctionRegex("_SEL_CONST_LOOP", 1)
	s = selConstLoop.ReplaceAllString(s, `{{template "selConstLoop" buildDict "Global" $ "HasNulls" $1 "Overload" .}}`)
	selLoop := makeFunctionRegex("_SEL_LOOP", 1)
	s = selLoop.ReplaceAllString(s, `{{template "selLoop" buildDict "Global" $ "HasNulls" $1 "Overload" .}}`)

	return template.New("selection_ops").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
}

func genSelectionOps(inputFileContents string, wr io.Writer) error {
	tmpl, err := getSelectionOpsTmpl(inputFileContents)
	if err != nil {
		return err
	}
	return tmpl.Execute(wr, twoArgsResolvedOverloadsInfo)
}

func init() {
	registerGenerator(genSelectionOps, "selection_ops.eg.go", selectionOpsTmpl)
}
