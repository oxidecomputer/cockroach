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

package fuzzystrmatch

import (
	"math/rand"
	"testing"
)

func TestSoundex(t *testing.T) {
	tt := []struct {
		Source   string
		Expected string
	}{
		{
			Source:   "hello world!",
			Expected: "H464",
		},
		{
			Source:   "Anne",
			Expected: "A500",
		},
		{
			Source:   "Ann",
			Expected: "A500",
		},
		{
			Source:   "Andrew",
			Expected: "A536",
		},
		{
			Source:   "Margaret",
			Expected: "M626",
		},
		{
			Source:   "🌞",
			Expected: "",
		},
		{
			Source:   "😄 🐃 🐯 🕣 💲 🏜 👞 🔠 🌟 📌",
			Expected: "",
		},
		{
			Source:   "zażółćx",
			Expected: "Z200",
		},
		{
			Source:   "K😋",
			Expected: "K000",
		},
		// Regression test for #82640, just ensure we don't panic.
		{
			Source:   "l�qă�_��",
			Expected: "L200",
		},
	}

	for _, tc := range tt {
		got := Soundex(tc.Source)
		if tc.Expected != got {
			t.Fatalf("error convert string to its Soundex code with source=%q"+
				" expected %s got %s", tc.Source, tc.Expected, got)
		}
	}

	// Run some random test cases to make sure we don't panic.

	for i := 0; i < 1000; i++ {
		l := rand.Int31n(10)
		b := make([]byte, l)
		//lint:ignore SA1019 "math/rand".Rand is acceptable in test code
		rand.Read(b)

		soundex(string(b))
	}
}

func TestDifference(t *testing.T) {
	tt := []struct {
		Source   string
		Target   string
		Expected int
	}{
		{
			Source:   "Anne",
			Target:   "Ann",
			Expected: 4,
		},
		{
			Source:   "Anne",
			Target:   "Andrew",
			Expected: 2,
		},
		{
			Source:   "Anne",
			Target:   "Margaret",
			Expected: 0,
		},
	}

	for _, tc := range tt {
		got := Difference(tc.Source, tc.Target)
		if tc.Expected != got {
			t.Fatalf("error reports the number of matching code positions with source=%q"+
				" target=%q: expected %d got %d", tc.Source, tc.Target, tc.Expected, got)
		}
	}
}
