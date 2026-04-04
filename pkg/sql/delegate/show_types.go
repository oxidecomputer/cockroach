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

func (d *delegator) delegateShowTypes() (tree.Statement, error) {
	// TODO (SQL Features, SQL Exec): Once more user defined types are added
	//  they should be added here.
	return parse(`
SELECT
  schema, name, owner
FROM
  [SHOW ENUMS]
ORDER BY
  (schema, name)`)
}

func (d *delegator) delegateShowCreateAllTypes() (tree.Statement, error) {
	sqltelemetry.IncrementShowCounter(sqltelemetry.Create)

	const showCreateAllTypesQuery = `
	SELECT crdb_internal.show_create_all_types(%[1]s) AS create_statement;
`
	databaseLiteral := d.evalCtx.SessionData().Database

	query := fmt.Sprintf(showCreateAllTypesQuery,
		lexbase.EscapeSQLString(databaseLiteral),
	)

	return parse(query)
}
