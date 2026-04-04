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

package tree

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestRoleSpecValidation(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	testCases := []struct {
		username   string
		normalized string
		err        string
		sqlstate   pgcode.Code
	}{
		{"Abc123", "abc123", "", pgcode.Code{}},
		{"0123121132", "0123121132", "", pgcode.Code{}},
		{"HeLlO", "hello", "", pgcode.Code{}},
		{"Ομηρος", "ομηρος", "", pgcode.Code{}},
		{"_HeLlO", "_hello", "", pgcode.Code{}},
		{"a-BC-d", "a-bc-d", "", pgcode.Code{}},
		{"A.Bcd", "a.bcd", "", pgcode.Code{}},
		{"WWW.BIGSITE.NET", "www.bigsite.net", "", pgcode.Code{}},
		{"", "", `"": username is empty`, pgcode.InvalidName},
		{"-ABC", "-abc", `"-abc": username is invalid`, pgcode.InvalidName},
		{".ABC", ".abc", `".abc": username is invalid`, pgcode.InvalidName},
		{"*.wildcard", "*.wildcard", `"\*.wildcard": username is invalid`, pgcode.InvalidName},
		{"foofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoof",
			"foofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoof",
			`"foofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoofoof": username is too long`, pgcode.NameTooLong},
		{"M", "m", "", pgcode.Code{}},
		{".", ".", `".": username is invalid`, pgcode.InvalidName},
	}

	for _, tc := range testCases {
		roleSpec := RoleSpec{RoleSpecType: RoleName, Name: tc.username}
		normalized, err := roleSpec.ToSQLUsername(&sessiondata.SessionData{}, security.UsernameCreation)
		if !testutils.IsError(err, tc.err) {
			t.Errorf("%q: expected %q, got %v", tc.username, tc.err, err)
			continue
		}
		if err != nil {
			if pgcode := pgerror.GetPGCode(err); pgcode != tc.sqlstate {
				t.Errorf("%q: expected SQLSTATE %s, got %s", tc.username, tc.sqlstate, pgcode)
				continue
			}
		}
		if normalized.Normalized() != tc.normalized {
			t.Errorf("%q: expected %q, got %q", tc.username, tc.normalized, normalized.Normalized())
		}
	}
}
