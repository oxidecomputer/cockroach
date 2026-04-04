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

package ptpb

import (
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
)

// MakeClusterTarget returns a target, which when used in a Record, will
// protect the entire keyspace of the cluster.
func MakeClusterTarget() *Target {
	return &Target{Union: &Target_Cluster{Cluster: &Target_ClusterTarget{}}}
}

// MakeTenantsTarget returns a target, which when used in a Record, will
// protect the keyspace of all tenants in ids.
func MakeTenantsTarget(ids []roachpb.TenantID) *Target {
	return &Target{Union: &Target_Tenants{Tenants: &Target_TenantsTarget{IDs: ids}}}
}

// MakeSchemaObjectsTarget returns a target, which when used in a Record,
// will protect the keyspace of all schema objects (database/table).
func MakeSchemaObjectsTarget(ids descpb.IDs) *Target {
	return &Target{Union: &Target_SchemaObjects{SchemaObjects: &Target_SchemaObjectsTarget{IDs: ids}}}
}
