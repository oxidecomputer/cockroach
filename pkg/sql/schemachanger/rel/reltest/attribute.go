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

package reltest

import (
	"sort"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/rel"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// AttributeTestCases exercise the mappings of fields to attribute values.
type AttributeTestCases []AttributeTestCase

// AttributeTestCase is a case specifying how an entity in the registry
// should be mapped
type AttributeTestCase struct {
	Entity   string
	Expected map[rel.Attr]interface{}
}

func (c AttributeTestCase) run(t *testing.T, s Suite) {
	got := make(map[rel.Attr]interface{})
	e := s.Registry.MustGetByName(t, c.Entity)
	require.NoError(t, s.Schema.IterateAttributes(e, func(
		attribute rel.Attr, value interface{},
	) error {
		got[attribute] = value
		return nil
	}))
	require.Equal(t, c.Expected, got)
	for a, v := range got {
		gotVal, err := s.Schema.GetAttribute(a, e)
		require.NoError(t, err)
		require.Equalf(t, v, gotVal, "%T[%v]", v, a.String())
	}
}
func (c AttributeTestCase) encode(t *testing.T, s Suite) *yaml.Node {
	type kv struct {
		k string
		v *yaml.Node
	}
	var kvs []kv
	for attr, value := range c.Expected {
		var v *yaml.Node
		if name, ok := s.Registry.GetName(value); ok {
			v = scalarYAML(name)
		} else {
			v = s.Registry.EncodeToYAML(t, value)
		}
		kvs = append(kvs, kv{attr.String(), v})
	}
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].k < kvs[j].k })
	var n yaml.Node
	n.Kind = yaml.MappingNode
	n.Style = yaml.FlowStyle
	for _, kv := range kvs {
		n.Content = append(n.Content,
			scalarYAML(kv.k),
			kv.v,
		)
	}
	return &n
}
