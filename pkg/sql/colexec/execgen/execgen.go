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
	"go/parser"
	"go/token"
	"strings"

	"github.com/dave/dst/decorator"
)

// Generate transforms the string contents of an input execgen template by
// processing all supported // execgen annotations.
func Generate(inputFileContents string) (string, error) {
	f, err := decorator.ParseFile(token.NewFileSet(), "", inputFileContents, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Generate template variants: // execgen:template
	expandTemplates(f)

	// Inline functions: // execgen:inline
	inlineFuncs(f)

	// Produce output string.
	var sb strings.Builder
	if err := decorator.Fprint(&sb, f); err != nil {
		panic(err)
	}
	return sb.String(), nil
}
