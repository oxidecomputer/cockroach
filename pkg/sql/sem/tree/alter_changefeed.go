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

package tree

// AlterChangefeed represents an ALTER CHANGEFEED statement.
type AlterChangefeed struct {
	Jobs Expr
	Cmds AlterChangefeedCmds
}

var _ Statement = &AlterChangefeed{}

// Format implements the NodeFormatter interface.
func (node *AlterChangefeed) Format(ctx *FmtCtx) {
	ctx.WriteString(`ALTER CHANGEFEED `)
	ctx.FormatNode(node.Jobs)
	ctx.FormatNode(&node.Cmds)
}

// AlterChangefeedCmds represents a list of changefeed alterations
type AlterChangefeedCmds []AlterChangefeedCmd

// Format implements the NodeFormatter interface.
func (node *AlterChangefeedCmds) Format(ctx *FmtCtx) {
	for i, n := range *node {
		if i > 0 {
			ctx.WriteString(" ")
		}
		ctx.FormatNode(n)
	}
}

// AlterChangefeedCmd represents a changefeed modification operation.
type AlterChangefeedCmd interface {
	NodeFormatter
	// Placeholder function to ensure that only desired types
	// (AlterChangefeed*) conform to the AlterChangefeedCmd interface.
	alterChangefeedCmd()
}

func (*AlterChangefeedAddTarget) alterChangefeedCmd()    {}
func (*AlterChangefeedDropTarget) alterChangefeedCmd()   {}
func (*AlterChangefeedSetOptions) alterChangefeedCmd()   {}
func (*AlterChangefeedUnsetOptions) alterChangefeedCmd() {}

var _ AlterChangefeedCmd = &AlterChangefeedAddTarget{}
var _ AlterChangefeedCmd = &AlterChangefeedDropTarget{}
var _ AlterChangefeedCmd = &AlterChangefeedSetOptions{}
var _ AlterChangefeedCmd = &AlterChangefeedUnsetOptions{}

// AlterChangefeedAddTarget represents an ADD <targets> command
type AlterChangefeedAddTarget struct {
	Targets ChangefeedTargets
	Options KVOptions
}

// Format implements the NodeFormatter interface.
func (node *AlterChangefeedAddTarget) Format(ctx *FmtCtx) {
	ctx.WriteString(" ADD ")
	ctx.FormatNode(&node.Targets)
	if node.Options != nil {
		ctx.WriteString(" WITH ")
		ctx.FormatNode(&node.Options)
	}
}

// AlterChangefeedDropTarget represents an DROP <targets> command
type AlterChangefeedDropTarget struct {
	Targets ChangefeedTargets
}

// Format implements the NodeFormatter interface.
func (node *AlterChangefeedDropTarget) Format(ctx *FmtCtx) {
	ctx.WriteString(" DROP ")
	ctx.FormatNode(&node.Targets)
}

// AlterChangefeedSetOptions represents an SET <options> command
type AlterChangefeedSetOptions struct {
	Options KVOptions
}

// Format implements the NodeFormatter interface.
func (node *AlterChangefeedSetOptions) Format(ctx *FmtCtx) {
	ctx.WriteString(" SET ")
	ctx.FormatNode(&node.Options)
}

// AlterChangefeedUnsetOptions represents an UNSET <options> command
type AlterChangefeedUnsetOptions struct {
	Options NameList
}

// Format implements the NodeFormatter interface.
func (node *AlterChangefeedUnsetOptions) Format(ctx *FmtCtx) {
	ctx.WriteString(" UNSET ")
	ctx.FormatNode(&node.Options)
}
