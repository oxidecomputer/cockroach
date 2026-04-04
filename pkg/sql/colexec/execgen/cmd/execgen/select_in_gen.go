// Copyright 2019 The Cockroach Authors.
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

const selectInTmpl = "pkg/sql/colexec/select_in_tmpl.go"

func genSelectIn(inputFileContents string, wr io.Writer) error {
	r := strings.NewReplacer(
		"_CANONICAL_TYPE_FAMILY", "{{.CanonicalTypeFamilyStr}}",
		"_TYPE_WIDTH", typeWidthReplacement,
		"_GOTYPESLICE", "{{.GoTypeSliceName}}",
		"_GOTYPE", "{{.GoType}}",
		"_TYPE", "{{.VecMethod}}",
		"TemplateType", "{{.VecMethod}}",
	)
	s := r.Replace(inputFileContents)

	compare := makeFunctionRegex("_COMPARE", 5)
	s = compare.ReplaceAllString(s, makeTemplateFunctionCall("Compare", 5))

	s = replaceManipulationFuncs(s)

	tmpl, err := template.New("select_in").Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, sameTypeComparisonOpToOverloads[treecmp.EQ])
}

func init() {
	registerGenerator(genSelectIn, "select_in.eg.go", selectInTmpl)
}
