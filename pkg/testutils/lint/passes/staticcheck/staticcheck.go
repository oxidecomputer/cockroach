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

// Package staticcheck provides utilities for consuming `staticcheck` checks in
// `nogo`.
package staticcheck

import (
	"strings"

	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/analysis/facts/directives"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/analysis/report"
)

// MungeAnalyzer updates an Analyzer from `staticcheck` to make it work w/ nogo.
// The staticcheck analyzers don't look at "lint:ignore" directives, so if you
// integrate them into `nogo` unchanged, you'll get spurious build failures for
// issues that are actually explicitly ignored. So for each staticcheck analyzer
// we add `directives.Analyzer` to the list of dependencies, then cross-check
// each reported diagnostic to make sure it's not ignored before allowing it
// through.
func MungeAnalyzer(analyzer *analysis.Analyzer) {
	// Add facts.directives to the list of dependencies for this analyzer.
	analyzer.Requires = analyzer.Requires[0:len(analyzer.Requires):len(analyzer.Requires)]
	analyzer.Requires = append(analyzer.Requires, directives.Analyzer)
	oldRun := analyzer.Run
	analyzer.Run = func(p *analysis.Pass) (interface{}, error) {
		pass := *p
		oldReport := p.Report
		pass.Report = func(diag analysis.Diagnostic) {
			dirs := pass.ResultOf[directives.Analyzer].([]lint.Directive)
			for _, dir := range dirs {
				cmd := dir.Command
				args := dir.Arguments
				switch cmd {
				case "ignore":
					ignorePos := report.DisplayPosition(pass.Fset, dir.Node.Pos())
					nodePos := report.DisplayPosition(pass.Fset, diag.Pos)
					if ignorePos.Filename != nodePos.Filename || ignorePos.Line != nodePos.Line {
						continue
					}
					for _, check := range strings.Split(args[0], ",") {
						if check == analyzer.Name {
							// Skip reporting the diagnostic.
							return
						}
					}
				case "file-ignore":
					ignorePos := report.DisplayPosition(pass.Fset, dir.Node.Pos())
					nodePos := report.DisplayPosition(pass.Fset, diag.Pos)
					if ignorePos.Filename == nodePos.Filename {
						// Skip reporting the diagnostic.
						return
					}
				default:
					// Unknown directive, ignore
					continue
				}
			}
			oldReport(diag)
		}
		return oldRun(&pass)
	}
}
