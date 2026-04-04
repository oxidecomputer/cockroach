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

// GetRuleForTest retrieves individual rule for testing purposes from
// the rule registry. The rules are retrieved in O(#rules) time.
func (r *RuleRegistry) GetRuleForTest(ruleName string) Rule {
	r.Lock()
	defer r.Unlock()
	for _, rule := range r.rules {
		if rule.Name() == ruleName {
			return rule
		}
	}
	return nil
}

// GetRuleCountForTest returns the total number of rules.
func (r *RuleRegistry) GetRuleCountForTest() int {
	return len(r.rules)
}
