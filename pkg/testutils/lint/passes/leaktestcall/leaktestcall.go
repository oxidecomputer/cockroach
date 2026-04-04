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

// Package leaktestcall defines an Analyzer that detects correct use
// of leaktest.AfterTest(t).
package leaktestcall

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is an analysis pass that checks that the return value of
// leaktest.AfterFunc(t) is called in defer statements.
var Analyzer = &analysis.Analyzer{
	Name:     "leaktestcall",
	Doc:      "Check that the closure returned by leaktest.AfterFunc(t) is called",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	astInspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	filter := []ast.Node{
		(*ast.DeferStmt)(nil),
	}

	astInspector.Preorder(filter, func(n ast.Node) {
		def := n.(*ast.DeferStmt)
		switch funCall := def.Call.Fun.(type) {
		case *ast.SelectorExpr:
			packageIdent, ok := funCall.X.(*ast.Ident)
			if !ok {
				return
			}
			if packageIdent.Name == "leaktest" && funCall.Sel != nil && funCall.Sel.Name == "AfterTest" {
				pass.Reportf(def.Call.Pos(), "leaktest.AfterTest return not called")
			}
		}
	})

	return nil, nil
}
