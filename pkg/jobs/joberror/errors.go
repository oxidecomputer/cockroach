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

package joberror

import (
	"strings"

	circuitbreaker "github.com/cockroachdb/circuitbreaker"
	"github.com/cockroachdb/cockroach/pkg/kv/kvclient/kvcoord"
	"github.com/cockroachdb/cockroach/pkg/sql/flowinfra"
	"github.com/cockroachdb/cockroach/pkg/util/circuit"
	"github.com/cockroachdb/cockroach/pkg/util/grpcutil"
	"github.com/cockroachdb/cockroach/pkg/util/sysutil"
	"github.com/cockroachdb/errors"
)

// IsDistSQLRetryableError returns true if the supplied error, or any of its parent
// causes is an rpc error.
// This is an unfortunate implementation that should be looking for a more
// specific error.
func IsDistSQLRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// TODO(knz): this is a bad implementation. Make it go away
	// by avoiding string comparisons.

	errStr := err.Error()
	// When a crdb node dies, any DistSQL flows with processors scheduled on
	// it get an error with "rpc error" in the message from the call to
	// `(*DistSQLPlanner).Run`.
	return strings.Contains(errStr, `rpc error`)
}

// isBreakerOpenError returns true if err is a circuit.ErrBreakerOpen.
//
// NB: Two packages have ErrBreakerOpen error types.  The cicruitbreaker package
// is used by the nodedialer. The circuit package is used by kvserver.
func isBreakerOpenError(err error) bool {
	return errors.Is(err, circuit.ErrBreakerOpen) || errors.Is(err, circuitbreaker.ErrBreakerOpen)
}

// IsPermanentBulkJobError returns true if the error results in a permanent
// failure of a bulk job (IMPORT, BACKUP, RESTORE). This function is an
// allowlist instead of a blocklist: only known safe errors are confirmed to not
// be permanent errors. Anything unknown is assumed to be permanent.
func IsPermanentBulkJobError(err error) bool {
	if err == nil {
		return false
	}

	return !IsDistSQLRetryableError(err) &&
		!grpcutil.IsClosedConnection(err) &&
		!flowinfra.IsNoInboundStreamConnectionError(err) &&
		!kvcoord.IsSendError(err) &&
		!isBreakerOpenError(err) &&
		!sysutil.IsErrConnectionReset(err) &&
		!sysutil.IsErrConnectionRefused(err)
}
