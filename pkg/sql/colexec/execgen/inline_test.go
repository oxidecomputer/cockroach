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

package execgen

import (
	"testing"

	"github.com/dave/dst"
	"github.com/stretchr/testify/assert"
)

func TestGetFormalParamReassignments(t *testing.T) {
	tests := []struct {
		funcDecl string
		callExpr string

		expectedStmts string
	}{
		{
			funcDecl:      `func a() {}`,
			callExpr:      `a()`,
			expectedStmts: ``,
		},
		{
			funcDecl:      `func a(a int) {}`,
			callExpr:      `a(b)`,
			expectedStmts: `var a int = b`,
		},
		{
			funcDecl: `func a(a int, b int) {}`,
			callExpr: `a(x, y)`,
			expectedStmts: `var (
		a int = x
		b int = y
	)`,
		},
		{
			funcDecl:      `func a(a int, b int) {}`,
			callExpr:      `a(a, c)`,
			expectedStmts: `var b int = c`,
		},
	}
	for _, tt := range tests {
		callExpr := parseStmts(tt.callExpr)[0].(*dst.ExprStmt).X.(*dst.CallExpr)
		funcDecl := parseDecls(tt.funcDecl)[0].(*dst.FuncDecl)
		stmt := getFormalParamReassignments(funcDecl, callExpr)
		actual := prettyPrintStmts(stmt)
		assert.Equal(t, tt.expectedStmts, actual)
	}
}

func TestExtractReturnValues(t *testing.T) {
	tests := []struct {
		decl             string
		expectedRetDecls string
	}{
		{
			decl:             "func foo(a int) {}",
			expectedRetDecls: "",
		},
		{
			decl: "func foo(a int) (int, string) {}",
			expectedRetDecls: `var (
		__retval_0 int
		__retval_1 string
	)`,
		},
		{
			decl:             "func foo(a int) int {}",
			expectedRetDecls: `var __retval_0 int`,
		},
		{
			decl: "func foo(a int) (a int, b string) {}",
			expectedRetDecls: `var (
		__retval_a int
		__retval_b string
	)`,
		},
	}
	for _, tt := range tests {
		decl := parseDecls(tt.decl)[0].(*dst.FuncDecl)
		retValDecl, retValNames := extractReturnValues(decl)
		if _, ok := retValDecl.(*dst.EmptyStmt); ok {
			assert.Equal(t, 0, len(retValNames))
		} else {
			specs := retValDecl.(*dst.DeclStmt).Decl.(*dst.GenDecl).Specs
			assert.Equal(t, len(specs), len(retValNames))
			for i := range retValNames {
				assert.Equal(t, retValNames[i], specs[i].(*dst.ValueSpec).Names[0].Name)
			}
		}
		assert.Equal(t, tt.expectedRetDecls, prettyPrintStmts(retValDecl))
	}
}
