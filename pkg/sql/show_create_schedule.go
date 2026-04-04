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

package sql

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cockroachdb/cockroach/pkg/jobs"
	"github.com/cockroachdb/cockroach/pkg/scheduledjobs"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/colinfo"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/sql/sqltelemetry"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

var showCreateTableColumns = colinfo.ResultColumns{
	{Name: "schedule_id", Typ: types.Int},
	{Name: "create_statement", Typ: types.String},
}

const (
	scheduleID = iota
	createStmt
)

func loadSchedules(params runParams, n *tree.ShowCreateSchedules) ([]*jobs.ScheduledJob, error) {
	env := JobSchedulerEnv(params.ExecCfg())
	var schedules []*jobs.ScheduledJob
	var rows []tree.Datums
	var cols colinfo.ResultColumns

	if n.ScheduleID != nil {
		sjID, err := strconv.Atoi(n.ScheduleID.String())
		if err != nil {
			return nil, err
		}

		datums, columns, err := params.ExecCfg().InternalExecutor.QueryRowExWithCols(
			params.ctx,
			"load-schedules",
			params.EvalContext().Txn, sessiondata.InternalExecutorOverride{User: security.RootUserName()},
			fmt.Sprintf("SELECT * FROM %s WHERE schedule_id = $1", env.ScheduledJobsTableName()),
			tree.NewDInt(tree.DInt(sjID)))
		if err != nil {
			return nil, err
		}
		rows = append(rows, datums)
		cols = columns
	} else {
		datums, columns, err := params.ExecCfg().InternalExecutor.QueryBufferedExWithCols(
			params.ctx,
			"load-schedules",
			params.EvalContext().Txn, sessiondata.InternalExecutorOverride{User: security.RootUserName()},
			fmt.Sprintf("SELECT * FROM %s", env.ScheduledJobsTableName()))
		if err != nil {
			return nil, err
		}
		rows = append(rows, datums...)
		cols = columns
	}

	for _, row := range rows {
		schedule := jobs.NewScheduledJob(env)
		if err := schedule.InitFromDatums(row, cols); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

func (p *planner) ShowCreateSchedule(
	ctx context.Context, n *tree.ShowCreateSchedules,
) (planNode, error) {
	// Only admin users can execute SHOW CREATE SCHEDULE
	if userIsAdmin, err := p.UserHasAdminRole(ctx, p.User()); err != nil {
		return nil, err
	} else if !userIsAdmin {
		return nil, pgerror.Newf(pgcode.InsufficientPrivilege,
			"user %s does not have admin role", p.User())
	}

	sqltelemetry.IncrementShowCounter(sqltelemetry.CreateSchedule)

	return &delayedNode{
		name:    fmt.Sprintf("SHOW CREATE SCHEDULE %d", n.ScheduleID),
		columns: showCreateTableColumns,
		constructor: func(ctx context.Context, p *planner) (planNode, error) {
			scheduledJobs, err := loadSchedules(
				runParams{ctx: ctx, p: p, extendedEvalCtx: &p.extendedEvalCtx}, n)
			if err != nil {
				return nil, err
			}

			var rows []tree.Datums
			for _, sj := range scheduledJobs {
				ex, err := jobs.GetScheduledJobExecutor(sj.ExecutorType())
				if err != nil {
					return nil, err
				}

				createStmtStr, err := ex.GetCreateScheduleStatement(
					ctx,
					scheduledjobs.ProdJobSchedulerEnv,
					p.Txn(),
					p.Descriptors(),
					sj,
					p.ExecCfg().InternalExecutor,
				)
				if err != nil {
					return nil, err
				}

				row := tree.Datums{
					scheduleID: tree.NewDInt(tree.DInt(sj.ScheduleID())),
					createStmt: tree.NewDString(createStmtStr),
				}
				rows = append(rows, row)
			}

			v := p.newContainerValuesNode(showCreateTableColumns, len(rows))
			for _, row := range rows {
				if _, err := v.rows.AddRow(ctx, row); err != nil {
					v.Close(ctx)
					return nil, err
				}
			}
			return v, nil
		},
	}, nil
}
