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

package colexecsel

import (
	"github.com/cockroachdb/cockroach/pkg/sql/colexec/colexeccmp"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/errors"
)

// GetLikeOperator returns a selection operator which applies the specified LIKE
// pattern, or NOT LIKE if the negate argument is true. The implementation
// varies depending on the complexity of the pattern.
func GetLikeOperator(
	ctx *tree.EvalContext, input colexecop.Operator, colIdx int, pattern string, negate bool,
) (colexecop.Operator, error) {
	likeOpType, pattern, err := colexeccmp.GetLikeOperatorType(pattern, negate)
	if err != nil {
		return nil, err
	}
	pat := []byte(pattern)
	base := selConstOpBase{
		OneInputHelper: colexecop.MakeOneInputHelper(input),
		colIdx:         colIdx,
	}
	switch likeOpType {
	case colexeccmp.LikeConstant:
		return &selEQBytesBytesConstOp{
			selConstOpBase: base,
			constArg:       pat,
		}, nil
	case colexeccmp.LikeConstantNegate:
		return &selNEBytesBytesConstOp{
			selConstOpBase: base,
			constArg:       pat,
		}, nil
	case colexeccmp.LikeNeverMatch:
		// Use an empty not-prefix operator to get correct NULL behavior.
		return &selNotPrefixBytesBytesConstOp{
			selConstOpBase: base,
			constArg:       []byte{},
		}, nil
	case colexeccmp.LikeAlwaysMatch:
		// Use an empty prefix operator to get correct NULL behavior.
		return &selPrefixBytesBytesConstOp{
			selConstOpBase: base,
			constArg:       []byte{},
		}, nil
	case colexeccmp.LikeSuffix:
		return &selSuffixBytesBytesConstOp{
			selConstOpBase: base,
			constArg:       pat,
		}, nil
	case colexeccmp.LikeSuffixNegate:
		return &selNotSuffixBytesBytesConstOp{
			selConstOpBase: base,
			constArg:       pat,
		}, nil
	case colexeccmp.LikePrefix:
		return &selPrefixBytesBytesConstOp{
			selConstOpBase: base,
			constArg:       pat,
		}, nil
	case colexeccmp.LikePrefixNegate:
		return &selNotPrefixBytesBytesConstOp{
			selConstOpBase: base,
			constArg:       pat,
		}, nil
	case colexeccmp.LikeContains:
		return &selContainsBytesBytesConstOp{
			selConstOpBase: base,
			constArg:       pat,
		}, nil
	case colexeccmp.LikeContainsNegate:
		return &selNotContainsBytesBytesConstOp{
			selConstOpBase: base,
			constArg:       pat,
		}, nil
	case colexeccmp.LikeRegexp:
		re, err := tree.ConvertLikeToRegexp(ctx, pattern, false, '\\')
		if err != nil {
			return nil, err
		}
		return &selRegexpBytesBytesConstOp{
			selConstOpBase: base,
			constArg:       re,
		}, nil
	case colexeccmp.LikeRegexpNegate:
		re, err := tree.ConvertLikeToRegexp(ctx, pattern, false, '\\')
		if err != nil {
			return nil, err
		}
		return &selNotRegexpBytesBytesConstOp{
			selConstOpBase: base,
			constArg:       re,
		}, nil
	default:
		return nil, errors.AssertionFailedf("unsupported like op type %d", likeOpType)
	}
}
