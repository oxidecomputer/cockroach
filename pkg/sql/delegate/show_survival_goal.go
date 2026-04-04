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

package delegate

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/sql/lexbase"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqltelemetry"
)

// delegateShowRanges implements the SHOW REGIONS statement.
func (d *delegator) delegateShowSurvivalGoal(n *tree.ShowSurvivalGoal) (tree.Statement, error) {
	sqltelemetry.IncrementShowCounter(sqltelemetry.SurvivalGoal)
	dbName := string(n.DatabaseName)
	if dbName == "" {
		dbName = d.evalCtx.SessionData().Database
	}
	query := fmt.Sprintf(
		`SELECT
	name AS "database",
	survival_goal
FROM crdb_internal.databases
WHERE name = %s`,
		lexbase.EscapeSQLString(dbName),
	)
	return parse(query)
}
