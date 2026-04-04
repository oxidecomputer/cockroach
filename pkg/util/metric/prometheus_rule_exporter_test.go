// Copyright 2015 The Cockroach Authors.
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

package metric

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestPrometheusRuleExporter tests that the PrometheusRuleExporter
// generates valid YAML output for a given set of alerting and
// aggregation rules.
func TestPrometheusRuleExporter(t *testing.T) {
	ctx := context.Background()
	rules, expectedYAMLText := getRulesAndExpectedYAML(t)
	registry := NewRuleRegistry()
	registry.AddRules(rules)
	ruleExporter := NewPrometheusRuleExporter(registry)
	ruleExporter.ScrapeRegistry(ctx)
	yaml, err := ruleExporter.PrintAsYAML()
	require.NoError(t, err)
	require.Equal(t, expectedYAMLText, string(yaml))
}

func getRulesAndExpectedYAML(t *testing.T) (rules []Rule, expectedYAML string) {
	aggRule, err := NewAggregationRule("test_aggregation_rule_1", "sum without(store) (capacity_available)", nil, "", false)
	if err != nil {
		t.Fatal(err)
	}
	rules = append(rules, aggRule)
	var annotations []LabelPair
	description := "description"
	value := "{{ $labels.instance }} for cluster {{ $labels.cluster }} has\n been down for more than 15 minutes."
	annotation := LabelPair{
		Name:  &description,
		Value: &value,
	}
	annotations = append(annotations, annotation)
	alertRule, err := NewAlertingRule("test_alert_rule", "resets(sys_uptime[10m]) > 0 and resets(sys_uptime[10m]) < 5", annotations, nil, time.Minute, "", false)
	if err != nil {
		t.Fatal(err)
	}
	rules = append(rules, alertRule)
	expectedYAML = "rules/alerts:\n    rules:\n        - alert: test_alert_rule\n          " +
		"expr: resets(sys_uptime[10m]) > 0 and resets(sys_uptime[10m]) < 5\n          for: 1m0s\n          " +
		"annotations:\n            description: |-\n                {{ $labels.instance }} for " +
		"cluster {{ $labels.cluster }} has\n                 been down for more than 15 minutes.\nrules/recording:" +
		"\n    rules:\n        - record: test_aggregation_rule_1\n          expr: sum without(store) (capacity_available)\n"

	return rules, expectedYAML
}
