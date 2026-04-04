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
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqltelemetry"
)

// delegateShowRoles implements SHOW ROLES which returns all the roles.
// Privileges: SELECT on system.users.
func (d *delegator) delegateShowRoles() (tree.Statement, error) {
	sqltelemetry.IncrementShowCounter(sqltelemetry.Roles)
	return parse(`
SELECT
	u.username,
	IFNULL(string_agg(o.option || COALESCE('=' || o.value, ''), ', ' ORDER BY o.option), '') AS options,
	ARRAY (SELECT role FROM system.role_members AS rm WHERE rm.member = u.username ORDER BY 1) AS member_of
FROM
	system.users AS u LEFT JOIN system.role_options AS o ON u.username = o.username
GROUP BY
	u.username
ORDER BY 1;
`)
}
