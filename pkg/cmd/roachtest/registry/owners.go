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

package registry

// Owner is a valid entry for the Owners field of a roachtest. They should be
// teams, not individuals.
type Owner string

// The allowable values of Owner.
const (
	OwnerSQLFoundations   Owner = `sql-foundations`
	OwnerDisasterRecovery Owner = `disaster-recovery`
	OwnerCDC              Owner = `cdc`
	OwnerKV               Owner = `kv`
	OwnerReplication      Owner = `replication`
	OwnerMultiRegion      Owner = `multiregion`
	OwnerObsInf           Owner = `obs-inf-prs`
	OwnerServer           Owner = `server`
	OwnerSQLQueries       Owner = `sql-queries`
	OwnerStorage          Owner = `storage`
	OwnerTestEng          Owner = `test-eng`
	OwnerDevInf           Owner = `dev-inf`
	OwnerMultiTenant      Owner = `multi-tenant`
)
