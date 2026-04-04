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
	"fmt"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

func prettyPrintStmts(stmts ...dst.Stmt) string {
	if len(stmts) == 0 {
		return ""
	}
	f := &dst.File{
		Name: dst.NewIdent("main"),
		Decls: []dst.Decl{
			&dst.FuncDecl{
				Name: dst.NewIdent("test"),
				Type: &dst.FuncType{},
				Body: &dst.BlockStmt{
					List: stmts,
				},
			},
		},
	}
	var ret strings.Builder
	_ = decorator.Fprint(&ret, f)
	prelude := `package main

func test() {
`
	postlude := `}
`
	s := ret.String()
	return strings.TrimSpace(s[len(prelude) : len(s)-len(postlude)])
}

func prettyPrintExprs(exprs ...dst.Expr) string {
	stmts := make([]dst.Stmt, len(exprs))
	for i := range exprs {
		stmts[i] = &dst.ExprStmt{X: exprs[i]}
	}
	return prettyPrintStmts(stmts...)
}

func parseStmt(stmt string) (dst.Stmt, error) {
	f, err := decorator.Parse(fmt.Sprintf(
		`package main
func test() {
	%s
}`, stmt))
	if err != nil {
		return nil, err
	}
	return f.Decls[0].(*dst.FuncDecl).Body.List[0], nil
}

func mustParseStmt(stmt string) dst.Stmt {
	ret, err := parseStmt(stmt)
	if err != nil {
		panic(err)
	}
	return ret
}
