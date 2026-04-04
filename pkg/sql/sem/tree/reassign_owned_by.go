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

package tree

// ReassignOwnedBy represents a REASSIGN OWNED BY <name> TO <name> statement.
type ReassignOwnedBy struct {
	OldRoles RoleSpecList
	NewRole  RoleSpec
}

var _ Statement = &ReassignOwnedBy{}

// Format implements the NodeFormatter interface.
func (node *ReassignOwnedBy) Format(ctx *FmtCtx) {
	ctx.WriteString("REASSIGN OWNED BY ")
	for i := range node.OldRoles {
		if i > 0 {
			ctx.WriteString(", ")
		}
		node.OldRoles[i].Format(ctx)
	}
	ctx.WriteString(" TO ")
	ctx.FormatNode(&node.NewRole)
}
