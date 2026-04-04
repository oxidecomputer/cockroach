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

package sql

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfrapb"
	"github.com/cockroachdb/cockroach/pkg/sql/rowexec"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/errors"
)

// PlanAndRunCTAS plans and runs the CREATE TABLE AS command.
func PlanAndRunCTAS(
	ctx context.Context,
	dsp *DistSQLPlanner,
	planner *planner,
	txn *kv.Txn,
	isLocal bool,
	in planMaybePhysical,
	out execinfrapb.ProcessorCoreUnion,
	recv *DistSQLReceiver,
) {
	distribute := DistributionType(DistributionTypeNone)
	if !isLocal {
		distribute = DistributionTypeSystemTenantOnly
	}
	planCtx := dsp.NewPlanningCtx(ctx, planner.ExtendedEvalContext(), planner,
		txn, distribute)
	planCtx.stmtType = tree.Rows

	physPlan, cleanup, err := dsp.createPhysPlan(ctx, planCtx, in)
	defer cleanup()
	if err != nil {
		recv.SetError(errors.Wrapf(err, "constructing distSQL plan"))
		return
	}
	physPlan.AddNoGroupingStage(
		out, execinfrapb.PostProcessSpec{}, rowexec.CTASPlanResultTypes, execinfrapb.Ordering{},
	)

	// The bulk row writers will emit a binary encoded BulkOpSummary.
	physPlan.PlanToStreamColMap = []int{0}

	// Make copy of evalCtx as Run might modify it.
	evalCtxCopy := planner.ExtendedEvalContextCopy()
	dsp.FinalizePlan(planCtx, physPlan)
	dsp.Run(ctx, planCtx, txn, physPlan, recv, evalCtxCopy, nil /* finishedSetupFn */)()
}
