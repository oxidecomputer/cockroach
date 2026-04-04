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

package colflow

import (
	"context"
	"math/rand"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecerror"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
	"github.com/cockroachdb/errors"
)

// panicInjector is a helper Operator that will randomly inject panics into
// Init and Next methods of the wrapped operator.
type panicInjector struct {
	colexecop.OneInputNode
	colexecop.InitHelper
	rng *rand.Rand
}

var _ colexecop.Operator = &panicInjector{}

const (
	// These constants were chosen arbitrarily with the guiding thought that
	// Init() methods are called less frequently, so the probability of
	// injecting in Init() should be higher. At the same time, we don't want
	// for the vectorized flows to always run into these injected panics, so
	// both numbers are relatively low.
	initPanicInjectionProbability = 0.001
	nextPanicInjectionProbability = 0.00001
)

// newPanicInjector creates a new panicInjector.
func newPanicInjector(input colexecop.Operator) colexecop.Operator {
	rng, _ := randutil.NewTestRand()
	return &panicInjector{
		OneInputNode: colexecop.OneInputNode{Input: input},
		rng:          rng,
	}
}

func (i *panicInjector) Init(ctx context.Context) {
	if !i.InitHelper.Init(ctx) {
		return
	}
	if i.rng.Float64() < initPanicInjectionProbability {
		log.Info(i.Ctx, "injecting panic in Init")
		colexecerror.ExpectedError(errors.New("injected panic in Init"))
	}
	i.Input.Init(i.Ctx)
}

func (i *panicInjector) Next() coldata.Batch {
	if i.rng.Float64() < nextPanicInjectionProbability {
		log.Info(i.Ctx, "injecting panic in Next")
		colexecerror.ExpectedError(errors.New("injected panic in Next"))
	}
	return i.Input.Next()
}
