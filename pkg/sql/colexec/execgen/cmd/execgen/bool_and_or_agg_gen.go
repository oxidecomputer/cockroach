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

	"github.com/cockroachdb/cockroach/pkg/sql/colexecerror"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/errors"
)

type booleanAggTmplInfo struct {
	aggTmplInfoBase
	IsAnd bool
}

func (b booleanAggTmplInfo) AssignBoolOp(target, l, r string) string {
	switch b.IsAnd {
	case true:
		return fmt.Sprintf("%s = %s && %s", target, l, r)
	case false:
		return fmt.Sprintf("%s = %s || %s", target, l, r)
	default:
		colexecerror.InternalError(errors.AssertionFailedf("unsupported boolean agg type"))
		// This code is unreachable, but the compiler cannot infer that.
		return ""
	}
}

func (b booleanAggTmplInfo) OpType() string {
	if b.IsAnd {
		return "And"
	}
	return "Or"
}

func (b booleanAggTmplInfo) DefaultVal() string {
	if b.IsAnd {
		return "true"
	}
	return "false"
}

// Avoid unused warnings. These methods are used in the template.
var (
	_ = booleanAggTmplInfo{}.AssignBoolOp
	_ = booleanAggTmplInfo{}.OpType
	_ = booleanAggTmplInfo{}.DefaultVal
)

const boolAggTmpl = "pkg/sql/colexec/colexecagg/bool_and_or_agg_tmpl.go"

func genBooleanAgg(inputFileContents string, wr io.Writer) error {
	r := strings.NewReplacer(
		"_OP_TYPE", "{{.OpType}}",
		"_DEFAULT_VAL", "{{.DefaultVal}}",
	)
	s := r.Replace(inputFileContents)

	accumulateBoolean := makeFunctionRegex("_ACCUMULATE_BOOLEAN", 5)
	s = accumulateBoolean.ReplaceAllString(s, `{{template "accumulateBoolean" buildDict "Global" . "HasNulls" $4 "HasSel" $5}}`)

	assignBoolRe := makeFunctionRegex("_ASSIGN_BOOL_OP", 3)
	s = assignBoolRe.ReplaceAllString(s, makeTemplateFunctionCall(`AssignBoolOp`, 3))

	s = replaceManipulationFuncs(s)

	tmpl, err := template.New("bool_and_or_agg").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, []booleanAggTmplInfo{
		{
			aggTmplInfoBase: aggTmplInfoBase{canonicalTypeFamily: types.BoolFamily},
			IsAnd:           true,
		},
		{
			aggTmplInfoBase: aggTmplInfoBase{canonicalTypeFamily: types.BoolFamily},
			IsAnd:           false,
		},
	})
}

func init() {
	registerAggGenerator(
		genBooleanAgg, "bool_and_or_agg.eg.go", boolAggTmpl, true /* genWindowVariant */)
}
