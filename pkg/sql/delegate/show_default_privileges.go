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

package delegate

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/sql/lexbase"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// delegateShowDefaultPrivileges implements SHOW DEFAULT PRIVILEGES
// which returns default privileges for a specified role.
func (d *delegator) delegateShowDefaultPrivileges(
	n *tree.ShowDefaultPrivileges,
) (tree.Statement, error) {
	currentDatabase, err := d.getSpecifiedOrCurrentDatabase("")
	if err != nil {
		return nil, err
	}

	schemaClause := " AND schema_name IS NULL"
	if n.Schema != "" {
		schemaClause = fmt.Sprintf(" AND schema_name = %s", lexbase.EscapeSQLString(string(n.Schema)))
	}

	query := fmt.Sprintf(
		"SELECT role, for_all_roles, object_type, grantee, privilege_type FROM crdb_internal.default_privileges WHERE database_name = %s%s",
		lexbase.EscapeSQLString(string(currentDatabase)),
		schemaClause,
	)
	if d.evalCtx.Settings.Version.IsActive(d.ctx, clusterversion.ValidateGrantOption) {
		query = fmt.Sprintf(
			"SELECT role, for_all_roles, object_type, grantee, privilege_type, is_grantable "+
				"FROM crdb_internal.default_privileges WHERE database_name = %s%s",
			lexbase.EscapeSQLString(string(currentDatabase)),
			schemaClause,
		)
	}

	if n.ForAllRoles {
		query += " AND for_all_roles=true"
	} else if len(n.Roles) > 0 {
		targetRoles, err := n.Roles.ToSQLUsernames(d.evalCtx.SessionData(), security.UsernameValidation)
		if err != nil {
			return nil, err
		}

		query = fmt.Sprintf("%s AND for_all_roles=false AND role IN (", query)
		for i, role := range targetRoles {
			if i != 0 {
				query += fmt.Sprintf(", '%s'", role.Normalized())
			} else {
				query += fmt.Sprintf("'%s'", role.Normalized())
			}
		}

		query += ")"
	} else {
		query = fmt.Sprintf("%s AND for_all_roles=false AND role = '%s'",
			query, d.evalCtx.SessionData().User())
	}
	query += " ORDER BY 1,2,3,4,5"
	return parse(query)
}
