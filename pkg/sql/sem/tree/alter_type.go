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

import (
	"github.com/cockroachdb/cockroach/pkg/server/telemetry"
	"github.com/cockroachdb/cockroach/pkg/sql/sqltelemetry"
)

// AlterType represents an ALTER TYPE statement.
type AlterType struct {
	Type *UnresolvedObjectName
	Cmd  AlterTypeCmd
}

// Format implements the NodeFormatter interface.
func (node *AlterType) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER TYPE ")
	ctx.FormatNode(node.Type)
	ctx.FormatNode(node.Cmd)
}

// AlterTypeCmd represents a type modification operation.
type AlterTypeCmd interface {
	NodeFormatter
	alterTypeCmd()
	// TelemetryCounter returns the telemetry counter to increment
	// when this command is used.
	TelemetryCounter() telemetry.Counter
}

func (*AlterTypeAddValue) alterTypeCmd()    {}
func (*AlterTypeRenameValue) alterTypeCmd() {}
func (*AlterTypeRename) alterTypeCmd()      {}
func (*AlterTypeSetSchema) alterTypeCmd()   {}
func (*AlterTypeOwner) alterTypeCmd()       {}
func (*AlterTypeDropValue) alterTypeCmd()   {}

var _ AlterTypeCmd = &AlterTypeAddValue{}
var _ AlterTypeCmd = &AlterTypeRenameValue{}
var _ AlterTypeCmd = &AlterTypeRename{}
var _ AlterTypeCmd = &AlterTypeSetSchema{}
var _ AlterTypeCmd = &AlterTypeOwner{}
var _ AlterTypeCmd = &AlterTypeDropValue{}

// AlterTypeAddValue represents an ALTER TYPE ADD VALUE command.
type AlterTypeAddValue struct {
	NewVal      EnumValue
	IfNotExists bool
	Placement   *AlterTypeAddValuePlacement
}

// Format implements the NodeFormatter interface.
func (node *AlterTypeAddValue) Format(ctx *FmtCtx) {
	ctx.WriteString(" ADD VALUE ")
	if node.IfNotExists {
		ctx.WriteString("IF NOT EXISTS ")
	}
	ctx.FormatNode(&node.NewVal)
	if node.Placement != nil {
		if node.Placement.Before {
			ctx.WriteString(" BEFORE ")
		} else {
			ctx.WriteString(" AFTER ")
		}
		ctx.FormatNode(&node.Placement.ExistingVal)
	}
}

// TelemetryCounter implements the AlterTypeCmd interface.
func (node *AlterTypeAddValue) TelemetryCounter() telemetry.Counter {
	return sqltelemetry.SchemaChangeAlterCounterWithExtra("type", "add_value")
}

// AlterTypeAddValuePlacement represents the placement clause for an ALTER
// TYPE ADD VALUE command ([BEFORE | AFTER] value).
type AlterTypeAddValuePlacement struct {
	Before      bool
	ExistingVal EnumValue
}

// AlterTypeRenameValue represents an ALTER TYPE RENAME VALUE command.
type AlterTypeRenameValue struct {
	OldVal EnumValue
	NewVal EnumValue
}

// Format implements the NodeFormatter interface.
func (node *AlterTypeRenameValue) Format(ctx *FmtCtx) {
	ctx.WriteString(" RENAME VALUE ")
	ctx.FormatNode(&node.OldVal)
	ctx.WriteString(" TO ")
	ctx.FormatNode(&node.NewVal)
}

// TelemetryCounter implements the AlterTypeCmd interface.
func (node *AlterTypeRenameValue) TelemetryCounter() telemetry.Counter {
	return sqltelemetry.SchemaChangeAlterCounterWithExtra("type", "rename_value")
}

// AlterTypeDropValue represents an ALTER TYPE DROP VALUE command.
type AlterTypeDropValue struct {
	Val EnumValue
}

// Format implements the NodeFormatter interface.
func (node *AlterTypeDropValue) Format(ctx *FmtCtx) {
	ctx.WriteString(" DROP VALUE ")
	ctx.FormatNode(&node.Val)
}

// TelemetryCounter implements the AlterTypeCmd interface.
func (node *AlterTypeDropValue) TelemetryCounter() telemetry.Counter {
	return sqltelemetry.SchemaChangeAlterCounterWithExtra("type", "drop_value")
}

// AlterTypeRename represents an ALTER TYPE RENAME command.
type AlterTypeRename struct {
	NewName Name
}

// Format implements the NodeFormatter interface.
func (node *AlterTypeRename) Format(ctx *FmtCtx) {
	ctx.WriteString(" RENAME TO ")
	ctx.FormatNode(&node.NewName)
}

// TelemetryCounter implements the AlterTypeCmd interface.
func (node *AlterTypeRename) TelemetryCounter() telemetry.Counter {
	return sqltelemetry.SchemaChangeAlterCounterWithExtra("type", "rename")
}

// AlterTypeSetSchema represents an ALTER TYPE SET SCHEMA command.
type AlterTypeSetSchema struct {
	Schema Name
}

// Format implements the NodeFormatter interface.
func (node *AlterTypeSetSchema) Format(ctx *FmtCtx) {
	ctx.WriteString(" SET SCHEMA ")
	ctx.FormatNode(&node.Schema)
}

// TelemetryCounter implements the AlterTypeCmd interface.
func (node *AlterTypeSetSchema) TelemetryCounter() telemetry.Counter {
	return sqltelemetry.SchemaChangeAlterCounterWithExtra("type", "set_schema")
}

// AlterTypeOwner represents an ALTER TYPE OWNER TO command.
type AlterTypeOwner struct {
	Owner RoleSpec
}

// Format implements the NodeFormatter interface.
func (node *AlterTypeOwner) Format(ctx *FmtCtx) {
	ctx.WriteString(" OWNER TO ")
	ctx.FormatNode(&node.Owner)
}

// TelemetryCounter implements the AlterTypeCmd interface.
func (node *AlterTypeOwner) TelemetryCounter() telemetry.Counter {
	return sqltelemetry.SchemaChangeAlterCounterWithExtra("type", "owner")
}
