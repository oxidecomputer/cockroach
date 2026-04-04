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

// DropOwnedBy represents a DROP OWNED BY command.
type DropOwnedBy struct {
	Roles        RoleSpecList
	DropBehavior DropBehavior
}

var _ Statement = &DropOwnedBy{}

// Format implements the NodeFormatter interface.
func (node *DropOwnedBy) Format(ctx *FmtCtx) {
	ctx.WriteString("DROP OWNED BY ")
	for i := range node.Roles {
		if i > 0 {
			ctx.WriteString(", ")
		}
		node.Roles[i].Format(ctx)
	}
	if node.DropBehavior != DropDefault {
		ctx.WriteString(" ")
		ctx.WriteString(node.DropBehavior.String())
	}
}
