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
)

const castTmpl = "pkg/sql/colexec/colexecbase/cast_tmpl.go"

const castOpInvocation = `{{template "castOp" buildDict "Global" . "FromInfo" $fromInfo "FromFamily" $fromFamily "ToFamily" $toFamily}}`

func genCastOperators(inputFileContents string, wr io.Writer) error {
	r := strings.NewReplacer(
		"_TYPE_FAMILY", "{{.TypeFamily}}",
		"_TYPE_WIDTH", typeWidthReplacement,
		"_TO_GO_TYPE", "{{.GoType}}",
		"_FROM_TYPE", "{{$fromInfo.VecMethod}}",
		"_TO_TYPE", "{{.VecMethod}}",
		"_NAME", "{{$fromInfo.TypeName}}{{.TypeName}}",
		"_GENERATE_CAST_OP", castOpInvocation,
	)
	s := r.Replace(inputFileContents)

	castRe := makeFunctionRegex("_CAST", 4)
	s = castRe.ReplaceAllString(s, makeTemplateFunctionCall("Cast", 4))

	tmpl, err := template.New("cast").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, getCastFromTmplInfos())
}

func init() {
	registerGenerator(genCastOperators, "cast.eg.go", castTmpl)
}
