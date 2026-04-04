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

package typedesc

import (
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/redact"
)

// SafeMessage makes immutable a SafeMessager.
func (desc *immutable) SafeMessage() string {
	return formatSafeType("typedesc.immutable", desc)
}

// SafeMessage makes Mutable a SafeMessager.
func (desc *Mutable) SafeMessage() string {
	return formatSafeType("typedesc.Mutable", desc)
}

func formatSafeType(typeName string, desc catalog.TypeDescriptor) string {
	var buf redact.StringBuilder
	buf.Printf(typeName + ": {")
	formatSafeTypeProperties(&buf, desc)
	buf.Printf("}")
	return buf.String()
}

func formatSafeTypeProperties(w *redact.StringBuilder, desc catalog.TypeDescriptor) {
	catalog.FormatSafeDescriptorProperties(w, desc)
	td := desc.TypeDesc()
	w.Printf(", Kind: %s", td.Kind)
	if len(td.EnumMembers) > 0 {
		w.Printf(", NumEnumMembers: %d", len(td.EnumMembers))
	}
	if td.Alias != nil {
		w.Printf(", Alias: %d", td.Alias.Oid())
	}
	if td.ArrayTypeID != 0 {
		w.Printf(", ArrayTypeID: %d", td.ArrayTypeID)
	}
	for i := range td.ReferencingDescriptorIDs {
		w.Printf(", ")
		if i == 0 {
			w.Printf("ReferencingDescriptorIDs: [")
		}
		w.Printf("%d", td.ReferencingDescriptorIDs[i])
	}
	if len(td.ReferencingDescriptorIDs) > 0 {
		w.Printf("]")
	}
}
