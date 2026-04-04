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

package optbuilder

import (
	"github.com/cockroachdb/cockroach/pkg/server/telemetry"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/colinfo"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/memo"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqltelemetry"
	"github.com/cockroachdb/errors"
)

func (b *Builder) buildExplain(explain *tree.Explain, inScope *scope) (outScope *scope) {
	if _, ok := explain.Statement.(*tree.Execute); ok {
		panic(pgerror.New(
			pgcode.FeatureNotSupported, "EXPLAIN EXECUTE is not supported; use EXPLAIN ANALYZE",
		))
	}

	stmtScope := b.buildStmtAtRoot(explain.Statement, nil /* desiredTypes */)

	outScope = inScope.push()

	switch explain.Mode {
	case tree.ExplainPlan:
		telemetry.Inc(sqltelemetry.ExplainPlanUseCounter)

	case tree.ExplainDistSQL:
		telemetry.Inc(sqltelemetry.ExplainDistSQLUseCounter)

	case tree.ExplainOpt:
		if explain.Flags[tree.ExplainFlagVerbose] {
			telemetry.Inc(sqltelemetry.ExplainOptVerboseUseCounter)
		} else {
			telemetry.Inc(sqltelemetry.ExplainOptUseCounter)
		}

	case tree.ExplainVec:
		telemetry.Inc(sqltelemetry.ExplainVecUseCounter)

	case tree.ExplainDDL:
		if explain.Flags[tree.ExplainFlagViz] {
			telemetry.Inc(sqltelemetry.ExplainDDLViz)
		} else if explain.Flags[tree.ExplainFlagVerbose] {
			telemetry.Inc(sqltelemetry.ExplainDDLVerbose)
		} else {
			telemetry.Inc(sqltelemetry.ExplainDDL)
		}

	case tree.ExplainGist:
		telemetry.Inc(sqltelemetry.ExplainGist)

	default:
		panic(errors.Errorf("EXPLAIN mode %s not supported", explain.Mode))
	}
	b.synthesizeResultColumns(outScope, colinfo.ExplainPlanColumns)

	input := stmtScope.expr
	private := memo.ExplainPrivate{
		Options:  explain.ExplainOptions,
		ColList:  colsToColList(outScope.cols),
		Props:    stmtScope.makePhysicalProps(),
		StmtType: explain.Statement.StatementReturnType(),
	}
	outScope.expr = b.factory.ConstructExplain(input, &private)
	return outScope
}
