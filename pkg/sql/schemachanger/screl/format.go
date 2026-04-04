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

package screl

import (
	"fmt"
	"reflect"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/rel"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/redact"
)

// NodeString formats a node as a string by invoking FormatNode.
func NodeString(n *Node) string {
	var v redact.StringBuilder
	if err := FormatNode(&v, n); err != nil {
		return fmt.Sprintf("failed for format node: %v", err)
	}
	return v.String()
}

// FormatNode formats the node into the SafeWriter.
func FormatNode(w redact.SafeWriter, e *Node) (err error) {
	w.SafeString("[[")
	if err := FormatElement(w, e.Element()); err != nil {
		return err
	}
	w.SafeString(", ")
	w.SafeString(redact.SafeString(e.Target.TargetStatus.String()))
	w.SafeString("], ")
	w.SafeString(redact.SafeString(e.CurrentStatus.String()))
	w.SafeString("]")
	return nil
}

// ElementString formats an element as a string by invoking FormatElement.
func ElementString(e scpb.Element) string {
	var v redact.StringBuilder
	if err := FormatElement(&v, e); err != nil {
		return fmt.Sprintf("failed for format element %T: %v", e, err)
	}
	return v.String()
}

// FormatElement formats the element into the SafeWriter.
func FormatElement(w redact.SafeWriter, e scpb.Element) (err error) {
	if e == nil {
		return errors.Errorf("nil element")
	}
	w.SafeString(redact.SafeString(reflect.TypeOf(e).Elem().Name()))
	w.SafeString(":{")
	var written int
	if err := Schema.IterateAttributes(e, func(attr rel.Attr, value interface{}) error {
		if written > 0 {
			w.SafeString(", ")
		}
		written++
		// Change the type of strings so that they get quoted appropriately.
		if str, isStr := value.(string); isStr {
			value = tree.Name(str)
		}
		w.SafeString(redact.SafeString(attr.String()))
		w.SafeString(": ")
		w.Printf("%v", value)
		return nil
	}); err != nil {
		return err
	}
	w.SafeRune('}')
	return nil
}
