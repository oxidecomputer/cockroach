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

package colconv

import (
	"strings"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
)

func init() {
	coldata.VecsToStringWithRowPrefix = vecsToStringWithRowPrefix
}

// vecsToStringWithRowPrefix returns a pretty representation of the vectors with
// each row being in a separate string.
func vecsToStringWithRowPrefix(vecs []coldata.Vec, length int, sel []int, prefix string) []string {
	var builder strings.Builder
	converter := NewAllVecToDatumConverter(len(vecs))
	defer converter.Release()
	converter.ConvertVecs(vecs, length, sel)
	result := make([]string, length)
	strs := make([]string, len(vecs))
	for i := 0; i < length; i++ {
		builder.Reset()
		rowIdx := i
		if sel != nil {
			rowIdx = sel[i]
		}
		builder.WriteString(prefix + "[")
		for colIdx := 0; colIdx < len(vecs); colIdx++ {
			strs[colIdx] = converter.GetDatumColumn(colIdx)[rowIdx].String()
		}
		builder.WriteString(strings.Join(strs, " "))
		builder.WriteString("]")
		result[i] = builder.String()
	}
	return result
}
