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

const vecTmpl = "pkg/col/coldata/vec_tmpl.go"

func genVec(inputFileContents string, wr io.Writer) error {
	r := strings.NewReplacer("_CANONICAL_TYPE_FAMILY", "{{.CanonicalTypeFamilyStr}}",
		"_TYPE_WIDTH", typeWidthReplacement,
		"_GOTYPESLICE", "{{.GoTypeSliceNameInColdata}}",
		"_GOTYPE", "{{.GoType}}",
		"_TYPE", "{{.VecMethod}}",
		"TemplateType", "{{.VecMethod}}",
	)
	s := r.Replace(inputFileContents)

	copyWithReorderedSource := makeFunctionRegex("_COPY_WITH_REORDERED_SOURCE", 1)
	s = copyWithReorderedSource.ReplaceAllString(s, `{{template "copyWithReorderedSource" buildDict "SrcHasNulls" $1}}`)

	s = replaceManipulationFuncs(s)

	// Now, generate the op, from the template.
	tmpl, err := template.New("vec_op").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	// It doesn't matter that we're passing in all overloads of Equality
	// comparison operator - we simply need to iterate over all supported
	// types.
	return tmpl.Execute(wr, sameTypeComparisonOpToOverloads[treecmp.EQ])
}
func init() {
	registerGenerator(genVec, "vec.eg.go", vecTmpl)
}
