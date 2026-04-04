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

import "github.com/cockroachdb/cockroach/pkg/cli/exit"

//go:generate mockgen -package=log -destination=mocks_generated_test.go --mock_names=TestingLogSink=MockLogSink . TestingLogSink

// TestingLogSink is exported for mock generation.
// This is painful, but it seems to be the only way, for the moment, to
// generate this mock.
//
// The reason is that there's no way to inject build tags into the current
// bazel rules for gomock.
type TestingLogSink = logSink

// sinkOutputOptions provides various options for a logSink.output call.
type sinkOutputOptions struct {
	// extraFlush invites an explicit flush of any buffering.
	extraFlush bool
	// ignoreErrors disables internal error handling (i.e. fail fast).
	ignoreErrors bool
	// forceSync forces synchronous operation of this output operation.
	// That is, it will block until the output has been handled.
	forceSync bool
}

// logSink abstracts the destination of logging events, after all
// their details have been collected into a logpb.Entry.
//
// Each logger can have zero or more logSinks attached to it.
type logSink interface {
	// active returns true if this sink is currently active.
	active() bool

	// attachHints attaches some hints about the location of the message
	// to the stack message.
	attachHints([]byte) []byte

	// output emits some formatted bytes to this sink.
	// the sink is invited to perform an extra flush if indicated
	// by the argument. This is set to true for e.g. Fatal
	// entries.
	//
	// The parent logger's outputMu is held during this operation: log
	// sinks must not recursively call into logging when implementing
	// this method.
	output(b []byte, opts sinkOutputOptions) error

	// exitCode returns the exit code to use if the logger decides
	// to terminate because of an error in output().
	exitCode() exit.Code

	// emergencyOutput attempts to emit some formatted bytes, and
	// ignores any errors.
	//
	// The parent logger's outputMu is held during this operation: log
	// sinks must not recursively call into logging when implementing
	// this method.
	// emergencyOutput([]byte)
}

var _ logSink = (*stderrSink)(nil)
var _ logSink = (*fileSink)(nil)
var _ logSink = (*fluentSink)(nil)
var _ logSink = (*httpSink)(nil)
var _ logSink = (*bufferedSink)(nil)
