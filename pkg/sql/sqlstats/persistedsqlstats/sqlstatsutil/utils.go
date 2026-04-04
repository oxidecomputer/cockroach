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

package sqlstatsutil

import (
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/encoding"
)

// DatumToUint64 Convert a bytes datum to uint64.
func DatumToUint64(d tree.Datum) (uint64, error) {
	b := []byte(tree.MustBeDBytes(d))

	_, val, err := encoding.DecodeUint64Ascending(b)
	if err != nil {
		return 0, err
	}

	return val, nil
}
