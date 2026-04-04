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

package delegate

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/sql/lexbase"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// ShowRoleGrants returns role membership details for the specified roles and grantees.
// Privileges: SELECT on system.role_members.
//   Notes: postgres does not have a SHOW GRANTS ON ROLES statement.
func (d *delegator) delegateShowRoleGrants(n *tree.ShowRoleGrants) (tree.Statement, error) {
	const selectQuery = `
SELECT role AS role_name,
       member,
       "isAdmin" AS is_admin
 FROM system.role_members`

	var query bytes.Buffer
	query.WriteString(selectQuery)

	if n.Roles != nil {
		var roles []string
		sqlUsernames, err := n.Roles.ToSQLUsernames(d.evalCtx.SessionData(), security.UsernameValidation)
		if err != nil {
			return nil, err
		}
		for _, r := range sqlUsernames {
			roles = append(roles, lexbase.EscapeSQLString(r.Normalized()))
		}
		fmt.Fprintf(&query, ` WHERE "role" IN (%s)`, strings.Join(roles, ","))
	}

	if n.Grantees != nil {
		if n.Roles == nil {
			// No roles specified: we need a WHERE clause.
			query.WriteString(" WHERE ")
		} else {
			// We have a WHERE clause for roles.
			query.WriteString(" AND ")
		}

		var grantees []string
		granteeSQLUsernames, err := n.Grantees.ToSQLUsernames(d.evalCtx.SessionData(), security.UsernameValidation)
		if err != nil {
			return nil, err
		}
		for _, g := range granteeSQLUsernames {
			grantees = append(grantees, lexbase.EscapeSQLString(g.Normalized()))
		}
		fmt.Fprintf(&query, ` member IN (%s)`, strings.Join(grantees, ","))

	}

	return parse(query.String())
}
