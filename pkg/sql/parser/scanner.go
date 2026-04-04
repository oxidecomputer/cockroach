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

package parser

import (
	"github.com/cockroachdb/cockroach/pkg/sql/lexbase"
	"github.com/cockroachdb/cockroach/pkg/sql/scanner"
)

func makeScanner(str string) scanner.Scanner {
	var s scanner.Scanner
	s.Init(str)
	return s
}

// SplitFirstStatement returns the length of the prefix of the string up to and
// including the first semicolon that separates statements. If there is no
// including the first semicolon that separates statements. If there is no
// semicolon, returns ok=false.
func SplitFirstStatement(sql string) (pos int, ok bool) {
	s := makeScanner(sql)
	var lval = &sqlSymType{}
	for {
		s.Scan(lval)
		switch lval.ID() {
		case 0, lexbase.ERROR:
			return 0, false
		case ';':
			return s.Pos(), true
		}
	}
}

// Tokens decomposes the input into lexical tokens.
func Tokens(sql string) (tokens []TokenString, ok bool) {
	s := makeScanner(sql)
	for {
		var lval = &sqlSymType{}
		s.Scan(lval)
		if lval.ID() == lexbase.ERROR {
			return nil, false
		}
		if lval.ID() == 0 {
			break
		}
		tokens = append(tokens, TokenString{TokenID: lval.ID(), Str: lval.Str()})
	}
	return tokens, true
}

// TokensIgnoreErrors decomposes the input into lexical tokens and
// ignores errors.
func TokensIgnoreErrors(sql string) (tokens []TokenString) {
	s := makeScanner(sql)
	for {
		var lval = &sqlSymType{}
		s.Scan(lval)
		if lval.ID() == 0 {
			break
		}
		tokens = append(tokens, TokenString{TokenID: lval.ID(), Str: lval.Str()})
	}
	return tokens
}

// TokenString is the unit value returned by Tokens.
type TokenString struct {
	TokenID int32
	Str     string
}
