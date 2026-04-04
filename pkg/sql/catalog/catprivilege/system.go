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

package catprivilege

import (
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catconstants"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/privilege"
)

var (
	readSystemTables = []catconstants.SystemTableName{
		catconstants.NamespaceTableName,
		catconstants.DescriptorTableName,
		catconstants.DescIDSequenceTableName,
		catconstants.TenantsTableName,
		catconstants.ProtectedTimestampsMetaTableName,
		catconstants.ProtectedTimestampsRecordsTableName,
		catconstants.StatementStatisticsTableName,
		catconstants.TransactionStatisticsTableName,
		// TODO(postamar): remove in 21.2
		catconstants.PreMigrationNamespaceTableName,
	}

	readWriteSystemTables = []catconstants.SystemTableName{
		catconstants.UsersTableName,
		catconstants.ZonesTableName,
		catconstants.SettingsTableName,
		catconstants.LeaseTableName,
		catconstants.EventLogTableName,
		catconstants.RangeEventTableName,
		catconstants.UITableName,
		catconstants.JobsTableName,
		catconstants.WebSessionsTableName,
		catconstants.TableStatisticsTableName,
		catconstants.LocationsTableName,
		catconstants.RoleMembersTableName,
		catconstants.CommentsTableName,
		catconstants.ReportsMetaTableName,
		catconstants.ReplicationConstraintStatsTableName,
		catconstants.ReplicationCriticalLocalitiesTableName,
		catconstants.ReplicationStatsTableName,
		catconstants.RoleOptionsTableName,
		catconstants.StatementBundleChunksTableName,
		catconstants.StatementDiagnosticsRequestsTableName,
		catconstants.StatementDiagnosticsTableName,
		catconstants.ScheduledJobsTableName,
		catconstants.SqllivenessTableName,
		catconstants.MigrationsTableName,
		catconstants.JoinTokensTableName,
		catconstants.DatabaseRoleSettingsTableName,
		catconstants.TenantUsageTableName,
		catconstants.SQLInstancesTableName,
		catconstants.SpanConfigurationsTableName,
		catconstants.TenantSettingsTableName,
		catconstants.SpanCountTableName,
	}

	// RestoreCopySystemTablePrefix is the prefix of the table name that we give
	// to the copy of the system table we are moving to a higher ID during
	// restore.
	//
	// TODO(adityamaru,dt): Remove once we fix the handling of dynamic system
	// table IDs during restore.
	RestoreCopySystemTablePrefix = "crdb_internal_copy"

	systemSuperuserPrivileges = func() map[descpb.NameInfo]privilege.List {
		m := make(map[descpb.NameInfo]privilege.List)
		tableKey := descpb.NameInfo{
			ParentID:       keys.SystemDatabaseID,
			ParentSchemaID: keys.SystemPublicSchemaID,
		}
		for _, rw := range readWriteSystemTables {
			tableKey.Name = string(rw)
			m[tableKey] = privilege.ReadWriteData
		}
		for _, r := range readSystemTables {
			tableKey.Name = string(r)
			m[tableKey] = privilege.ReadData
		}
		m[descpb.NameInfo{Name: catconstants.SystemDatabaseName}] = privilege.List{privilege.CONNECT}
		return m
	}()
)

// SystemSuperuserPrivileges returns the privilege list for super-users found
// for the given system descriptor name key. Returns nil if none was found.
func SystemSuperuserPrivileges(nameKey catalog.NameKey) privilege.List {
	key := descpb.NameInfo{
		ParentID:       nameKey.GetParentID(),
		ParentSchemaID: nameKey.GetParentSchemaID(),
		Name:           nameKey.GetName(),
	}
	return systemSuperuserPrivileges[key]
}
