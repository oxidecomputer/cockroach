// Copyright 2020 The Cockroach Authors.
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

package serverpb

import (
	context "context"

	"github.com/cockroachdb/cockroach/pkg/util/errorutil"
)

// SQLStatusServer is a smaller version of the serverpb.StatusInterface which
// includes only the methods used by the SQL subsystem.
type SQLStatusServer interface {
	ListSessions(context.Context, *ListSessionsRequest) (*ListSessionsResponse, error)
	ListLocalSessions(context.Context, *ListSessionsRequest) (*ListSessionsResponse, error)
	CancelQuery(context.Context, *CancelQueryRequest) (*CancelQueryResponse, error)
	CancelQueryByKey(context.Context, *CancelQueryByKeyRequest) (*CancelQueryByKeyResponse, error)
	CancelSession(context.Context, *CancelSessionRequest) (*CancelSessionResponse, error)
	ListContentionEvents(context.Context, *ListContentionEventsRequest) (*ListContentionEventsResponse, error)
	ListLocalContentionEvents(context.Context, *ListContentionEventsRequest) (*ListContentionEventsResponse, error)
	ResetSQLStats(context.Context, *ResetSQLStatsRequest) (*ResetSQLStatsResponse, error)
	CombinedStatementStats(context.Context, *CombinedStatementsStatsRequest) (*StatementsResponse, error)
	Statements(context.Context, *StatementsRequest) (*StatementsResponse, error)
	StatementDetails(context.Context, *StatementDetailsRequest) (*StatementDetailsResponse, error)
	ListDistSQLFlows(context.Context, *ListDistSQLFlowsRequest) (*ListDistSQLFlowsResponse, error)
	ListLocalDistSQLFlows(context.Context, *ListDistSQLFlowsRequest) (*ListDistSQLFlowsResponse, error)
	Profile(context.Context, *ProfileRequest) (*JSONResponse, error)
	IndexUsageStatistics(context.Context, *IndexUsageStatisticsRequest) (*IndexUsageStatisticsResponse, error)
	ResetIndexUsageStats(context.Context, *ResetIndexUsageStatsRequest) (*ResetIndexUsageStatsResponse, error)
	TableIndexStats(context.Context, *TableIndexStatsRequest) (*TableIndexStatsResponse, error)
	UserSQLRoles(context.Context, *UserSQLRolesRequest) (*UserSQLRolesResponse, error)
	TxnIDResolution(context.Context, *TxnIDResolutionRequest) (*TxnIDResolutionResponse, error)
	TransactionContentionEvents(context.Context, *TransactionContentionEventsRequest) (*TransactionContentionEventsResponse, error)
	NodesList(context.Context, *NodesListRequest) (*NodesListResponse, error)
}

// OptionalNodesStatusServer is a StatusServer that is only optionally present
// inside the SQL subsystem. In practice, it is present on the system tenant,
// and not present on "regular" tenants.
type OptionalNodesStatusServer struct {
	w errorutil.TenantSQLDeprecatedWrapper // stores serverpb.StatusServer
}

// MakeOptionalNodesStatusServer initializes and returns an
// OptionalNodesStatusServer. The provided server will be returned via
// OptionalNodesStatusServer() if and only if it is not nil.
func MakeOptionalNodesStatusServer(s NodesStatusServer) OptionalNodesStatusServer {
	return OptionalNodesStatusServer{
		// Return the status server from OptionalSQLStatusServer() only if one was provided.
		// We don't have any calls to .Deprecated().
		w: errorutil.MakeTenantSQLDeprecatedWrapper(s, s != nil /* exposed */),
	}
}

// NodesStatusServer is an endpoint that allows the SQL subsystem
// to observe node descriptors.
// It is unavailable to tenants.
type NodesStatusServer interface {
	ListNodesInternal(context.Context, *NodesRequest) (*NodesResponse, error)
}

// RegionsServer is the subset of the serverpb.StatusInterface that is used
// by the SQL system to query for available regions.
// It is available for tenants.
type RegionsServer interface {
	Regions(context.Context, *RegionsRequest) (*RegionsResponse, error)
}

// TenantStatusServer is the subset of the serverpb.StatusInterface that is
// used by tenants to query for debug information, such as tenant-specific
// range reports.
//
// It is available for all tenants.
type TenantStatusServer interface {
	TenantRanges(context.Context, *TenantRangesRequest) (*TenantRangesResponse, error)
}

// OptionalNodesStatusServer returns the wrapped NodesStatusServer, if it is
// available. If it is not, an error referring to the optionally supplied issues
// is returned.
func (s *OptionalNodesStatusServer) OptionalNodesStatusServer(
	issue int,
) (NodesStatusServer, error) {
	v, err := s.w.OptionalErr(issue)
	if err != nil {
		return nil, err
	}
	return v.(NodesStatusServer), nil
}
