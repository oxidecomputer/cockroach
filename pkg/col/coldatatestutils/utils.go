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

package coldatatestutils

import (
	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

// CopyBatch copies the original batch and returns that copy. However, note that
// the underlying capacity might be different (a new batch is created only with
// capacity original.Length()).
// NOTE: memory accounting is not performed.
func CopyBatch(
	original coldata.Batch, typs []*types.T, factory coldata.ColumnFactory,
) coldata.Batch {
	b := coldata.NewMemBatchWithCapacity(typs, original.Length(), factory)
	b.SetLength(original.Length())
	for colIdx, col := range original.ColVecs() {
		b.ColVec(colIdx).Copy(coldata.SliceArgs{
			Src:       col,
			SrcEndIdx: original.Length(),
		})
	}
	return b
}
