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
	"github.com/cockroachdb/cockroach/pkg/cli/exit"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
)

// Type of a stderr copy sink.
type stderrSink struct {
	// the --no-color flag. When set it disables escapes code on the
	// stderr copy.
	noColor syncutil.AtomicBool
}

// activeAtSeverity implements the logSink interface.
func (l *stderrSink) active() bool { return true }

// attachHints implements the logSink interface.
func (l *stderrSink) attachHints(stacks []byte) []byte {
	return stacks
}

// output implements the logSink interface.
func (l *stderrSink) output(b []byte, _ sinkOutputOptions) error {
	_, err := OrigStderr.Write(b)
	return err
}

// exitCode implements the logSink interface.
func (l *stderrSink) exitCode() exit.Code {
	return exit.LoggingStderrUnavailable()
}
