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

	"github.com/cockroachdb/cockroach/pkg/sql/lexbase"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/cat"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqltelemetry"
)

// delegateShowSchemas implements SHOW SCHEMAS which returns all the schemas in
// the given or current database.
// Privileges: None.
func (d *delegator) delegateShowSchemas(n *tree.ShowSchemas) (tree.Statement, error) {
	name, err := d.getSpecifiedOrCurrentDatabase(n.Database)
	if err != nil {
		return nil, err
	}
	getSchemasQuery := fmt.Sprintf(`
      SELECT nspname AS schema_name, rolname AS owner
      FROM %[1]s.information_schema.schemata i
      INNER JOIN %[1]s.pg_catalog.pg_namespace n ON (n.nspname = i.schema_name)
      LEFT JOIN %[1]s.pg_catalog.pg_roles r ON (n.nspowner = r.oid)
			WHERE catalog_name = %[2]s
			ORDER BY schema_name`,
		name.String(), // note: (tree.Name).String() != string(name)
		lexbase.EscapeSQLString(string(name)),
	)

	return parse(getSchemasQuery)
}

func (d *delegator) delegateShowCreateAllSchemas() (tree.Statement, error) {
	sqltelemetry.IncrementShowCounter(sqltelemetry.Create)

	const showCreateAllSchemasQuery = `
	SELECT crdb_internal.show_create_all_schemas(%[1]s) AS create_statement;
`
	databaseLiteral := d.evalCtx.SessionData().Database

	query := fmt.Sprintf(showCreateAllSchemasQuery,
		lexbase.EscapeSQLString(databaseLiteral),
	)

	return parse(query)
}

// getSpecifiedOrCurrentDatabase returns the name of the specified database, or
// of the current database if the specified name is empty.
//
// Returns an error if there is no current database, or if the specified
// database doesn't exist.
func (d *delegator) getSpecifiedOrCurrentDatabase(specifiedDB tree.Name) (tree.Name, error) {
	var name cat.SchemaName
	if specifiedDB != "" {
		// Note: the schema name may be interpreted as database name,
		// see name_resolution.go.
		name.SchemaName = specifiedDB
		name.ExplicitSchema = true
	}

	flags := cat.Flags{AvoidDescriptorCaches: true}
	_, resName, err := d.catalog.ResolveSchema(d.ctx, flags, &name)
	if err != nil {
		return "", err
	}
	return resName.CatalogName, nil
}
