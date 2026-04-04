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

package treebin

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

// BinaryOperator represents a unary operator used in a BinaryExpr.
type BinaryOperator struct {
	Symbol BinaryOperatorSymbol
	// IsExplicitOperator is true if OPERATOR(symbol) is used.
	IsExplicitOperator bool
}

// MakeBinaryOperator creates a BinaryOperator given a symbol.
func MakeBinaryOperator(symbol BinaryOperatorSymbol) BinaryOperator {
	return BinaryOperator{Symbol: symbol}
}

func (o BinaryOperator) String() string {
	if o.IsExplicitOperator {
		return fmt.Sprintf("OPERATOR(%s)", o.Symbol.String())
	}
	return o.Symbol.String()
}

// Operator implements tree.Operator.
func (BinaryOperator) Operator() {}

// BinaryOperatorSymbol is a symbol for a binary operator.
type BinaryOperatorSymbol uint8

// BinaryExpr.Operator
const (
	Bitand BinaryOperatorSymbol = iota
	Bitor
	Bitxor
	Plus
	Minus
	Mult
	Div
	FloorDiv
	Mod
	Pow
	Concat
	LShift
	RShift
	JSONFetchVal
	JSONFetchText
	JSONFetchValPath
	JSONFetchTextPath

	NumBinaryOperatorSymbols
)

var _ = NumBinaryOperatorSymbols

var binaryOpName = [...]string{
	Bitand:            "&",
	Bitor:             "|",
	Bitxor:            "#",
	Plus:              "+",
	Minus:             "-",
	Mult:              "*",
	Div:               "/",
	FloorDiv:          "//",
	Mod:               "%",
	Pow:               "^",
	Concat:            "||",
	LShift:            "<<",
	RShift:            ">>",
	JSONFetchVal:      "->",
	JSONFetchText:     "->>",
	JSONFetchValPath:  "#>",
	JSONFetchTextPath: "#>>",
}

// IsPadded returns whether the binary operator needs to be padded.
func (i BinaryOperatorSymbol) IsPadded() bool {
	return !(i == JSONFetchVal || i == JSONFetchText || i == JSONFetchValPath || i == JSONFetchTextPath)
}

func (i BinaryOperatorSymbol) String() string {
	if i > BinaryOperatorSymbol(len(binaryOpName)-1) {
		return fmt.Sprintf("BinaryOp(%d)", i)
	}
	return binaryOpName[i]
}

// BinaryOpName returns the name of op.
func BinaryOpName(op BinaryOperatorSymbol) string {
	if int(op) >= len(binaryOpName) || binaryOpName[op] == "" {
		panic(errors.AssertionFailedf("missing name for operator %q", op.String()))
	}
	return binaryOpName[op]
}
