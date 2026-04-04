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

package colexecjoin

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/colexecerror"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfra"
	"github.com/cockroachdb/errors"
)

// newJoinHelper returns an execinfra.OpNode with two Operator inputs.
func newJoinHelper(inputOne, inputTwo colexecop.Operator) *joinHelper {
	return &joinHelper{inputOne: inputOne, inputTwo: inputTwo}
}

type joinHelper struct {
	colexecop.InitHelper
	inputOne colexecop.Operator
	inputTwo colexecop.Operator
}

// init initializes both inputs and returns true if this is the first time init
// was called.
func (h *joinHelper) init(ctx context.Context) bool {
	if !h.Init(ctx) {
		return false
	}
	h.inputOne.Init(h.Ctx)
	h.inputTwo.Init(h.Ctx)
	return true
}

func (h *joinHelper) ChildCount(verbose bool) int {
	return 2
}

func (h *joinHelper) Child(nth int, verbose bool) execinfra.OpNode {
	switch nth {
	case 0:
		return h.inputOne
	case 1:
		return h.inputTwo
	}
	colexecerror.InternalError(errors.AssertionFailedf("invalid idx %d", nth))
	// This code is unreachable, but the compiler cannot infer that.
	return nil
}
