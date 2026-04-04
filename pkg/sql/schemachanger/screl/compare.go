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
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/rel"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
)

// equalityAttrs are used to sort elements.
var equalityAttrs = func() []rel.Attr {
	s := make([]rel.Attr, 0, AttrMax)
	s = append(s, rel.Type)
	for a := Attr(1); a <= AttrMax; a++ {
		s = append(s, a)
	}
	return s
}()

// EqualElements returns true if the two elements are equal.
func EqualElements(a, b scpb.Element) bool {
	return Schema.EqualOn(equalityAttrs, a, b)
}

// CompareElements orders two elements.
func CompareElements(a, b scpb.Element) (less, eq bool) {
	return Schema.CompareOn(equalityAttrs, a, b)
}
