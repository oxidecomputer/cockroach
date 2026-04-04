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

package descpb

import (
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/errors"
)

// PrettyString returns the locking strength as a user-readable string.
func (s ScanLockingStrength) PrettyString() string {
	switch s {
	case ScanLockingStrength_FOR_NONE:
		return "for none"
	case ScanLockingStrength_FOR_KEY_SHARE:
		return "for key share"
	case ScanLockingStrength_FOR_SHARE:
		return "for share"
	case ScanLockingStrength_FOR_NO_KEY_UPDATE:
		return "for no key update"
	case ScanLockingStrength_FOR_UPDATE:
		return "for update"
	default:
		panic(errors.AssertionFailedf("unexpected strength"))
	}
}

// ToScanLockingStrength converts a tree.LockingStrength to its corresponding
// ScanLockingStrength.
func ToScanLockingStrength(s tree.LockingStrength) ScanLockingStrength {
	switch s {
	case tree.ForNone:
		return ScanLockingStrength_FOR_NONE
	case tree.ForKeyShare:
		return ScanLockingStrength_FOR_KEY_SHARE
	case tree.ForShare:
		return ScanLockingStrength_FOR_SHARE
	case tree.ForNoKeyUpdate:
		return ScanLockingStrength_FOR_NO_KEY_UPDATE
	case tree.ForUpdate:
		return ScanLockingStrength_FOR_UPDATE
	default:
		panic(errors.AssertionFailedf("unknown locking strength %s", s))
	}
}

// PrettyString returns the locking strength as a user-readable string.
func (wp ScanLockingWaitPolicy) PrettyString() string {
	switch wp {
	case ScanLockingWaitPolicy_BLOCK:
		return "block"
	case ScanLockingWaitPolicy_SKIP:
		return "skip locked"
	case ScanLockingWaitPolicy_ERROR:
		return "nowait"
	default:
		panic(errors.AssertionFailedf("unexpected wait policy"))
	}
}

// ToScanLockingWaitPolicy converts a tree.LockingWaitPolicy to its
// corresponding ScanLockingWaitPolicy.
func ToScanLockingWaitPolicy(wp tree.LockingWaitPolicy) ScanLockingWaitPolicy {
	switch wp {
	case tree.LockWaitBlock:
		return ScanLockingWaitPolicy_BLOCK
	case tree.LockWaitSkip:
		return ScanLockingWaitPolicy_SKIP
	case tree.LockWaitError:
		return ScanLockingWaitPolicy_ERROR
	default:
		panic(errors.AssertionFailedf("unknown locking wait policy %s", wp))
	}
}
