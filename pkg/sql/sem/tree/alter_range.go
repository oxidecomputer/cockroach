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

package tree

import "github.com/cockroachdb/errors"

// RelocateRange represents an `ALTER RANGE .. RELOCATE ..`
// statement.
type RelocateRange struct {
	Rows            *Select
	ToStoreID       Expr
	FromStoreID     Expr
	SubjectReplicas RelocateSubject
}

// RelocateSubject indicates what replicas of a range should be relocated.
type RelocateSubject int8

const (
	// RelocateLease indicates that leases should be relocated.
	RelocateLease RelocateSubject = iota
	// RelocateVoters indicates what voter replicas should be relocated.
	RelocateVoters
	// RelocateNonVoters indicates that non-voter replicas should be relocated.
	RelocateNonVoters
)

// Format implementsthe NodeFormatter interface.
func (n *RelocateSubject) Format(ctx *FmtCtx) {
	ctx.WriteString(n.String())
}

func (n RelocateSubject) String() string {
	switch n {
	case RelocateLease:
		return "LEASE"
	case RelocateVoters:
		return "VOTERS"
	case RelocateNonVoters:
		return "NONVOTERS"
	default:
		panic(errors.AssertionFailedf("programming error: unhandled case %d", int(n)))
	}
}

// Format implements the NodeFormatter interface.
func (n *RelocateRange) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER RANGE RELOCATE ")
	ctx.FormatNode(&n.SubjectReplicas)
	// When relocating leases, the origin store is implicit.
	if n.SubjectReplicas != RelocateLease {
		ctx.WriteString(" FROM ")
		ctx.FormatNode(n.FromStoreID)
	}
	ctx.WriteString(" TO ")
	ctx.FormatNode(n.ToStoreID)
	ctx.WriteString(" FOR ")
	ctx.FormatNode(n.Rows)
}
