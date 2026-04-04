// Copyright 2022 The Cockroach Authors.
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

package delegate

import (
	"reflect"
	"testing"
)

func TestCompletions(t *testing.T) {
	tests := []struct {
		stmt                string
		offset              int
		expectedCompletions []string
	}{
		{
			stmt:                "creat",
			expectedCompletions: []string{"CREATE", "CREATEDB", "CREATELOGIN", "CREATEROLE"},
		},
		{
			stmt:                "CREAT",
			expectedCompletions: []string{"CREATE", "CREATEDB", "CREATELOGIN", "CREATEROLE"},
		},
		{
			stmt:                "creat ",
			expectedCompletions: []string{},
		},
		{
			stmt:                "SHOW CREAT",
			expectedCompletions: []string{"CREATE", "CREATEDB", "CREATELOGIN", "CREATEROLE"},
		},
		{
			stmt:                "show creat",
			expectedCompletions: []string{"CREATE", "CREATEDB", "CREATELOGIN", "CREATEROLE"},
		},
		{
			stmt: "se",
			expectedCompletions: []string{
				"SEARCH", "SECOND", "SELECT", "SEQUENCE", "SEQUENCES",
				"SERIALIZABLE", "SERVER", "SESSION", "SESSIONS", "SESSION_USER",
				"SET", "SETS", "SETTING", "SETTINGS",
			},
		},
		{
			stmt:                "sel",
			expectedCompletions: []string{"SELECT"},
		},
		{
			stmt:                "create ta",
			expectedCompletions: []string{"TABLE", "TABLES", "TABLESPACE"},
		},
		{
			stmt:                "create ta",
			expectedCompletions: []string{"CREATE"},
			offset:              3,
		},
		{
			stmt:                "select",
			expectedCompletions: []string{"SELECT"},
			offset:              2,
		},
		{
			stmt:                "select ",
			expectedCompletions: []string{},
			offset:              7,
		},
		{
			stmt:                "你好，我的名字是鲍勃 SELECT",
			expectedCompletions: []string{"你好，我的名字是鲍勃"},
			offset:              2,
		},
		{
			stmt:                "你好，我的名字是鲍勃 SELECT",
			expectedCompletions: []string{},
			offset:              11,
		},
		{
			stmt:                "你好，我的名字是鲍勃 SELECT",
			expectedCompletions: []string{"SELECT"},
			offset:              12,
		},
		{
			stmt:                "😋😋😋 😋😋😋",
			expectedCompletions: []string{},
			offset:              25,
		},
		{
			stmt:                "Jalapeño",
			expectedCompletions: []string{},
			offset:              9,
		},
	}
	for _, tc := range tests {
		offset := tc.offset
		if tc.offset == 0 {
			offset = len(tc.stmt)
		}
		completions, err := RunShowCompletions(tc.stmt, offset)
		if err != nil {
			t.Error(err)
		}
		if !(len(completions) == 0 && len(tc.expectedCompletions) == 0) &&
			!reflect.DeepEqual(completions, tc.expectedCompletions) {
			t.Errorf("expected %v, got %v", tc.expectedCompletions, completions)
		}
	}
}
