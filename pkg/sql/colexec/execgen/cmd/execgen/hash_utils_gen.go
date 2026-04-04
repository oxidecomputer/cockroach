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

const hashUtilsTmpl = "pkg/sql/colexec/colexechash/hash_utils_tmpl.go"

func genHashUtils(inputFileContents string, wr io.Writer) error {

	r := strings.NewReplacer(
		"_CANONICAL_TYPE_FAMILY", "{{.CanonicalTypeFamilyStr}}",
		"_TYPE_WIDTH", typeWidthReplacement,
		"_TYPE", "{{.VecMethod}}",
		"TemplateType", "{{.VecMethod}}",
	)
	s := r.Replace(inputFileContents)

	assignHash := makeFunctionRegex("_ASSIGN_HASH", 4)
	s = assignHash.ReplaceAllString(s, makeTemplateFunctionCall("Global.AssignHash", 4))

	rehash := makeFunctionRegex("_REHASH_BODY", 7)
	s = rehash.ReplaceAllString(s, `{{template "rehashBody" buildDict "Global" . "HasSel" $6 "HasNulls" $7}}`)

	s = replaceManipulationFuncsAmbiguous(".Global", s)

	tmpl, err := template.New("hash_utils").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, hashOverloads)
}

func init() {
	registerGenerator(genHashUtils, "hash_utils.eg.go", hashUtilsTmpl)
}
