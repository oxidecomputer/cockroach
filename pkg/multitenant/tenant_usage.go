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

package multitenant

import (
	"context"
	"time"

	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/metric"
)

// TenantUsageServer is an interface through which tenant usage is reported and
// controlled, used on the host server side. Its implementation lives in the
// tenantcostserver CCL package.
type TenantUsageServer interface {
	// TokenBucketRequest implements the TokenBucket API of the roachpb.Internal
	// service. Used to to service requests coming from tenants (through the
	// kvtenant.Connector)
	TokenBucketRequest(
		ctx context.Context, tenantID roachpb.TenantID, in *roachpb.TokenBucketRequest,
	) *roachpb.TokenBucketResponse

	// ReconfigureTokenBucket updates a tenant's token bucket settings.
	//
	// Arguments:
	//
	//  - availableRU is the amount of Request Units that the tenant can consume at
	//    will. Also known as "burst RUs".
	//
	//  - refillRate is the amount of Request Units per second that the tenant
	//    receives.
	//
	//  - maxBurstRU is the maximum amount of Request Units that can be accumulated
	//    from the refill rate, or 0 if there is no limit.
	//
	//  - asOf is a timestamp; the reconfiguration request is assumed to be based on
	//    the consumption at that time. This timestamp is used to compensate for any
	//    refill that would have happened in the meantime.
	//
	//  - asOfConsumedRequestUnits is the total number of consumed RUs based on
	//    which the reconfiguration values were calculated (i.e. at the asOf time).
	//    It is used to adjust availableRU with the consumption that happened in the
	//    meantime.
	//
	ReconfigureTokenBucket(
		ctx context.Context,
		txn *kv.Txn,
		tenantID roachpb.TenantID,
		availableRU float64,
		refillRate float64,
		maxBurstRU float64,
		asOf time.Time,
		asOfConsumedRequestUnits float64,
	) error

	// Metrics returns the top-level metrics.
	Metrics() metric.Struct
}
