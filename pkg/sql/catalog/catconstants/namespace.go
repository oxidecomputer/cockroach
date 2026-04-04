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

package catconstants

const (
	// NamespaceTableFamilyID is the column family of the namespace table which is
	// actually written to.
	NamespaceTableFamilyID = 4

	// NamespaceTablePrimaryIndexID is the id of the primary index of the
	// namespace table.
	NamespaceTablePrimaryIndexID = 1

	// PreMigrationNamespaceTableName is the name that was used on the descriptor
	// of the current namespace table before the DeprecatedNamespaceTable was
	// migrated away.
	PreMigrationNamespaceTableName = "namespace2"
)
