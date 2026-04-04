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

import (
	"strconv"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
)

// TenantID represents a tenant ID that can be pretty-printed.
type TenantID struct {
	roachpb.TenantID

	// Specified is set to true when the TENANT clause was specified in
	// the backup target input syntax. We need this, instead of relying
	// on roachpb.TenantID.IsSet(), because we need a special marker for
	// the case when the value was anonymized. In other places, the
	// value used for anonymized integer literals is 0, but we can't use
	// 0 for TenantID as this is refused during parsing, nor can we use
	// any other value since all non-zero integers are valid tenant IDs.
	//
	// TODO(knz): None of this complexity would be needed if the tenant
	// ID was specified using an Expr here. This can be removed once the
	// backup target syntax accepts expressions for tenant targets.
	Specified bool
}

var _ NodeFormatter = (*TenantID)(nil)

// Format implements the NodeFormatter interface.
func (t *TenantID) Format(ctx *FmtCtx) {
	if ctx.flags.HasFlags(FmtHideConstants) || !t.IsSet() {
		// The second part of the condition above is conceptually
		// redundant, but is sadly needed here because we need to be able
		// to re-print a TENANT clause after it has been anonymized,
		// and we can't emit the value zero which is unparseable.
		//
		// TODO(knz): This branch can be removed once the TENANT clause
		// accepts expressions.
		ctx.WriteByte('_')
	} else {
		ctx.WriteString(strconv.FormatUint(t.ToUint64(), 10))
	}
}
