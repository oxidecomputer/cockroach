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

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func valueForYAML(v interface{}) interface{} {
	switch v := v.(type) {
	case fmt.Stringer:
		return v.String()
	default:
		return v
	}
}

func (v valueExpr) encoded() interface{} {
	return valueForYAML(v.value)
}

func (a anyExpr) encoded() interface{} {
	ret := make([]interface{}, 0, len(a))
	for _, v := range a {
		ret = append(ret, valueForYAML(v))
	}
	return ret
}

func (v Var) encoded() interface{} {
	return "$" + string(v)
}

func (e *eqDecl) MarshalYAML() (interface{}, error) {
	return clauseStr("$"+string(e.v), e.expr)
}

func exprToString(e expr) (string, error) {
	var expr yaml.Node
	if err := expr.Encode(e.encoded()); err != nil {
		return "", err
	}
	expr.Style = yaml.FlowStyle
	out, err := yaml.Marshal(&expr)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), err
}

func (f *tripleDecl) MarshalYAML() (interface{}, error) {
	return clauseStr(fmt.Sprintf("$%s[%s]", f.entity, f.attribute), f.value)
}

func clauseStr(lhs string, rhs expr) (string, error) {
	rhsStr, err := exprToString(rhs)
	if err != nil {
		return "", err
	}
	op := "="
	if _, isAny := rhs.(anyExpr); isAny {
		op = "IN"
	}
	return fmt.Sprintf("%s %s %s", lhs, op, rhsStr), nil
}

func (f filterDecl) MarshalYAML() (interface{}, error) {
	var buf strings.Builder
	buf.WriteString(f.name)
	buf.WriteString("(")
	ft := reflect.TypeOf(f.predicateFunc)
	for i := 0; i < ft.NumIn(); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(ft.In(i).String())
	}
	buf.WriteString(")(")
	for i, v := range f.vars {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("$")
		buf.WriteString(string(v))
	}
	buf.WriteString(")")
	return buf.String(), nil
}
