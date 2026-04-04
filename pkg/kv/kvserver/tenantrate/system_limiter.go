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

package tenantrate

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/multitenant/tenantcostmodel"
)

// systemLimiter implements Limiter for the use of tracking metrics for the
// system tenant. It does not actually perform any rate-limiting.
type systemLimiter struct {
	tenantMetrics
}

func (s systemLimiter) Wait(ctx context.Context, reqInfo tenantcostmodel.RequestInfo) error {
	if reqInfo.IsWrite() {
		s.writeBatchesAdmitted.Inc(1)
		s.writeRequestsAdmitted.Inc(reqInfo.WriteCount())
		s.writeBytesAdmitted.Inc(reqInfo.WriteBytes())
	}
	return nil
}

func (s systemLimiter) RecordRead(ctx context.Context, respInfo tenantcostmodel.ResponseInfo) {
	if respInfo.IsRead() {
		s.readBatchesAdmitted.Inc(1)
		s.readRequestsAdmitted.Inc(respInfo.ReadCount())
		s.readBytesAdmitted.Inc(respInfo.ReadBytes())
	}
}

var _ Limiter = (*systemLimiter)(nil)
