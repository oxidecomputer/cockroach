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

package scpb

const (
	// PlaceHolderComment placeholder string for non-fetched comments.
	PlaceHolderComment string = "__placeholder_comment__"
	// PlaceHolderRoleName placeholder string for non-fetched role names.
	PlaceHolderRoleName string = "__placeholder_role_name__"
)

// TargetStatus is the subset of Status values which serve as valid target
// statuses.
type TargetStatus Status

const (
	// InvalidTarget indicates that the element doesn't have a target status
	// for the current schema change.
	InvalidTarget TargetStatus = TargetStatus(Status_UNKNOWN)

	// ToAbsent indicates that the element should have status ABSENT after
	// completion of the current schema change, meaning it should no longer exist.
	ToAbsent TargetStatus = TargetStatus(Status_ABSENT)

	// ToPublic indicates that the element should have status PUBLIC after
	// completion of the current schema change, meaning it should exist and not
	// in any kind of transitory state.
	ToPublic TargetStatus = TargetStatus(Status_PUBLIC)

	// Transient is like ToAbsent in that the element should no longer exist once
	// the schema change is done. The element is also not assumed to exist prior
	// to the schema change, so this target status is used to ensure it comes into
	// existence before disappearing again. Otherwise, an element whose current
	// and target statuses are both ABSENT won't experience any state transitions.
	Transient TargetStatus = TargetStatus(Status_TRANSIENT)
)

// Status returns the TargetStatus as a Status.
func (t TargetStatus) Status() Status {
	return Status(t)
}

// AsTargetStatus returns the Status as a TargetStatus.
func AsTargetStatus(s Status) TargetStatus {
	switch s {
	case Status_ABSENT:
		return ToAbsent
	case Status_PUBLIC:
		return ToPublic
	case Status_TRANSIENT:
		return Transient
	default:
		return InvalidTarget
	}
}
