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

const defaultCmpSelOpsTmpl = "pkg/sql/colexec/colexecsel/default_cmp_sel_ops_tmpl.go"

func genDefaultCmpSelOps(inputFileContents string, wr io.Writer) error {
	s := strings.ReplaceAll(inputFileContents, "_KIND", "{{.Kind}}")

	tmpl, err := template.New("default_cmp_sel_ops").Parse(s)
	if err != nil {
		return err
	}
	return tmpl.Execute(wr, []struct {
		HasConst bool
		Kind     string
	}{
		{HasConst: false, Kind: ""},
		{HasConst: true, Kind: "Const"},
	})
}

func init() {
	registerGenerator(genDefaultCmpSelOps, "default_cmp_sel_ops.eg.go", defaultCmpSelOpsTmpl)
}
