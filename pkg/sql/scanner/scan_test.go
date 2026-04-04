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

package scanner

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/lexbase"
	"github.com/stretchr/testify/require"
)

func TestHasMultipleStatements(t *testing.T) {
	testCases := []struct {
		in       string
		expected bool
	}{
		{`a b c`, false},
		{`a; b c`, true},
		{`a b; b c`, true},
		{`a b; b c;`, true},
		{`a b;`, false},
		{`SELECT 123; SELECT 123`, true},
		{`SELECT 123; SELECT 123;`, true},
	}

	for _, tc := range testCases {
		actual, err := HasMultipleStatements(tc.in)
		if err != nil {
			t.Error(err)
		}

		if actual != tc.expected {
			t.Errorf("%q: expected %v, got %v", tc.in, tc.expected, actual)
		}
	}
}

func TestFirstLexicalToken(t *testing.T) {
	tests := []struct {
		s   string
		res int
	}{
		{
			s:   "",
			res: 0,
		},
		{
			s:   " /* comment */ ",
			res: 0,
		},
		{
			s:   "SELECT",
			res: lexbase.SELECT,
		},
		{
			s:   "SELECT 1",
			res: lexbase.SELECT,
		},
		{
			s:   "SELECT 1;",
			res: lexbase.SELECT,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			tok := FirstLexicalToken(tc.s)
			require.Equal(t, tc.res, tok)
		})
	}
}

func TestLastLexicalToken(t *testing.T) {
	tests := []struct {
		s   string
		res int
	}{
		{
			s:   "",
			res: 0,
		},
		{
			s:   " /* comment */ ",
			res: 0,
		},
		{
			s:   "SELECT",
			res: lexbase.SELECT,
		},
		{
			s:   "SELECT 1",
			res: lexbase.ICONST,
		},
		{
			s:   "SELECT 1;",
			res: ';',
		},
		{
			s:   "SELECT 1; /* comment */",
			res: ';',
		},
		{
			s: `SELECT 1;
			    -- comment`,
			res: ';',
		},
		{
			s: `
				--SELECT 1, 2, 3;
				SELECT 4, 5
				--blah`,
			res: lexbase.ICONST,
		},
		{
			s: `
				--SELECT 1, 2, 3;
				SELECT 4, 5;
				--blah`,
			res: ';',
		},
		{
			s:   `SELECT 'unfinished`,
			res: lexbase.ERROR,
		},
		{
			s:   `SELECT e'\xaa';`, // invalid token but last token is semicolon
			res: ';',
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			tok, ok := LastLexicalToken(tc.s)
			if !ok && tok != 0 {
				t.Fatalf("!ok but nonzero tok")
			}
			if tc.res != tok {
				t.Errorf("expected %d but got %d", tc.res, tok)
			}
		})
	}
}
