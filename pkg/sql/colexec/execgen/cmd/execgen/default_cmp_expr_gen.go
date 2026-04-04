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

type defaultCmpExprTmplInfo struct {
	NullableArgs bool
	FlippedArgs  bool
	Negate       bool
}

const defaultCmpExprTmpl = "pkg/sql/colexec/colexeccmp/default_cmp_expr_tmpl.go"

func genDefaultCmpExpr(inputFileContents string, wr io.Writer) error {
	s := strings.ReplaceAll(
		inputFileContents, "_EXPR_NAME", "cmp{{if .NullableArgs}}Nullable{{end}}"+
			"{{if .FlippedArgs}}Flipped{{end}}{{if .Negate}}Negate{{end}}ExprAdapter",
	)

	tmpl, err := template.New("default_cmp_expr").Parse(s)
	if err != nil {
		return err
	}
	var info []defaultCmpExprTmplInfo
	for _, nullable := range []bool{false, true} {
		for _, flipped := range []bool{false, true} {
			for _, negate := range []bool{false, true} {
				info = append(info, defaultCmpExprTmplInfo{
					NullableArgs: nullable,
					FlippedArgs:  flipped,
					Negate:       negate,
				})
			}
		}
	}
	return tmpl.Execute(wr, info)
}

func init() {
	registerGenerator(genDefaultCmpExpr, "default_cmp_expr.eg.go", defaultCmpExprTmpl)
}
