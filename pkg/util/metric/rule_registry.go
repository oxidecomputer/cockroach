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

package metric

import "github.com/cockroachdb/cockroach/pkg/util/syncutil"

// RuleRegistry is a list of all rules (AlertingRule and AggregationRule).
//
// All defined rules should be registered in the RuleRegistry to be exported
// as Prometheus alert/recording rules.
type RuleRegistry struct {
	syncutil.Mutex
	rules []Rule
}

// NewRuleRegistry creates a new RuleRegistry.
func NewRuleRegistry() *RuleRegistry {
	return &RuleRegistry{
		rules: []Rule{},
	}
}

// AddRule adds a rule to the registry.
func (r *RuleRegistry) AddRule(rule Rule) {
	r.Lock()
	defer r.Unlock()
	r.rules = append(r.rules, rule)
}

// AddRules adds multiple rules to the registry.
func (r *RuleRegistry) AddRules(rules []Rule) {
	r.Lock()
	defer r.Unlock()
	r.rules = append(r.rules, rules...)
}

// Each calls the given closure for all rules.
func (r *RuleRegistry) Each(f func(rule Rule)) {
	r.Lock()
	defer r.Unlock()
	for _, currentRule := range r.rules {
		f(currentRule)
	}
}
