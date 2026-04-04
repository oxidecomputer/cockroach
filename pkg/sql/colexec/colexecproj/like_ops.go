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

package colexecproj

import (
	"github.com/cockroachdb/cockroach/pkg/sql/colexec/colexeccmp"
	"github.com/cockroachdb/cockroach/pkg/sql/colexec/colexecutils"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
	"github.com/cockroachdb/cockroach/pkg/sql/colmem"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/errors"
)

// GetLikeProjectionOperator returns a projection operator which projects the
// result of the specified LIKE pattern, or NOT LIKE if the negate argument is
// true. The implementation varies depending on the complexity of the pattern.
func GetLikeProjectionOperator(
	allocator *colmem.Allocator,
	ctx *tree.EvalContext,
	input colexecop.Operator,
	colIdx int,
	resultIdx int,
	pattern string,
	negate bool,
) (colexecop.Operator, error) {
	likeOpType, pattern, err := colexeccmp.GetLikeOperatorType(pattern, negate)
	if err != nil {
		return nil, err
	}
	pat := []byte(pattern)
	input = colexecutils.NewVectorTypeEnforcer(allocator, input, types.Bool, resultIdx)
	base := projConstOpBase{
		OneInputHelper: colexecop.MakeOneInputHelper(input),
		allocator:      allocator,
		colIdx:         colIdx,
		outputIdx:      resultIdx,
	}
	switch likeOpType {
	case colexeccmp.LikeConstant:
		return &projEQBytesBytesConstOp{
			projConstOpBase: base,
			constArg:        pat,
		}, nil
	case colexeccmp.LikeConstantNegate:
		return &projNEBytesBytesConstOp{
			projConstOpBase: base,
			constArg:        pat,
		}, nil
	case colexeccmp.LikeNeverMatch:
		// Use an empty not-prefix operator to get correct NULL behavior.
		return &projNotPrefixBytesBytesConstOp{
			projConstOpBase: base,
			constArg:        []byte{},
		}, nil
	case colexeccmp.LikeAlwaysMatch:
		// Use an empty prefix operator to get correct NULL behavior.
		return &projPrefixBytesBytesConstOp{
			projConstOpBase: base,
			constArg:        []byte{},
		}, nil
	case colexeccmp.LikeSuffix:
		return &projSuffixBytesBytesConstOp{
			projConstOpBase: base,
			constArg:        pat,
		}, nil
	case colexeccmp.LikeSuffixNegate:
		return &projNotSuffixBytesBytesConstOp{
			projConstOpBase: base,
			constArg:        pat,
		}, nil
	case colexeccmp.LikePrefix:
		return &projPrefixBytesBytesConstOp{
			projConstOpBase: base,
			constArg:        pat,
		}, nil
	case colexeccmp.LikePrefixNegate:
		return &projNotPrefixBytesBytesConstOp{
			projConstOpBase: base,
			constArg:        pat,
		}, nil
	case colexeccmp.LikeContains:
		return &projContainsBytesBytesConstOp{
			projConstOpBase: base,
			constArg:        pat,
		}, nil
	case colexeccmp.LikeContainsNegate:
		return &projNotContainsBytesBytesConstOp{
			projConstOpBase: base,
			constArg:        pat,
		}, nil
	case colexeccmp.LikeRegexp:
		re, err := tree.ConvertLikeToRegexp(ctx, pattern, false, '\\')
		if err != nil {
			return nil, err
		}
		return &projRegexpBytesBytesConstOp{
			projConstOpBase: base,
			constArg:        re,
		}, nil
	case colexeccmp.LikeRegexpNegate:
		re, err := tree.ConvertLikeToRegexp(ctx, pattern, false, '\\')
		if err != nil {
			return nil, err
		}
		return &projNotRegexpBytesBytesConstOp{
			projConstOpBase: base,
			constArg:        re,
		}, nil
	default:
		return nil, errors.AssertionFailedf("unsupported like op type %d", likeOpType)
	}
}
