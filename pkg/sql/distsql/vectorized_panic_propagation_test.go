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

package distsql

import (
	"context"
	"sync"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/colexec"
	"github.com/cockroachdb/cockroach/pkg/sql/colexec/colexecargs"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
	"github.com/cockroachdb/cockroach/pkg/sql/colflow"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfra"
	"github.com/cockroachdb/cockroach/pkg/sql/flowinfra"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/stretchr/testify/require"
)

// TestNonVectorizedPanicDoesntHangServer verifies that propagating a non
// vectorized panic doesn't result in a hang as described in:
// https://github.com/cockroachdb/cockroach/issues/39779
func TestNonVectorizedPanicDoesntHangServer(t *testing.T) {
	defer leaktest.AfterTest(t)()
	ctx := context.Background()
	st := cluster.MakeTestingClusterSettings()
	evalCtx := tree.MakeTestingEvalContext(st)
	defer evalCtx.Stop(ctx)

	flowCtx := execinfra.FlowCtx{
		EvalCtx: &evalCtx,
		Cfg:     &execinfra.ServerConfig{Settings: cluster.MakeTestingClusterSettings()},
	}
	base := flowinfra.NewFlowBase(
		flowCtx,
		nil, /* sp */
		nil, /* flowReg */
		nil, /* rowSyncFlowConsumer */
		nil, /* batchSyncFlowConsumer */
		nil, /* localProcessors */
		nil, /* onFlowCleanup */
		"",  /* statementSQL */
	)
	flow := colflow.NewVectorizedFlow(base)

	mat := colexec.NewMaterializer(
		&flowCtx,
		0, /* processorID */
		colexecargs.OpWithMetaInfo{Root: &colexecop.CallbackOperator{
			NextCb: func() coldata.Batch {
				panic("")
			},
		}},
		nil, /* typs */
	)

	ctx, _, err := base.Setup(ctx, nil, flowinfra.FuseAggressively)
	require.NoError(t, err)

	base.SetProcessors([]execinfra.Processor{mat})
	// This test specifically verifies that a flow doesn't get stuck in Wait for
	// asynchronous components that haven't been signaled to exit. To simulate
	// this we just create a mock startable.
	flow.AddStartable(
		flowinfra.StartableFn(func(ctx context.Context, wg *sync.WaitGroup, _ context.CancelFunc) {
			wg.Add(1)
			go func() {
				// Ensure context is canceled.
				<-ctx.Done()
				wg.Done()
			}()
		}),
	)

	require.Panics(t, func() { flow.Run(ctx, nil) })
}
