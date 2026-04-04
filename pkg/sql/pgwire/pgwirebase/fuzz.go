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

//go:build gofuzz
// +build gofuzz

package pgwirebase

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

var (
	// Compile a slice of all typs.
	typs = func() []*types.T {
		var ret []*types.T
		for _, typ := range types.OidToType {
			ret = append(ret, typ)
		}
		return ret
	}()
)

func FuzzDecodeDatum(data []byte) int {
	if len(data) < 2 {
		return 0
	}

	typ := typs[int(data[1])%len(typs)]
	code := FormatCode(data[0]) % (FormatBinary + 1)
	b := data[2:]

	evalCtx := tree.NewTestingEvalContext(cluster.MakeTestingClusterSettings())
	defer evalCtx.Stop(context.Background())

	_, err := DecodeDatum(evalCtx, typ, code, b)
	if err != nil {
		return 0
	}
	return 1
}
