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
)

type nTileTmplInfo struct {
	HasPartition bool
	String       string
}

const nTileTmpl = "pkg/sql/colexec/colexecwindow/ntile_tmpl.go"

func genNTileOp(inputFileContents string, wr io.Writer) error {
	s := strings.ReplaceAll(inputFileContents, "_NTILE_STRING", "{{.String}}")

	// Now, generate the op, from the template.
	tmpl, err := template.New("ntile_op").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	nTileTmplInfos := []nTileTmplInfo{
		{HasPartition: false, String: "nTileNoPartition"},
		{HasPartition: true, String: "nTileWithPartition"},
	}
	return tmpl.Execute(wr, nTileTmplInfos)
}

func init() {
	registerGenerator(genNTileOp, "ntile.eg.go", nTileTmpl)
}
