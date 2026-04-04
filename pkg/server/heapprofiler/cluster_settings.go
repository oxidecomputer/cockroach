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

package heapprofiler

import "github.com/cockroachdb/cockroach/pkg/settings"

// ActiveQueryDumpsEnabled wraps "diagnostics.active_query_dumps.enabled"
//
// diagnostics.active_query_dumps.enabled enables the periodic writing of
// active queries on a node to disk, in *.csv format, if a node is determined to
// be under memory pressure.
//
// Note: this feature only works for nodes running on unix hosts with cgroups
// enabled.
var ActiveQueryDumpsEnabled = settings.RegisterBoolSetting(
	settings.SystemOnly,
	"diagnostics.active_query_dumps.enabled",
	"experimental: enable dumping of anonymized active queries to disk when node is under memory pressure",
	true,
).WithPublic()
