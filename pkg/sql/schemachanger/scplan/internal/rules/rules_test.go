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

package rules

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/datadriven"
	"gopkg.in/yaml.v3"
)

// TestRulesYAML outputs the rules to yaml as a way to visualize changes.
func TestRulesYAML(t *testing.T) {
	datadriven.Walk(t, testutils.TestDataPath(t), func(t *testing.T, path string) {
		datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
			switch d.Cmd {
			case "deprules":
				out, err := yaml.Marshal(registry.depRules)
				if err != nil {
					d.Fatalf(t, "failed to marshal deprules: %v", err)
				}
				return string(out)
			case "oprules":
				out, err := yaml.Marshal(registry.opRules)
				if err != nil {
					d.Fatalf(t, "failed to marshal oprules: %v", err)
				}
				return string(out)
			}
			d.Fatalf(t, "deprules and oprules are the only commands")
			return ""
		})
	})
}

func (r registeredDepRule) MarshalYAML() (interface{}, error) {
	var query yaml.Node
	if err := query.Encode(r.q.Clauses()); err != nil {
		return nil, err
	}
	return &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "name"},
			{Kind: yaml.ScalarNode, Value: r.name},
			{Kind: yaml.ScalarNode, Value: "from"},
			{Kind: yaml.ScalarNode, Value: string(r.from)},
			{Kind: yaml.ScalarNode, Value: "kind"},
			{Kind: yaml.ScalarNode, Value: r.kind.String()},
			{Kind: yaml.ScalarNode, Value: "to"},
			{Kind: yaml.ScalarNode, Value: string(r.to)},
			{Kind: yaml.ScalarNode, Value: "query"},
			&query,
		},
	}, nil
}

func (r registeredOpRule) MarshalYAML() (interface{}, error) {
	var query yaml.Node
	if err := query.Encode(r.q.Clauses()); err != nil {
		return nil, err
	}
	return &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "name"},
			{Kind: yaml.ScalarNode, Value: r.name},
			{Kind: yaml.ScalarNode, Value: "from"},
			{Kind: yaml.ScalarNode, Value: string(r.from)},
			{Kind: yaml.ScalarNode, Value: "query"},
			&query,
		},
	}, nil
}
