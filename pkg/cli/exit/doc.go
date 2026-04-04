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

// Package exit encapsulates calls to os.Exit to control the
// production of process exit status codes.
//
// Its goal is to ensure that all possible exit codes produced
// by the 'cockroach' process upon termination are documented.
// It achieves this by providing a type exit.Code and requiring that
// all possible values come from constructors in the package (see
// codes.go). A linter ensures that no direct call to os.Exit() can be
// present elsewhere.
//
// Note that due to the limited range of unix exit codes, it is not
// possible to map all possible error situations inside a CockroachDB
// server to a unique exit code.
// This is why the main mechanism to explain the cause of a process
// termination must remain the logging subsystem.
//
// The introduction of discrete exit codes here is thus meant to
// merely complement logging, in those cases where logging is unable
// to detail the reason why the process is terminating; for example:
//
// - before logging is initialized (e.g. during command-line parsing)
// - when a logging operation fails.
//
// For client commands, the situation is different: there are much
// fewer different exit situations, so we could envision discrete
// error codes for them. Additionally, different client commands
// can reuse the same numeric codes for different error situations,
// when they do not overlap.
//
// This package accommodates this as follows:
//
// - exit codes common to all commands should be allocated
//   incrementally starting from the last defined common error
//   in codes.go.
//
// - exit codes specific to one command should be allocated downwards
//   starting from 125.
//
package exit
