// Copyright 2016 The Cockroach Authors.
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

package limit

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/util/quotapool"
	"github.com/cockroachdb/cockroach/pkg/util/tracing"
	"github.com/cockroachdb/errors"
	"github.com/gogo/protobuf/types"
)

// ConcurrentRequestLimiter wraps a simple semaphore, adding a tracing span when
// a request is forced to wait.
type ConcurrentRequestLimiter struct {
	spanName string
	sem      *quotapool.IntPool
}

// Reservation is an allocation from a limiter which should be released once the
// limited task has been completed.
type Reservation interface {
	Release()
}

// MakeConcurrentRequestLimiter creates a ConcurrentRequestLimiter.
func MakeConcurrentRequestLimiter(spanName string, limit int) ConcurrentRequestLimiter {
	return ConcurrentRequestLimiter{
		spanName: spanName,
		sem:      quotapool.NewIntPool(spanName, uint64(limit)),
	}
}

// Begin attempts to reserve a spot in the pool, blocking if needed until the
// one is available or the context is canceled and adding a tracing span if it
// is forced to block.
func (l *ConcurrentRequestLimiter) Begin(ctx context.Context) (Reservation, error) {
	if err := ctx.Err(); err != nil {
		return nil, errors.Wrap(err, "limiter begin")
	}

	res, err := l.sem.TryAcquire(ctx, 1)
	if errors.Is(err, quotapool.ErrNotEnoughQuota) {
		var span *tracing.Span
		ctx, span = tracing.ChildSpan(ctx, l.spanName)
		defer span.Finish()
		span.RecordStructured(&types.StringValue{Value: fmt.Sprintf("%d requests are waiting", l.sem.Len())})
		res, err = l.sem.Acquire(ctx, 1)
	}
	return res, err
}

// SetLimit adjusts the size of the pool.
func (l *ConcurrentRequestLimiter) SetLimit(newLimit int) {
	l.sem.UpdateCapacity(uint64(newLimit))
}
