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

package scplan

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/jobs/jobspb"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scerrors"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scplan/internal/opgen"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scplan/internal/rules"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scplan/internal/scgraph"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scplan/internal/scstage"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/redact"
)

// Params holds the arguments for planning.
type Params struct {
	// InRollback is used to indicate whether we've already been reverted.
	// Note that when in rollback, there is no turning back and all work is
	// non-revertible. Theory dictates that this is fine because of how we
	// had carefully crafted stages to only allow entering rollback while it
	// remains safe to do so.
	InRollback bool

	// ExecutionPhase indicates the phase that the plan should be constructed for.
	ExecutionPhase scop.Phase

	// SchemaChangerJobIDSupplier is used to return the JobID for a
	// job if one should exist.
	SchemaChangerJobIDSupplier func() jobspb.JobID
}

// Exported internal types
type (
	// Graph is an exported alias of scgraph.Graph.
	Graph = scgraph.Graph

	// Stage is an exported alias of scstage.Stage.
	Stage = scstage.Stage
)

// A Plan is a schema change plan, primarily containing ops to be executed that
// are partitioned into stages.
type Plan struct {
	scpb.CurrentState
	Params Params
	Graph  *scgraph.Graph
	JobID  jobspb.JobID
	Stages []Stage
}

// StagesForCurrentPhase returns the stages in the execution phase specified in
// the plan params.
func (p Plan) StagesForCurrentPhase() []scstage.Stage {
	for i, s := range p.Stages {
		if s.Phase > p.Params.ExecutionPhase {
			return p.Stages[:i]
		}
	}
	return p.Stages
}

// MakePlan generates a Plan for a particular phase of a schema change, given
// the initial state for a set of targets.
// Returns an error when planning fails. It is up to the caller to wrap this
// error as an assertion failure and with useful debug information details.
func MakePlan(ctx context.Context, initial scpb.CurrentState, params Params) (p Plan, err error) {
	defer scerrors.StartEventf(
		ctx,
		"building declarative schema changer plan in %s (rollback=%v) for %s",
		redact.Safe(params.ExecutionPhase),
		redact.Safe(params.InRollback),
		redact.Safe(initial.StatementTags()),
	).HandlePanicAndLogError(ctx, &err)
	p = Plan{
		CurrentState: initial,
		Params:       params,
	}
	{
		start := timeutil.Now()
		p.Graph = buildGraph(p.CurrentState)
		if log.V(2) {
			log.Infof(context.TODO(), "graph generation took %v", timeutil.Since(start))
		}
	}
	{
		start := timeutil.Now()
		p.Stages = scstage.BuildStages(
			initial, params.ExecutionPhase, p.Graph, params.SchemaChangerJobIDSupplier,
		)
		if log.V(2) {
			log.Infof(context.TODO(), "stage generation took %v", timeutil.Since(start))
		}
	}
	if n := len(p.Stages); n > 0 && p.Stages[n-1].Phase > scop.PreCommitPhase {
		// Only get the job ID if it's actually been assigned already.
		p.JobID = params.SchemaChangerJobIDSupplier()
	}
	if err := scstage.ValidateStages(p.TargetState, p.Stages, p.Graph); err != nil {
		panic(errors.Wrapf(err, "invalid execution plan"))
	}
	return p, nil
}

func buildGraph(cs scpb.CurrentState) *scgraph.Graph {
	g, err := opgen.BuildGraph(cs)
	if err != nil {
		panic(errors.Wrapf(err, "build graph op edges"))
	}
	err = rules.ApplyDepRules(g)
	if err != nil {
		panic(errors.Wrapf(err, "build graph dep edges"))
	}
	err = g.Validate()
	if err != nil {
		panic(errors.Wrapf(err, "validate graph"))
	}
	g, err = rules.ApplyOpRules(g)
	if err != nil {
		panic(errors.Wrapf(err, "mark op edges as no-op"))
	}
	return g
}
