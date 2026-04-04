// Copyright 2017 The Cockroach Authors.
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

package cli

import (
	"os"

	"github.com/cockroachdb/cockroach/pkg/cli/exit"
)

// drainSignals are the signals that will cause the server to drain and exit.
var drainSignals = []os.Signal{os.Interrupt}

// termSignal is the signal that causes an idempotent graceful
// shutdown (i.e. second occurrence does not incur hard shutdown).
var termSignal os.Signal = nil

// quitSignal is the signal to recognize to dump Go stacks.
var quitSignal os.Signal = nil

// debugSignal is the signal to open a pprof debugging server.
var debugSignal os.Signal = nil

const backgroundFlagDefined = false

func handleSignalDuringShutdown(os.Signal) {
	// Windows doesn't indicate whether a process exited due to a signal in the
	// exit code, so we don't need to do anything but exit with a failing code.
	// The error message has already been printed.
	exit.WithCode(exit.UnspecifiedError())
}

func maybeRerunBackground() (bool, error) {
	return false, nil
}

func disableOtherPermissionBits() {
	// No-op on windows, which does not support umask.
}

func closeAllSockets() {
	// No-op on windows.
	// TODO(jackson): Is there something else we can do on Windows?
}
