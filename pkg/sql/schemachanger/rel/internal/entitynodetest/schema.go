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

package entitynodetest

import (
	"reflect"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/rel"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/rel/reltest"
	"gopkg.in/yaml.v3"
)

type entity struct {
	I8  int8
	PI8 *int8
	I16 int16
}

type node struct {
	Value       *entity
	Left, Right *node
}

func (n *node) EncodeToYAML(t *testing.T, r *reltest.Registry) interface{} {
	yn := yaml.Node{Kind: yaml.MappingNode, Style: yaml.FlowStyle}
	for _, f := range []struct {
		name  string
		field interface{}
		ok    bool
	}{
		{"value", n.Value, n.Value != nil},
		{"left", n.Left, n.Left != nil},
		{"right", n.Right, n.Right != nil},
	} {
		if !f.ok {
			continue
		}
		yn.Content = append(yn.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: f.name},
			&yaml.Node{Kind: yaml.ScalarNode, Value: r.MustGetName(t, f.field)},
		)
	}
	return &yn
}

var _ reltest.RegistryYAMLEncoder = (*node)(nil)

// testAttr is a rel.Attr used for testing.
type testAttr int8

var _ rel.Attr = testAttr(0)

//go:generate stringer --type testAttr  --tags test
const (
	i8 testAttr = iota
	pi8
	i16
	value
	left
	right
)

var schema = rel.MustSchema("testschema",
	rel.EntityMapping(reflect.TypeOf((*entity)(nil)),
		rel.EntityAttr(i8, "I8"),
		rel.EntityAttr(pi8, "PI8"),
		rel.EntityAttr(i16, "I16"),
	),
	rel.EntityMapping(reflect.TypeOf((*node)(nil)),
		rel.EntityAttr(value, "Value"),
		rel.EntityAttr(left, "Left"),
		rel.EntityAttr(right, "Right"),
	),
)
