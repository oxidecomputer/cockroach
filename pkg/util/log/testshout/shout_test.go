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

package testshout

import (
	"context"
	"os"

	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/log/logconfig"
	"github.com/cockroachdb/cockroach/pkg/util/log/severity"
)

// Example_shout_before_log verifies that Shout output emitted after
// the log flags were set, but before the first log message was
// output, properly appears on stderr.
//
// This test needs to occur in its own test package where there is no
// other activity on the log flags, and no other log activity,
// otherwise the test's behavior will break on `make stress`.
func Example_shout_before_log() {
	// Set up a configuration where only WARNING or above goes to stderr.
	cfg := logconfig.DefaultConfig()
	if err := cfg.Validate(nil /* no dir */); err != nil {
		panic(err)
	}
	cfg.Sinks.Stderr.Filter = severity.WARNING
	cleanup, err := log.ApplyConfig(cfg)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Redirect stderr to stdout so the reference output checking below
	// has something to work with.
	origStderr := log.OrigStderr
	log.OrigStderr = os.Stdout
	defer func() { log.OrigStderr = origStderr }()

	log.Shout(context.Background(), severity.INFO, "hello world")

	// output:
	// *
	// * INFO: hello world
	// *
}
