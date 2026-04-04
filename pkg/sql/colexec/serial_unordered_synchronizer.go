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

package colexec

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/colexec/colexecargs"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfra"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfrapb"
	"github.com/cockroachdb/cockroach/pkg/util/tracing"
)

// SerialUnorderedSynchronizer is an Operator that combines multiple Operator
// streams into one. It reads its inputs one by one until each one is exhausted,
// at which point it moves to the next input. See ParallelUnorderedSynchronizer
// for a parallel implementation. The serial one is used when concurrency is
// undesirable - for example when the whole query is planned on the gateway and
// we want to run it in the RootTxn.
type SerialUnorderedSynchronizer struct {
	colexecop.InitHelper
	span *tracing.Span

	inputs []colexecargs.OpWithMetaInfo
	// curSerialInputIdx indicates the index of the current input being consumed.
	curSerialInputIdx int
}

var (
	_ colexecop.Operator = &SerialUnorderedSynchronizer{}
	_ execinfra.OpNode   = &SerialUnorderedSynchronizer{}
	_ colexecop.Closer   = &SerialUnorderedSynchronizer{}
)

// ChildCount implements the execinfra.OpNode interface.
func (s *SerialUnorderedSynchronizer) ChildCount(verbose bool) int {
	return len(s.inputs)
}

// Child implements the execinfra.OpNode interface.
func (s *SerialUnorderedSynchronizer) Child(nth int, verbose bool) execinfra.OpNode {
	return s.inputs[nth].Root
}

// NewSerialUnorderedSynchronizer creates a new SerialUnorderedSynchronizer.
func NewSerialUnorderedSynchronizer(
	inputs []colexecargs.OpWithMetaInfo,
) *SerialUnorderedSynchronizer {
	return &SerialUnorderedSynchronizer{
		inputs: inputs,
	}
}

// Init is part of the colexecop.Operator interface.
func (s *SerialUnorderedSynchronizer) Init(ctx context.Context) {
	if !s.InitHelper.Init(ctx) {
		return
	}
	s.Ctx, s.span = execinfra.ProcessorSpan(s.Ctx, "serial unordered sync")
	for _, input := range s.inputs {
		input.Root.Init(s.Ctx)
	}
}

// Next is part of the colexecop.Operator interface.
func (s *SerialUnorderedSynchronizer) Next() coldata.Batch {
	for {
		if s.curSerialInputIdx == len(s.inputs) {
			return coldata.ZeroBatch
		}
		b := s.inputs[s.curSerialInputIdx].Root.Next()
		if b.Length() == 0 {
			s.curSerialInputIdx++
		} else {
			return b
		}
	}
}

// DrainMeta is part of the colexecop.MetadataSource interface.
func (s *SerialUnorderedSynchronizer) DrainMeta() []execinfrapb.ProducerMetadata {
	var bufferedMeta []execinfrapb.ProducerMetadata
	if s.span != nil {
		for i := range s.inputs {
			for _, stats := range s.inputs[i].StatsCollectors {
				s.span.RecordStructured(stats.GetStats())
			}
		}
		if meta := execinfra.GetTraceDataAsMetadata(s.span); meta != nil {
			bufferedMeta = append(bufferedMeta, *meta)
		}
	}
	for _, input := range s.inputs {
		bufferedMeta = append(bufferedMeta, input.MetadataSources.DrainMeta()...)
	}
	return bufferedMeta
}

// Close is part of the colexecop.ClosableOperator interface.
func (s *SerialUnorderedSynchronizer) Close(context.Context) error {
	// Note that we're using the context of the synchronizer rather than the
	// argument of Close() because the synchronizer derives its own tracing
	// span.
	ctx := s.EnsureCtx()
	var lastErr error
	for _, input := range s.inputs {
		if err := input.ToClose.Close(ctx); err != nil {
			lastErr = err
		}
	}
	if s.span != nil {
		s.span.Finish()
		s.span = nil
	}
	return lastErr
}
