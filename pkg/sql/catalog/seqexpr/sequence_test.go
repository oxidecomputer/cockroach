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

package seqexpr

import (
	"context"
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

func TestGetSequenceFromFunc(t *testing.T) {
	testData := []struct {
		expr     string
		expected *SeqIdentifier
	}{
		{`nextval('seq')`, &SeqIdentifier{SeqName: "seq"}},
		{`nextval(123::REGCLASS)`, &SeqIdentifier{SeqID: 123}},
		{`nextval(123)`, &SeqIdentifier{SeqID: 123}},
		{`nextval(123::OID::REGCLASS)`, &SeqIdentifier{SeqID: 123}},
		{`nextval(123::OID)`, &SeqIdentifier{SeqID: 123}},
	}

	ctx := context.Background()
	for i, test := range testData {
		t.Run(fmt.Sprintf("%d %s", i, test.expr), func(t *testing.T) {
			parsedExpr, err := parser.ParseExpr(test.expr)
			if err != nil {
				t.Fatal(err)
			}
			semaCtx := tree.MakeSemaContext()
			typedExpr, err := tree.TypeCheck(ctx, parsedExpr, &semaCtx, types.Any)
			if err != nil {
				t.Fatal(err)
			}
			funcExpr, ok := typedExpr.(*tree.FuncExpr)
			if !ok {
				t.Fatal("Expr is not a FuncExpr")
			}
			identifier, err := GetSequenceFromFunc(funcExpr)
			if err != nil {
				t.Fatal(err)
			}
			if identifier.IsByID() {
				if identifier.SeqID != test.expected.SeqID {
					t.Fatalf("expected %d, got %d", test.expected.SeqID, identifier.SeqID)
				}
			} else {
				if identifier.SeqName != test.expected.SeqName {
					t.Fatalf("expected %s, got %s", test.expected.SeqName, identifier.SeqName)
				}
			}
		})
	}
}

func TestGetUsedSequences(t *testing.T) {
	testData := []struct {
		expr     string
		expected []SeqIdentifier
	}{
		{`nextval('seq')`, []SeqIdentifier{
			{SeqName: "seq"},
		}},
		{`nextval(123::REGCLASS)`, []SeqIdentifier{
			{SeqID: 123},
		}},
		{`nextval(123::REGCLASS) + nextval('seq')`, []SeqIdentifier{
			{SeqID: 123},
			{SeqName: "seq"},
		}},
	}

	ctx := context.Background()
	for i, test := range testData {
		t.Run(fmt.Sprintf("%d %s", i, test.expr), func(t *testing.T) {
			parsedExpr, err := parser.ParseExpr(test.expr)
			if err != nil {
				t.Fatal(err)
			}
			semaCtx := tree.MakeSemaContext()
			typedExpr, err := tree.TypeCheck(ctx, parsedExpr, &semaCtx, types.Any)
			if err != nil {
				t.Fatal(err)
			}
			identifiers, err := GetUsedSequences(typedExpr)
			if err != nil {
				t.Fatal(err)
			}

			if len(identifiers) != len(test.expected) {
				t.Fatalf("expected %d identifiers, got %d", len(test.expected), len(identifiers))
			}

			for i, identifier := range identifiers {
				if identifier.IsByID() {
					if identifier.SeqID != test.expected[i].SeqID {
						t.Fatalf("expected %d, got %d", test.expected[i].SeqID, identifier.SeqID)
					}
				} else {
					if identifier.SeqName != test.expected[i].SeqName {
						t.Fatalf("expected %s, got %s", test.expected[i].SeqName, identifier.SeqName)
					}
				}
			}
		})
	}
}

func TestReplaceSequenceNamesWithIDs(t *testing.T) {
	namesToID := map[string]int64{
		"seq": 123,
	}

	testData := []struct {
		expr     string
		expected string
	}{
		{`nextval('seq')`, `nextval(123:::REGCLASS)`},
		{`nextval('non_existent')`, `nextval('non_existent')`},
		{`nextval(123::REGCLASS)`, `nextval(123::REGCLASS)`},
		{`nextval(123)`, `nextval(123)`},
	}

	ctx := context.Background()
	for i, test := range testData {
		t.Run(fmt.Sprintf("%d %s", i, test.expr), func(t *testing.T) {
			parsedExpr, err := parser.ParseExpr(test.expr)
			if err != nil {
				t.Fatal(err)
			}
			semaCtx := tree.MakeSemaContext()
			typedExpr, err := tree.TypeCheck(ctx, parsedExpr, &semaCtx, types.Any)
			if err != nil {
				t.Fatal(err)
			}
			newExpr, err := ReplaceSequenceNamesWithIDs(typedExpr, namesToID)
			if err != nil {
				t.Fatal(err)
			}
			if newExpr.String() != test.expected {
				t.Fatalf("expected %s, got %s", test.expected, newExpr.String())
			}
		})
	}
}
