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

package clisqlshell

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/cli/clicfg"
	"github.com/cockroachdb/cockroach/pkg/cli/clisqlclient"
	"github.com/cockroachdb/cockroach/pkg/cli/clisqlexec"
	"github.com/cockroachdb/cockroach/pkg/sql/scanner"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/assert"
)

func TestIsEndOfStatement(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	tests := []struct {
		in         string
		isEnd      bool
		isNotEmpty bool
	}{
		{
			in:         ";",
			isEnd:      true,
			isNotEmpty: true,
		},
		{
			in:         "; /* comment */",
			isEnd:      true,
			isNotEmpty: true,
		},
		{
			in:         "; SELECT",
			isNotEmpty: true,
		},
		{
			in:         "SELECT",
			isNotEmpty: true,
		},
		{
			in:         "SET; SELECT 1;",
			isEnd:      true,
			isNotEmpty: true,
		},
		{
			in:         "SELECT ''''; SET;",
			isEnd:      true,
			isNotEmpty: true,
		},
		{
			in: "  -- hello",
		},
		{
			in:         "select 'abc", // invalid token
			isNotEmpty: true,
		},
		{
			in:         "'abc", // invalid token
			isNotEmpty: true,
		},
		{
			in:         `SELECT e'\xaa';`, // invalid token but last token is semicolon
			isEnd:      true,
			isNotEmpty: true,
		},
	}

	for _, test := range tests {
		lastTok, isNotEmpty := scanner.LastLexicalToken(test.in)
		if isNotEmpty != test.isNotEmpty {
			t.Errorf("%q: isNotEmpty expected %v, got %v", test.in, test.isNotEmpty, isNotEmpty)
		}
		isEnd := isEndOfStatement(lastTok)
		if isEnd != test.isEnd {
			t.Errorf("%q: isEnd expected %v, got %v", test.in, test.isEnd, isEnd)
		}
	}
}

// Test handleCliCmd cases for client-side commands that are aliases for sql
// statements.
func TestHandleCliCmdSqlAlias(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	clientSideCommandTestsTable := []struct {
		commandString string
		wantSQLStmt   string
	}{
		{`\l`, `SHOW DATABASES`},
		{`\dt`, `SHOW TABLES`},
		{`\dT`, `SHOW TYPES`},
		{`\du`, `SHOW USERS`},
		{`\du myuser`, `SELECT * FROM [SHOW USERS] WHERE username = 'myuser'`},
		{`\d mytable`, `SHOW COLUMNS FROM mytable`},
		{`\d`, `SHOW TABLES`},
	}

	for _, tt := range clientSideCommandTestsTable {
		c := setupTestCliState()
		c.lastInputLine = tt.commandString
		gotState := c.doHandleCliCmd(cliStateEnum(0), cliStateEnum(1))

		assert.Equal(t, cliRunStatement, gotState)
		assert.Equal(t, tt.wantSQLStmt, c.concatLines)
	}
}

func TestHandleCliCmdSlashDInvalidSyntax(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	clientSideCommandTests := []string{`\d goodarg badarg`, `\dz`}

	for _, tt := range clientSideCommandTests {
		c := setupTestCliState()
		c.lastInputLine = tt
		gotState := c.doHandleCliCmd(cliStateEnum(0), cliStateEnum(1))

		assert.Equal(t, cliStateEnum(0), gotState)
		assert.Equal(t, errInvalidSyntax, c.exitErr)
	}
}

func TestHandleDemoNodeCommandsInvalidNodeName(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	demoNodeCommandTests := []string{"shutdown", "*"}

	c := setupTestCliState()
	c.handleDemoNodeCommands(demoNodeCommandTests, cliStateEnum(0), cliStateEnum(1))
	assert.Equal(t, errInvalidSyntax, c.exitErr)
}

func setupTestCliState() *cliState {
	cliCtx := &clicfg.Context{}
	sqlConnCtx := &clisqlclient.Context{CliCtx: cliCtx}
	sqlExecCtx := &clisqlexec.Context{
		CliCtx:             cliCtx,
		TableDisplayFormat: clisqlexec.TableDisplayTable,
	}
	sqlCtx := &Context{}
	c := NewShell(cliCtx, sqlConnCtx, sqlExecCtx, sqlCtx, nil).(*cliState)
	c.ins = noLineEditor
	return c
}
