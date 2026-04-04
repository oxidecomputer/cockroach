// Copyright 2014 The Cockroach Authors.
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

	"github.com/cockroachdb/cockroach/pkg/util/log/logpb"
	"github.com/cockroachdb/cockroach/pkg/util/tracing"
	"github.com/cockroachdb/errors"
)

func init() {
	errors.SetWarningFn(Warningf)
}

// Severity aliases a type.
type Severity = logpb.Severity

// V returns true if the logging verbosity is set to the specified level or
// higher.
//
// See also ExpensiveLogEnabled().
//
// TODO(andrei): Audit uses of V() and see which ones should actually use the
// newer ExpensiveLogEnabled().
func V(level Level) bool {
	return VDepth(level, 1)
}

// ExpensiveLogEnabled is used to test whether effort should be used to produce
// log messages whose construction has a measurable cost. It returns true if
// either the current context is recording the trace, or if the caller's
// verbosity is above level.
//
// NOTE: This doesn't take into consideration whether tracing is generally
// enabled or whether a trace.EventLog or a trace.Trace (i.e. sp.netTr) is
// attached to ctx. In particular, if OpenTelemetry collection is enabled,
// that, by itself, does NOT cause the expensive messages to
// be enabled. "SET tracing" and friends, on the other hand, does cause
// these messages to be enabled, as it shows that a user has expressed
// particular interest in a trace.
//
// Usage:
//
// if ExpensiveLogEnabled(ctx, 2) {
//   msg := constructExpensiveMessage()
//   log.VEventf(ctx, 2, msg)
// }
//
func ExpensiveLogEnabled(ctx context.Context, level Level) bool {
	if sp := tracing.SpanFromContext(ctx); sp != nil {
		if sp.IsVerbose() || sp.Tracer().HasExternalSink() {
			return true
		}
	}
	if VDepth(level, 1 /* depth */) {
		return true
	}
	return false
}
