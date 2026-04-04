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

package txnidcache

import "github.com/cockroachdb/cockroach/pkg/settings"

// MaxSize limits the maximum byte size can be used by the TxnIDCache.
var MaxSize = settings.RegisterByteSizeSetting(
	settings.TenantWritable,
	`sql.contention.txn_id_cache.max_size`,
	"the maximum byte size TxnID cache will use (set to 0 to disable)",
	64*1024*1024, // 64MiB
).WithPublic()
