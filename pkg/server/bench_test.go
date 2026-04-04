// Copyright 2014 The Cockroach Authors.
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

package server

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/testutils/skip"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/tracing"
	"google.golang.org/grpc/metadata"
)

func BenchmarkSetupSpanForIncomingRPC(b *testing.B) {
	skip.UnderDeadlock(b, "span reuse triggers false-positives in the deadlock detector")
	defer leaktest.AfterTest(b)()

	for _, tc := range []struct {
		name      string
		traceInfo bool
		grpcMeta  bool
	}{
		{name: "traceInfo", traceInfo: true},
		{name: "grpcMeta", grpcMeta: true},
		{name: "no parent"},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			ctx := context.Background()
			tr := tracing.NewTracerWithOpt(ctx,
				tracing.WithTracingMode(tracing.TracingModeActiveSpansRegistry),
				tracing.WithSpanReusePercent(100))
			parentSpan := tr.StartSpan("parent")
			defer parentSpan.Finish()

			ba := &roachpb.BatchRequest{}
			if tc.traceInfo {
				ba.TraceInfo = parentSpan.Meta().ToProto()
			} else if tc.grpcMeta {
				traceCarrier := tracing.MapCarrier{
					Map: make(map[string]string),
				}
				tr.InjectMetaInto(parentSpan.Meta(), traceCarrier)
				ctx = metadata.NewIncomingContext(ctx, metadata.New(traceCarrier.Map))
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, sp := setupSpanForIncomingRPC(ctx, roachpb.SystemTenantID, ba, tr)
				sp.finish(ctx, nil /* br */)
			}
		})
	}
}
