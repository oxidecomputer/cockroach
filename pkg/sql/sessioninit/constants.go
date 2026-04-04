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

package sessioninit

import (
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// UsersTableName represents system.users.
var UsersTableName = tree.NewTableNameWithSchema("system", tree.PublicSchemaName, "users")

// RoleOptionsTableName represents system.role_options.
var RoleOptionsTableName = tree.NewTableNameWithSchema("system", tree.PublicSchemaName, "role_options")

// DatabaseRoleSettingsTableName represents system.database_role_settings.
var DatabaseRoleSettingsTableName = tree.NewTableNameWithSchema("system", tree.PublicSchemaName, "database_role_settings")

// defaultDatabaseID is used in the settingsCache for entries that should
// apply to all database.
const defaultDatabaseID = 0

// defaultUsername is used in the settingsCache for entries that should
// apply to all roles.
var defaultUsername = security.MakeSQLUsernameFromPreNormalizedString("")
