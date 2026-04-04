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
//

// Package sqlfsm contains the definitions for the state labels of the
// conn executor FSM. These are in a separate package to ensure the
// SQL shell only uses a minimal dependency graph.
package sqlfsm

// Constants for the String() representation of the session states. Shared with
// the CLI code which needs to recognize them.
const (
	NoTxnStateStr         = "NoTxn"
	OpenStateStr          = "Open"
	AbortedStateStr       = "Aborted"
	CommitWaitStateStr    = "CommitWait"
	InternalErrorStateStr = "InternalError"
)
