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

package rel

// expr is an expression value in rel.
type expr interface {
	expr() // marker

	// encoded returns a value for use in serialization.
	encoded() interface{}
}

// Var is a variable name. Everything is convention, but, when you create
// clauses and use variable names which are not part of the defined scope of
// the rule, the new variableSlots which will be created will have a scope
// prefix to try to ensure that they are unique. Given that, don't put `:` in
// your variable names.
type Var string

// Var is an expr.
func (Var) expr() {}

type valueExpr struct {
	value interface{}
}

func (v valueExpr) expr() {}

type anyExpr []interface{}

func (a anyExpr) expr() {}
