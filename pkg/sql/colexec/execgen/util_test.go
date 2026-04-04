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

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

func parseStmts(stmts string) []dst.Stmt {
	inputStr := fmt.Sprintf(`package main
func test() {
  %s
}`, stmts)
	f, err := decorator.Parse(inputStr)
	if err != nil {
		panic(err)
	}
	return f.Decls[0].(*dst.FuncDecl).Body.List
}

func parseDecls(decls string) []dst.Decl {
	inputStr := fmt.Sprintf(`package main
%s
`, decls)
	f, err := decorator.Parse(inputStr)
	if err != nil {
		panic(err)
	}
	return f.Decls
}
