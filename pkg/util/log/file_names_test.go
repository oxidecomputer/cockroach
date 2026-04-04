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

package log

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/stretchr/testify/require"
)

func TestShortHostName(t *testing.T) {
	defer leaktest.AfterTest(t)()

	testCases := []struct {
		input  string
		output string
	}{
		{"abc", "abc"},
		{"www.example.com", "www"},
		{"", ""},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.output, shortHostname(tc.input))
	}
}

func TestNormalizeFileName(t *testing.T) {
	defer leaktest.AfterTest(t)()

	testCases := []struct {
		input                string
		outputWithHyphens    string
		outputWithoutHyphens string
	}{
		{"abc", "abc", "abc"},
		{"", "", ""},
		{"...", "", ""},
		{"www.example.com", "wwwexamplecom", "wwwexamplecom"},
		{"my-big/test", "my-bigtest", "mybigtest"},
		{"ελλάδα-☃︎..☀️", "ελλάδα-", "ελλάδα"},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.outputWithHyphens, normalizeFileName(tc.input, true))
		require.Equal(t, tc.outputWithoutHyphens, normalizeFileName(tc.input, false))
	}
}
