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

package kvserver

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/metric"
	"github.com/stretchr/testify/require"
)

// TestMetricRules tests the creation of metric rules related
// to KV metrics.
func TestMetricRules(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ruleRegistry := metric.NewRuleRegistry()
	CreateAndAddRules(context.Background(), ruleRegistry)
	require.NotNil(t, ruleRegistry.GetRuleForTest(unavailableRangesRuleName))
	require.NotNil(t, ruleRegistry.GetRuleForTest(trippedReplicaCircuitBreakersRuleName))
	require.NotNil(t, ruleRegistry.GetRuleForTest(underreplicatedRangesRuleName))
	require.NotNil(t, ruleRegistry.GetRuleForTest(requestsStuckInRaftRuleName))
	require.NotNil(t, ruleRegistry.GetRuleForTest(highOpenFDCountRuleName))
	require.NotNil(t, ruleRegistry.GetRuleForTest(nodeCapacityRuleName))
	require.NotNil(t, ruleRegistry.GetRuleForTest(clusterCapacityRuleName))
	require.NotNil(t, ruleRegistry.GetRuleForTest(nodeCapacityAvailableRuleName))
	require.NotNil(t, ruleRegistry.GetRuleForTest(clusterCapacityAvailableRuleName))
	require.NotNil(t, ruleRegistry.GetRuleForTest(capacityAvailableRatioRuleName))
	require.NotNil(t, ruleRegistry.GetRuleForTest(nodeCapacityAvailableRatioRuleName))
	require.NotNil(t, ruleRegistry.GetRuleForTest(clusterCapacityAvailableRatioRuleName))
	require.Equal(t, 12, ruleRegistry.GetRuleCountForTest())
}
