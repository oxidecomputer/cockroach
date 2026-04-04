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

type rankTmplInfo struct {
	Dense        bool
	HasPartition bool
	String       string
}

// UpdateRank is used to encompass the different logic between DENSE_RANK and
// RANK window functions when updating the internal state of rank operators. It
// is used by the template.
func (r rankTmplInfo) UpdateRank() string {
	if r.Dense {
		return `r.rank++`
	}
	return `r.rank += r.rankIncrement
r.rankIncrement = 1`
}

// UpdateRankIncrement is used to encompass the different logic between
// DENSE_RANK and RANK window functions when updating the internal state of
// rank operators. It is used by the template.
func (r rankTmplInfo) UpdateRankIncrement() string {
	if r.Dense {
		return ``
	}
	return `r.rankIncrement++`
}

// Avoid unused warnings. These methods are used in the template.
var (
	_ = rankTmplInfo{}.UpdateRank()
	_ = rankTmplInfo{}.UpdateRankIncrement()
)

const rankTmpl = "pkg/sql/colexec/colexecwindow/rank_tmpl.go"

func genRankOps(inputFileContents string, wr io.Writer) error {
	s := strings.ReplaceAll(inputFileContents, "_RANK_STRING", "{{.String}}")

	computeRankRe := makeFunctionRegex("_COMPUTE_RANK", 1)
	s = computeRankRe.ReplaceAllString(s, `{{template "computeRank" buildDict "Global" . "HasPartition" .HasPartition "HasSel" $1}}`)
	updateRankRe := makeFunctionRegex("_UPDATE_RANK", 0)
	s = updateRankRe.ReplaceAllString(s, makeTemplateFunctionCall("Global.UpdateRank", 0))
	updateRankIncrementRe := makeFunctionRegex("_UPDATE_RANK_INCREMENT", 0)
	s = updateRankIncrementRe.ReplaceAllString(s, makeTemplateFunctionCall("Global.UpdateRankIncrement", 0))

	// Now, generate the op, from the template.
	tmpl, err := template.New("rank_op").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	rankTmplInfos := []rankTmplInfo{
		{Dense: false, HasPartition: false, String: "rankNoPartition"},
		{Dense: false, HasPartition: true, String: "rankWithPartition"},
		{Dense: true, HasPartition: false, String: "denseRankNoPartition"},
		{Dense: true, HasPartition: true, String: "denseRankWithPartition"},
	}
	return tmpl.Execute(wr, rankTmplInfos)
}

func init() {
	registerGenerator(genRankOps, "rank.eg.go", rankTmpl)
}
