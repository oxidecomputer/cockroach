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
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

const substringTmpl = "pkg/sql/colexec/substring_tmpl.go"

func genSubstring(inputFileContents string, wr io.Writer) error {
	r := strings.NewReplacer(
		"_START_WIDTH", fmt.Sprintf("{{$startWidth}}{{if eq $startWidth %d}}: default{{end}}", anyWidth),
		"_LENGTH_WIDTH", fmt.Sprintf("{{$lengthWidth}}{{if eq $lengthWidth %d}}: default{{end}}", anyWidth),
		"_StartType", fmt.Sprintf("Int{{if eq $startWidth %d}}64{{else}}{{$startWidth}}{{end}}", anyWidth),
		"_LengthType", fmt.Sprintf("Int{{if eq $lengthWidth %d}}64{{else}}{{$lengthWidth}}{{end}}", anyWidth),
	)
	s := r.Replace(inputFileContents)

	tmpl, err := template.New("substring").Parse(s)
	if err != nil {
		return err
	}

	supportedIntWidths := supportedWidthsByCanonicalTypeFamily[types.IntFamily]
	intWidthsToIntWidths := make(map[int32][]int32)
	for _, intWidth := range supportedIntWidths {
		intWidthsToIntWidths[intWidth] = supportedIntWidths
	}
	return tmpl.Execute(wr, intWidthsToIntWidths)
}

func init() {
	registerGenerator(genSubstring, "substring.eg.go", substringTmpl)
}
