// Copyright 2016 The Cockroach Authors.
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

package util

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSingleRune(t *testing.T) {
	tests := []struct {
		s        string
		expected rune
		err      bool
	}{
		{"a", 'a', false},
		{"", 0, false},
		{"🐛"[:1], 0, true},
		{"aa", 'a', true},
	}
	for _, tc := range tests {
		t.Run(tc.s, func(t *testing.T) {
			got, err := GetSingleRune(tc.s)
			if (err != nil) != tc.err {
				t.Fatalf("got unexpected err: %v", err)
			}
			if tc.expected != got {
				t.Fatalf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestToLowerSingleByte(t *testing.T) {
	testCases := []struct {
		from     byte
		expected byte
	}{
		{'a', 'a'},
		{'A', 'a'},
		{'c', 'c'},
		{'C', 'c'},
		{'Z', 'z'},
		{'1', '1'},
		{'\n', '\n'},
	}

	for _, tc := range testCases {
		t.Run(string(tc.from), func(t *testing.T) {
			ret := ToLowerSingleByte(tc.from)
			require.Equal(t, tc.expected, ret)
		})
	}
}

func TestTruncateString(t *testing.T) {
	testCases := []struct {
		s string
		// res stores the expected results for maxRunes=0,1,2,3,etc.
		res []string
	}{
		{"", []string{"", ""}},
		{"abcd", []string{"", "a", "ab", "abc", "abcd", "abcd", "abcd"}},
		{"🐛🏠", []string{"", "🐛", "🐛🏠", "🐛🏠", "🐛🏠"}},
		{"a🐛b🏠c", []string{"", "a", "a🐛", "a🐛b", "a🐛b🏠", "a🐛b🏠c", "a🐛b🏠c"}},
		{
			// Test with an invalid UTF-8 sequence.
			"\xf0\x90\x28\xbc",
			[]string{"", "\xf0", "\xf0\x90", "\xf0\x90\x28", "\xf0\x90\x28\xbc", "\xf0\x90\x28\xbc"},
		},
	}

	for _, tc := range testCases {
		for i := range tc.res {
			if r := TruncateString(tc.s, i); r != tc.res[i] {
				t.Errorf("TruncateString(\"%q\", %d) = \"%q\"; expected \"%q\"", tc.s, i, r, tc.res[i])
			}
		}
	}
}

func TestRemoveTrailingSpaces(t *testing.T) {
	for _, tc := range []struct{ input, expected string }{
		{
			input:    "",
			expected: "\n",
		},
		{
			input:    "line 1  \nline 2   \nline 3 \n\n",
			expected: "line 1\nline 2\nline 3\n",
		},
		{
			input:    " line 1  \nline 2   \nline 3  ",
			expected: " line 1\nline 2\nline 3\n",
		},
		{
			input:    "line 1\n\n  \nline 2   \nline 3",
			expected: "line 1\n\n\nline 2\nline 3\n",
		},
	} {
		output := RemoveTrailingSpaces(tc.input)
		if output != tc.expected {
			t.Errorf("expected:\n%s\ngot:\n%s", tc.expected, output)
		}
	}
}

func TestStringListBuilder(t *testing.T) {
	var buf bytes.Buffer
	var b StringListBuilder
	expect := func(exp string) {
		t.Helper()
		if buf.String() != exp {
			t.Errorf("expected `%s`, got `%s`", exp, buf.String())
		}
		buf.Reset()
	}

	b = MakeStringListBuilder("(", ",", ")")
	b.Finish(&buf)
	expect("")

	b = MakeStringListBuilder("(", ",", ")")
	b.Add(&buf, "one")
	b.Finish(&buf)
	expect("(one)")

	b = MakeStringListBuilder("[", ", ", "]")
	b.Add(&buf, "one")
	b.Addf(&buf, "%s", "two")
	b.Finish(&buf)
	expect("[one, two]")
}
