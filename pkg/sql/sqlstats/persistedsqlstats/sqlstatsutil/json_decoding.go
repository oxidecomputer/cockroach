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

package sqlstatsutil

import (
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/json"
	"github.com/cockroachdb/errors"
)

// DecodeTxnStatsMetadataJSON decodes the 'metadata' field of the JSON
// representation of transaction statistics into
// roachpb.CollectedTransactionStatistics.
func DecodeTxnStatsMetadataJSON(
	metadata json.JSON, result *roachpb.CollectedTransactionStatistics,
) error {
	return jsonFields{
		{"stmtFingerprintIDs", (*stmtFingerprintIDArray)(&result.StatementFingerprintIDs)},
	}.decodeJSON(metadata)
}

// DecodeTxnStatsStatisticsJSON decodes the 'statistics' section of the
// transaction statistics JSON payload into roachpb.TransactionStatistics
// protobuf.
func DecodeTxnStatsStatisticsJSON(jsonVal json.JSON, result *roachpb.TransactionStatistics) error {
	return (*txnStats)(result).decodeJSON(jsonVal)
}

// DecodeStmtStatsMetadataJSON decodes the 'metadata' field of the JSON
// representation of the statement statistics into
// roachpb.CollectedStatementStatistics.
func DecodeStmtStatsMetadataJSON(
	metadata json.JSON, result *roachpb.CollectedStatementStatistics,
) error {
	return (*stmtStatsMetadata)(result).jsonFields().decodeJSON(metadata)
}

// DecodeAggregatedMetadataJSON decodes the 'aggregated metadata' represented by roachpb.AggregatedStatementMetadata.
func DecodeAggregatedMetadataJSON(
	metadata json.JSON, result *roachpb.AggregatedStatementMetadata,
) error {
	return (*aggregatedMetadata)(result).jsonFields().decodeJSON(metadata)
}

// DecodeStmtStatsStatisticsJSON decodes the 'statistics' field and the
// 'execution_statistics' field in the given json into
// roachpb.StatementStatistics.
func DecodeStmtStatsStatisticsJSON(jsonVal json.JSON, result *roachpb.StatementStatistics) error {
	return (*stmtStats)(result).decodeJSON(jsonVal)
}

// JSONToExplainTreePlanNode decodes the JSON-formatted ExplainTreePlanNode
// produced by ExplainTreePlanNodeToJSON.
func JSONToExplainTreePlanNode(jsonVal json.JSON) (*roachpb.ExplainTreePlanNode, error) {
	node := roachpb.ExplainTreePlanNode{}

	nameAttr, err := jsonVal.FetchValKey("Name")
	if err != nil {
		return nil, err
	}

	if nameAttr != nil {
		str, err := nameAttr.AsText()
		if err != nil {
			return nil, err
		}
		node.Name = *str
	}

	iter, err := jsonVal.ObjectIter()
	if err != nil {
		return nil, err
	}

	if iter == nil {
		return nil, errors.New("unable to deconstruct json object")
	}

	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		if key == "Name" {
			// We already handled the name, so we skip it.
			continue
		}

		if key == "Children" {
			for childIdx := 0; childIdx < value.Len(); childIdx++ {
				childJSON, err := value.FetchValIdx(childIdx)
				if err != nil {
					return nil, err
				}
				if childJSON != nil {
					child, err := JSONToExplainTreePlanNode(childJSON)
					if err != nil {
						return nil, err
					}
					node.Children = append(node.Children, child)
				}
			}
		} else {
			str, err := value.AsText()
			if err != nil {
				return nil, err
			}
			node.Attrs = append(node.Attrs, &roachpb.ExplainTreePlanNode_Attr{
				Key:   key,
				Value: *str,
			})
		}
	}

	return &node, nil
}
