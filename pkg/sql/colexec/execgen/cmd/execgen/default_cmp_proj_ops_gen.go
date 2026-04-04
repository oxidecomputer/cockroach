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

const defaultCmpProjOpsTmpl = "pkg/sql/colexec/colexecproj/default_cmp_proj_ops_tmpl.go"

func genDefaultCmpProjOps(inputFileContents string, wr io.Writer) error {
	s := strings.ReplaceAll(inputFileContents, "_KIND", "{{.Kind}}")

	tmpl, err := template.New("default_cmp_proj_ops").Parse(s)
	if err != nil {
		return err
	}
	return tmpl.Execute(wr, []struct {
		// Comparison operators are always normalized so that the constant is
		// on the right side, so we skip generating the code when the constant
		// is on the left.
		IsRightConst bool
		Kind         string
	}{
		{IsRightConst: false, Kind: ""},
		{IsRightConst: true, Kind: "RConst"},
	})
}

func init() {
	registerGenerator(genDefaultCmpProjOps, "default_cmp_proj_ops.eg.go", defaultCmpProjOpsTmpl)
}
