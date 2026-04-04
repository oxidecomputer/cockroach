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

package cyclegraphtest

import (
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/rel"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/rel/reltest"
)

// Suite is the suite for testproto.
var Suite = reltest.Suite{
	Name:          "cyclegraph",
	Schema:        schema,
	Registry:      reg,
	DatabaseTests: databaseTests,
}

var (
	reg = reltest.NewRegistry()

	// Here we create a bunch of cyclic relationships between these 4 entities.
	message1, message2, container1, container2 = func() (*struct1, *struct2, *container, *container) {
		s1, s2, c1, c2 := &struct1{}, &struct2{}, &container{}, &container{}
		*c1 = container{S1: s1}
		*c2 = container{S2: s2}
		*s1 = struct1{Name: "message1", S1: s1, S2: s2, C: c2}
		*s2 = struct2{Name: "message2", S1: s1, S2: s2, C: c1}

		reg.Register("message1", s1)
		reg.Register("message2", s2)
		reg.Register("container1", c1)
		reg.Register("container2", c2)

		return s1, s2, c1, c2
	}()

	databaseTests = []reltest.DatabaseTest{
		{
			Data: []string{"container1"}, // recursively will add it all, test that
			Indexes: [][][]rel.Attr{
				nil,
				{{s}, {c}, {name}},
			},
			QueryCases: []reltest.QueryTest{
				{
					Name: "oneOf member",
					Query: rel.Clauses{
						rel.Var("c").Type((*container)(nil)),
						rel.Var("c").AttrEqVar(s, "s"),
					},
					ResVars:  []rel.Var{"c", "s"},
					Entities: []rel.Var{"c"},
					Results: [][]interface{}{
						{container1, message1},
						{container2, message2},
					},
				},
				{
					Name: "oneOf member",
					Query: rel.Clauses{
						rel.Var("c").AttrEq(s, message1),
					},
					ResVars:  []rel.Var{"c"},
					Entities: []rel.Var{"c"},
					Results: [][]interface{}{
						{container1},
					},
				},
			},
		},
	}
)
