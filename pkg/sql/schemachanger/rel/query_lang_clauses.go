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

import "gopkg.in/yaml.v3"

// Clauses exists to handle flattening of a slice of clauses before marshaling.
type Clauses []Clause

func flattened(c Clauses) Clauses {
	hasAnd := func() bool {
		for _, cl := range c {
			if _, isAnd := cl.(and); isAnd {
				return true
			}
		}
		return false
	}
	if !hasAnd() {
		return c
	}
	var ret Clauses
	for _, cl := range c {
		switch cl := cl.(type) {
		case and:
			ret = append(ret, flattened(Clauses(cl))...)
		default:
			ret = append(ret, cl)
		}
	}
	return ret
}

// MarshalYAML marshals clauses to yaml.
func (c Clauses) MarshalYAML() (interface{}, error) {
	fc := flattened(c)
	var n yaml.Node
	if err := n.Encode([]Clause(fc)); err != nil {
		return nil, err
	}
	n.Style = yaml.LiteralStyle
	return &n, nil
}
