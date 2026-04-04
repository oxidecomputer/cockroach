// Copyright 2020 The Cockroach Authors.
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

package log

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/util/log/eventpb"
	"github.com/cockroachdb/cockroach/pkg/util/log/severity"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
)

// StructuredEvent emits a structured event to the debug log.
func StructuredEvent(ctx context.Context, event eventpb.EventPayload) {
	// Populate the missing common fields.
	common := event.CommonDetails()
	if common.Timestamp == 0 {
		common.Timestamp = timeutil.Now().UnixNano()
	}
	if len(common.EventType) == 0 {
		common.EventType = eventpb.GetEventTypeName(event)
	}

	entry := makeStructuredEntry(ctx,
		severity.INFO,
		event.LoggingChannel(),
		// Note: we use depth 0 intentionally here, so that structured
		// events can be reliably detected (their source filename will
		// always be log/event_log.go).
		0, /* depth */
		event)

	if sp, el, ok := getSpanOrEventLog(ctx); ok {
		// Prevent `entry` from moving to the heap when this branch is not taken.
		heapEntry := entry
		eventInternal(sp, el, entry.sev >= severity.ERROR, &heapEntry)
	}

	logger := logging.getLogger(entry.ch)
	logger.outputLogEntry(entry)
}
