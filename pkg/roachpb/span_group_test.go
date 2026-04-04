// Copyright 2019 The Cockroach Authors.
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

package roachpb

import (
	"reflect"
	"testing"
)

func TestSpanGroup(t *testing.T) {
	g := &SpanGroup{}

	be := makeSpan("b-e")

	if g.Contains(Key("a")) {
		t.Fatal("empty group should not contain a")
	}
	if g.Sub(be) {
		t.Fatalf("removing b-e from empty group should not expand it")
	}
	if !g.Add(be) {
		t.Fatal("adding b-e to empty should expand it")
	}
	if !g.Sub(be) {
		t.Fatalf("removing b-e from b-e should expand it")
	}
	g.Add(be)
	if g.Add(be) {
		t.Fatal("adding  b-e to b-e should not expand it")
	}
	if g.Add(makeSpan("c-d")) {
		t.Fatal("adding c-d to b-e should not expand it")
	}
	if g.Add(makeSpan("b-d")) {
		t.Fatal("adding b-d to b-e should not expand it")
	}
	if g.Add(makeSpan("d-e")) {
		t.Fatal("adding d-e to b-e should not expand it")
	}
	if got, expected := g.Len(), 1; got != expected {
		t.Fatalf("got %d, expected %d", got, expected)
	}
	if got, expected := g.Slice(), be; len(got) != 1 || !reflect.DeepEqual(got[0], expected) {
		t.Fatalf("got %v, expected %v", got, expected)
	}
	for _, k := range []string{"b", "c", "d"} {
		if !g.Contains(Key(k)) {
			t.Fatalf("span b-e should contain %q", k)
		}
	}
	for _, k := range []string{"a", "e", "f"} {
		if g.Contains(Key(k)) {
			t.Fatalf("span b-e should not contain %q", k)
		}
	}
	if !g.Encloses(makeSpan("b-d")) {
		t.Fatalf("span b-e should enclose b-d")
	}
	if g.Encloses(makeSpan("b-d"), makeSpan("c-f")) {
		t.Fatalf("span b-e should not enclose c-f")
	}
	if g.Encloses(makeSpan("b-d"), makeSpan("a-f")) {
		t.Fatalf("span b-e should not enclose a-f")
	}
	if g.Encloses(makeSpan("b-d"), makeSpan("e-f")) {
		t.Fatalf("span b-e should not enclose e-f")
	}
	if !g.Sub(makeSpan("d-e"), makeSpan("b-c")) {
		t.Fatalf("removing b-c and d-e from b-e should expand it")
	}
	if got, expected := g.Slice(), makeSpan("c-d"); len(got) != 1 || !reflect.DeepEqual(got[0], expected) {
		t.Fatalf("got %v, expected %v", got, expected)
	}
}
