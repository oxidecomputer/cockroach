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

package delegate

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

func (d *delegator) delegateShowDatabases(stmt *tree.ShowDatabases) (tree.Statement, error) {
	query := `SELECT
	name AS database_name, owner, primary_region, regions, survival_goal`

	if stmt.WithComment {
		query += `, comment`
	}

	query += `
FROM
  "".crdb_internal.databases d
`
	if stmt.WithComment {
		query += fmt.Sprintf(`
LEFT JOIN
	(
		SELECT 
			object_id, type, comment
		FROM
			system.comments
		WHERE
			type = %d
	) c
ON
	c.object_id = d.id`, keys.DatabaseCommentType)
	}

	query += `
ORDER BY
	database_name`

	return parse(query)
}
