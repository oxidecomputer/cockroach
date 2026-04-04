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

// Package lock provides type definitions for locking-related concepts used by
// concurrency control in the key-value layer.
package lock

import "github.com/cockroachdb/redact"

// SafeFormat implements redact.SafeFormatter.
func (lw Waiter) SafeFormat(w redact.SafePrinter, _ rune) {
	expand := w.Flag('+')

	txnIDRedactableString := redact.Sprint(nil)
	if lw.WaitingTxn != nil {
		if expand {
			txnIDRedactableString = redact.Sprint(lw.WaitingTxn.ID)
		} else {
			txnIDRedactableString = redact.Sprint(lw.WaitingTxn.Short())
		}
	}
	w.Printf("waiting_txn:%s active_waiter:%t strength:%s wait_duration:%s", txnIDRedactableString, lw.ActiveWaiter, lw.Strength, lw.WaitDuration)
}
