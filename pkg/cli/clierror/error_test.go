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

package clierror

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestOutputError(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	errBase := errors.New("woo")
	file, line, fn, _ := errors.GetOneLineSource(errBase)
	refLoc := fmt.Sprintf("%s, %s:%d", fn, file, line)
	testData := []struct {
		err                   error
		showSeverity, verbose bool
		exp                   string
	}{
		// Check the basic with/without severity.
		{errBase, false, false, "woo"},
		{errBase, true, false, "ERROR: woo"},
		{pgerror.WithCandidateCode(errBase, pgcode.Syntax), false, false, "woo\nSQLSTATE: " + pgcode.Syntax.String()},
		// Check the verbose output. This includes the uncategorized sqlstate.
		{errBase, false, true, "woo\nSQLSTATE: " + pgcode.Uncategorized.String() + "\nLOCATION: " + refLoc},
		{errBase, true, true, "ERROR: woo\nSQLSTATE: " + pgcode.Uncategorized.String() + "\nLOCATION: " + refLoc},
		// Check the same over pq.Error objects.
		{&pq.Error{Message: "woo"}, false, false, "woo"},
		{&pq.Error{Message: "woo"}, true, false, "ERROR: woo"},
		{&pq.Error{Message: "woo"}, false, true, "woo"},
		{&pq.Error{Message: "woo"}, true, true, "ERROR: woo"},
		{&pq.Error{Severity: "W", Message: "woo"}, false, false, "woo"},
		{&pq.Error{Severity: "W", Message: "woo"}, true, false, "W: woo"},
		// Check hint printed after message.
		{errors.WithHint(errBase, "hello"), false, false, "woo\nHINT: hello"},
		// Check sqlstate printed before hint, location after hint.
		{errors.WithHint(errBase, "hello"), false, true, "woo\nSQLSTATE: " + pgcode.Uncategorized.String() + "\nHINT: hello\nLOCATION: " + refLoc},
		// Check detail printed after message.
		{errors.WithDetail(errBase, "hello"), false, false, "woo\nDETAIL: hello"},
		// Check hint/detail collection, hint printed after detail.
		{errors.WithHint(
			errors.WithDetail(
				errors.WithHint(errBase, "a"),
				"b"),
			"c"), false, false, "woo\nDETAIL: b\nHINT: a\n--\nc"},
		{errors.WithDetail(
			errors.WithHint(
				errors.WithDetail(errBase, "a"),
				"b"),
			"c"), false, false, "woo\nDETAIL: a\n--\nc\nHINT: b"},
		// Check sqlstate printed before detail, location after hint.
		{errors.WithDetail(
			errors.WithHint(errBase, "a"), "b"),
			false, true, "woo\nSQLSTATE: " + pgcode.Uncategorized.String() + "\nDETAIL: b\nHINT: a\nLOCATION: " + refLoc},
	}

	for _, tc := range testData {
		var buf strings.Builder
		OutputError(&buf, tc.err, tc.showSeverity, tc.verbose)
		assert.Equal(t, tc.exp+"\n", buf.String())
	}
}

func TestFormatLocation(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	testData := []struct {
		file, line, fn string
		exp            string
	}{
		{"", "", "", ""},
		{"a.b", "", "", "a.b"},
		{"", "123", "", "<unknown>:123"},
		{"", "", "abc", "abc"},
		{"a.b", "", "abc", "abc, a.b"},
		{"a.b", "123", "", "a.b:123"},
		{"", "123", "abc", "abc, <unknown>:123"},
	}

	for _, tc := range testData {
		r := formatLocation(tc.file, tc.line, tc.fn)
		assert.Equal(t, tc.exp, r)
	}
}
