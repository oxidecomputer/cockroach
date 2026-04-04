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

package clicfg

// Context represents configuration common to all CLI commands.
//
// Note: when adding new configuration fields, consider adding
// them to the most specific configuration context possible.
// This one is not a dumping ground.
type Context struct {
	// IsInteractive indicates whether the session is interactive, that
	// is, the commands executed are extremely likely to be *input* from
	// a human user: the standard input is a terminal and `-e` was not
	// used (the shell has a prompt).
	//
	// Refer to cli/README.md to understand the general design
	// guidelines for CLI utilities with interactive vs non-interactive
	// input.
	IsInteractive bool

	// EmbeddedMode, when set, reduces the amount of informational
	// messages printed out to exclude details that are not under user's
	// control when the client command is run by a playground
	// environment.
	EmbeddedMode bool
}
