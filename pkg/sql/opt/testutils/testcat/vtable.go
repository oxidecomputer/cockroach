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

package testcat

import (
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/vtable"
	"github.com/cockroachdb/errors"
)

var informationSchemaMap = map[string]*tree.CreateTable{}
var pgCatalogMap = map[string]*tree.CreateTable{}

var informationSchemaTables = []string{
	vtable.InformationSchemaColumns,
	vtable.InformationSchemaAdministrableRoleAuthorizations,
	vtable.InformationSchemaApplicableRoles,
	vtable.InformationSchemaColumnPrivileges,
	vtable.InformationSchemaSchemata,
	vtable.InformationSchemaTables,
}

var pgCatalogTables = []string{
	vtable.PGCatalogAm,
	vtable.PGCatalogAttrDef,
	vtable.PGCatalogAttribute,
	vtable.PGCatalogCast,
	vtable.PGCatalogAuthID,
	vtable.PGCatalogAuthMembers,
	vtable.PGCatalogAvailableExtensions,
	vtable.PGCatalogClass,
	vtable.PGCatalogCollation,
	vtable.PGCatalogConstraint,
	vtable.PGCatalogConversion,
	vtable.PGCatalogDatabase,
	vtable.PGCatalogDefaultACL,
	vtable.PGCatalogDepend,
	vtable.PGCatalogDescription,
	vtable.PGCatalogSharedDescription,
	vtable.PGCatalogEnum,
	vtable.PGCatalogEventTrigger,
	vtable.PGCatalogExtension,
	vtable.PGCatalogForeignDataWrapper,
	vtable.PGCatalogForeignServer,
	vtable.PGCatalogForeignTable,
	vtable.PGCatalogIndex,
	vtable.PGCatalogIndexes,
	vtable.PGCatalogInherits,
	vtable.PGCatalogLanguage,
	vtable.PGCatalogLocks,
	vtable.PGCatalogMatViews,
	vtable.PGCatalogNamespace,
	vtable.PGCatalogOperator,
	vtable.PGCatalogPreparedXacts,
	vtable.PGCatalogPreparedStatements,
	vtable.PGCatalogProc,
	vtable.PGCatalogRange,
	vtable.PGCatalogRewrite,
	vtable.PGCatalogRoles,
	vtable.PGCatalogSecLabels,
	vtable.PGCatalogSequence,
	vtable.PGCatalogSettings,
	vtable.PGCatalogShdepend,
	vtable.PGCatalogTables,
	vtable.PGCatalogTablespace,
	vtable.PGCatalogTrigger,
	vtable.PGCatalogType,
	vtable.PGCatalogUser,
	vtable.PGCatalogUserMapping,
	vtable.PGCatalogStatActivity,
	vtable.PGCatalogSecurityLabel,
	vtable.PGCatalogSharedSecurityLabel,
	vtable.PGCatalogViews,
	vtable.PGCatalogAggregate,
}

func init() {
	// Build a map that maps the names of the various virtual tables
	// to their CREATE TABLE AST.
	buildMap := func(schemaName string, tableList []string, tableMap map[string]*tree.CreateTable) {
		for _, table := range tableList {
			parsed, err := parser.ParseOne(table)
			if err != nil {
				panic(errors.Wrap(err, "error initializing virtual table map"))
			}

			ct, ok := parsed.AST.(*tree.CreateTable)
			if !ok {
				panic(errors.New("virtual table schemas must be CREATE TABLE statements"))
			}

			ct.Table.SchemaName = tree.Name(schemaName)
			ct.Table.ExplicitSchema = true

			ct.Table.CatalogName = testDB
			ct.Table.ExplicitCatalog = true

			name := ct.Table
			tableMap[name.ObjectName.String()] = ct
		}
	}

	buildMap("information_schema", informationSchemaTables, informationSchemaMap)
	buildMap("pg_catalog", pgCatalogTables, pgCatalogMap)
}

// Resolve returns true and the AST node describing the virtual table referenced.
// TODO(justin): make this complete for all virtual tables.
func resolveVTable(name *tree.TableName) (*tree.CreateTable, bool) {
	switch name.SchemaName {
	case "information_schema":
		schema, ok := informationSchemaMap[name.ObjectName.String()]
		return schema, ok

	case "pg_catalog":
		schema, ok := pgCatalogMap[name.ObjectName.String()]
		return schema, ok
	}

	return nil, false
}
