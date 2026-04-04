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

// StreamIngestion represents a RESTORE FROM REPLICATION STREAM statement.
type StreamIngestion struct {
	Targets TargetList
	From    StringOrPlaceholderOptList
	AsOf    AsOfClause
}

var _ Statement = &StreamIngestion{}

// Format implements the NodeFormatter interface.
func (node *StreamIngestion) Format(ctx *FmtCtx) {
	ctx.WriteString("RESTORE ")
	ctx.FormatNode(&node.Targets)
	ctx.WriteString(" ")
	ctx.WriteString("FROM REPLICATION STREAM FROM ")
	ctx.FormatNode(&node.From)
	if node.AsOf.Expr != nil {
		ctx.WriteString(" ")
		ctx.FormatNode(&node.AsOf)
	}
}
