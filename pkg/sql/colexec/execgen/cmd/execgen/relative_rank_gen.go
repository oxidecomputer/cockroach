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

type relativeRankTmplInfo struct {
	IsPercentRank bool
	IsCumeDist    bool
	HasPartition  bool
	String        string
}

const relativeRankTmpl = "pkg/sql/colexec/colexecwindow/relative_rank_tmpl.go"

func genRelativeRankOps(inputFileContents string, wr io.Writer) error {
	s := strings.ReplaceAll(inputFileContents, "_RELATIVE_RANK_STRING", "{{.String}}")

	computePartitionsSizesRe := makeFunctionRegex("_COMPUTE_PARTITIONS_SIZES", 1)
	s = computePartitionsSizesRe.ReplaceAllString(s, `{{template "computePartitionsSizes" buildDict "HasSel" $1}}`)
	computePeerGroupsSizesRe := makeFunctionRegex("_COMPUTE_PEER_GROUPS_SIZES", 1)
	s = computePeerGroupsSizesRe.ReplaceAllString(s, `{{template "computePeerGroupsSizes" buildDict "HasSel" $1}}`)

	// Now, generate the op, from the template.
	tmpl, err := template.New("relative_rank_op").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	relativeRankTmplInfos := []relativeRankTmplInfo{
		{IsPercentRank: true, HasPartition: false, String: "percentRankNoPartition"},
		{IsPercentRank: true, HasPartition: true, String: "percentRankWithPartition"},
		{IsCumeDist: true, HasPartition: false, String: "cumeDistNoPartition"},
		{IsCumeDist: true, HasPartition: true, String: "cumeDistWithPartition"},
	}
	return tmpl.Execute(wr, relativeRankTmplInfos)
}

func init() {
	registerGenerator(genRelativeRankOps, "relative_rank.eg.go", relativeRankTmpl)
}
