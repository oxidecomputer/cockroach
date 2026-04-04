// Copyright 2022 The Cockroach Authors.
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

package sslocal

import "github.com/cockroachdb/cockroach/pkg/settings"

// AssociateStmtWithTxnFingerprint determines whether to segment
// per-statment statistics by transaction fingerprint. While enabled by
// default, it may be useful to disable for workloads that run the same
// statements across many (ad-hoc) transaction fingerprints, producing
// higher-cardinality data in the system.statement_statistics table than
// the cleanup job is able to keep up with. See #78338.
var AssociateStmtWithTxnFingerprint = settings.RegisterBoolSetting(
	settings.TenantWritable,
	"sql.stats.associate_stmt_with_txn_fingerprint.enabled",
	"whether to segment per-statement query statistics by transaction fingerprint",
	true,
)
