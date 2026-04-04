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

package security_test

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
)

func TestUserName(t *testing.T) {
	defer leaktest.AfterTest(t)()

	testCases := []struct {
		username   string
		normalized string
		err        string
	}{
		{"Abc123", "abc123", ""},
		{"0123121132", "0123121132", ""},
		{"HeLlO", "hello", ""},
		{"Ομηρος", "ομηρος", ""},
		{"_HeLlO", "_hello", ""},
		{"a-BC-d", "a-bc-d", ""},
		{"A.Bcd", "a.bcd", ""},
		{"WWW.BIGSITE.NET", "www.bigsite.net", ""},
		{"", "", `username is empty`},
		{"-ABC", "", `username is invalid`},
		{".ABC", "", `username is invalid`},
		{"*.wildcard", "", `username is invalid`},
		{"foofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoof", "", `username is too long`},
		{"M", "m", ""},
		{".", "", `username is invalid`},
	}

	for _, tc := range testCases {
		username, err := security.MakeSQLUsernameFromUserInput(tc.username, security.UsernameCreation)
		if !testutils.IsError(err, tc.err) {
			t.Errorf("%q: expected %q, got %v", tc.username, tc.err, err)
			continue
		}
		if err != nil {
			continue
		}
		if username.Normalized() != tc.normalized {
			t.Errorf("%q: expected %q, got %q", tc.username, tc.normalized, username)
		}
	}
}
