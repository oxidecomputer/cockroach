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

// AlterRole represents an `ALTER ROLE ... WITH options` statement.
type AlterRole struct {
	Name      RoleSpec
	IfExists  bool
	IsRole    bool
	KVOptions KVOptions
}

// Format implements the NodeFormatter interface.
func (node *AlterRole) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER")
	if node.IsRole {
		ctx.WriteString(" ROLE ")
	} else {
		ctx.WriteString(" USER ")
	}
	if node.IfExists {
		ctx.WriteString("IF EXISTS ")
	}
	ctx.FormatNode(&node.Name)

	if len(node.KVOptions) > 0 {
		ctx.WriteString(" WITH")
		node.KVOptions.formatAsRoleOptions(ctx)
	}
}

// AlterRoleSet represents an `ALTER ROLE ... SET` statement.
type AlterRoleSet struct {
	RoleName     RoleSpec
	IfExists     bool
	IsRole       bool
	AllRoles     bool
	DatabaseName Name
	SetOrReset   *SetVar
}

// Format implements the NodeFormatter interface.
func (node *AlterRoleSet) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER")
	if node.IsRole {
		ctx.WriteString(" ROLE ")
	} else {
		ctx.WriteString(" USER ")
	}
	if node.IfExists {
		ctx.WriteString("IF EXISTS ")
	}
	if node.AllRoles {
		ctx.WriteString("ALL ")
	} else {
		ctx.FormatNode(&node.RoleName)
		ctx.WriteString(" ")
	}
	if node.DatabaseName != "" {
		ctx.WriteString("IN DATABASE ")
		ctx.FormatNode(&node.DatabaseName)
		ctx.WriteString(" ")
	}
	ctx.FormatNode(node.SetOrReset)
}
