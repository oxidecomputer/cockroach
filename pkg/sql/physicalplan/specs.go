// Copyright 2018 The Cockroach Authors.
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

package physicalplan

import (
	"sync"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfrapb"
)

var flowSpecPool = sync.Pool{
	New: func() interface{} {
		return &execinfrapb.FlowSpec{}
	},
}

// NewFlowSpec returns a new FlowSpec, which may have non-zero capacity in its
// slice fields.
func NewFlowSpec(flowID execinfrapb.FlowID, gateway base.SQLInstanceID) *execinfrapb.FlowSpec {
	spec := flowSpecPool.Get().(*execinfrapb.FlowSpec)
	spec.FlowID = flowID
	spec.Gateway = gateway
	return spec
}

// ReleaseFlowSpec returns this FlowSpec back to the pool of FlowSpecs. It may
// not be used again after this call.
func ReleaseFlowSpec(spec *execinfrapb.FlowSpec) {
	for i := range spec.Processors {
		if tr := spec.Processors[i].Core.TableReader; tr != nil {
			releaseTableReaderSpec(tr)
		}
	}
	*spec = execinfrapb.FlowSpec{
		Processors: spec.Processors[:0],
	}
	flowSpecPool.Put(spec)
}

var trSpecPool = sync.Pool{
	New: func() interface{} {
		return &execinfrapb.TableReaderSpec{}
	},
}

// NewTableReaderSpec returns a new TableReaderSpec.
func NewTableReaderSpec() *execinfrapb.TableReaderSpec {
	return trSpecPool.Get().(*execinfrapb.TableReaderSpec)
}

// releaseTableReaderSpec puts this TableReaderSpec back into its sync pool. It
// may not be used again after Release returns.
func releaseTableReaderSpec(s *execinfrapb.TableReaderSpec) {
	*s = execinfrapb.TableReaderSpec{}
	trSpecPool.Put(s)
}
