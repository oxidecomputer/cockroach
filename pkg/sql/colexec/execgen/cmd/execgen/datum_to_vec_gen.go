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

const datumToVecTmpl = "pkg/sql/colconv/datum_to_vec_tmpl.go"

func genDatumToVec(inputFileContents string, wr io.Writer) error {
	r := strings.NewReplacer(
		"_TYPE_FAMILY", "{{.TypeFamily}}",
		"_TYPE_WIDTH", typeWidthReplacement,
	)
	s := r.Replace(inputFileContents)

	preludeRe := makeFunctionRegex("_PRELUDE", 1)
	s = preludeRe.ReplaceAllString(s, makeTemplateFunctionCall("Prelude", 1))
	convertRe := makeFunctionRegex("_CONVERT", 1)
	s = convertRe.ReplaceAllString(s, makeTemplateFunctionCall("Convert", 1))

	tmpl, err := template.New("utils").Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, getRowToVecTmplInfos())
}

func init() {
	registerGenerator(genDatumToVec, "datum_to_vec.eg.go", datumToVecTmpl)
}
