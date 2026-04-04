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

package protectedts

import (
	"time"

	"github.com/cockroachdb/cockroach/pkg/settings"
)

// Records and their spans are stored in memory on every host so it's best
// not to let this data size be unbounded.

// MaxBytes controls the maximum number of bytes worth of spans and metadata
// which can be protected by all protected timestamp records.
var MaxBytes = settings.RegisterIntSetting(
	settings.TenantWritable,
	"kv.protectedts.max_bytes",
	"if non-zero the limit of the number of bytes of spans and metadata which can be protected",
	1<<20, // 1 MiB
	settings.NonNegativeInt,
)

// MaxSpans controls the maximum number of spans which can be protected
// by all protected timestamp records.
var MaxSpans = settings.RegisterIntSetting(
	settings.TenantWritable,
	"kv.protectedts.max_spans",
	"if non-zero the limit of the number of spans which can be protected",
	32768,
	settings.NonNegativeInt,
)

// PollInterval defines how frequently the protectedts state is polled by the
// Tracker.
var PollInterval = settings.RegisterDurationSetting(
	settings.TenantWritable,
	"kv.protectedts.poll_interval",
	// TODO(ajwerner): better description.
	"the interval at which the protectedts subsystem state is polled",
	2*time.Minute, settings.NonNegativeDuration)

func init() {
	MaxBytes.SetVisibility(settings.Reserved)
	MaxSpans.SetVisibility(settings.Reserved)
}
