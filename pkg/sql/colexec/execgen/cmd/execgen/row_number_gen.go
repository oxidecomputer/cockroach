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

type rowNumberTmplInfo struct {
	HasPartition bool
	String       string
}

const rowNumberTmpl = "pkg/sql/colexec/colexecwindow/row_number_tmpl.go"

func genRowNumberOp(inputFileContents string, wr io.Writer) error {
	s := strings.ReplaceAll(inputFileContents, "_ROW_NUMBER_STRING", "{{.String}}")

	computeRowNumberRe := makeFunctionRegex("_COMPUTE_ROW_NUMBER", 1)
	s = computeRowNumberRe.ReplaceAllString(s, `{{template "computeRowNumber" buildDict "HasPartition" .HasPartition "HasSel" $1}}`)

	// Now, generate the op, from the template.
	tmpl, err := template.New("row_number_op").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	rowNumberTmplInfos := []rowNumberTmplInfo{
		{HasPartition: false, String: "rowNumberNoPartition"},
		{HasPartition: true, String: "rowNumberWithPartition"},
	}
	return tmpl.Execute(wr, rowNumberTmplInfos)
}

func init() {
	registerGenerator(genRowNumberOp, "row_number.eg.go", rowNumberTmpl)
}
