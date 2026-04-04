// Copyright 2018 The Cockroach Authors.
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

package tracing

import (
	"context"
	"testing"

	"github.com/cockroachdb/logtags"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestLogTags(t *testing.T) {
	tr := NewTracerWithOpt(context.Background(), WithTracingMode(TracingModeActiveSpansRegistry))
	sr := tracetest.NewSpanRecorder()
	otelTr := otelsdk.NewTracerProvider(otelsdk.WithSpanProcessor(sr)).Tracer("test")
	tr.SetOpenTelemetryTracer(otelTr)

	l := logtags.SingleTagBuffer("tag1", "val1")
	l = l.Add("tag2", "val2")
	sp1 := tr.StartSpan("foo", WithLogTags(l))
	sp1.SetRecordingType(RecordingVerbose)
	require.NoError(t, CheckRecordedSpans(sp1.FinishAndGetRecording(RecordingVerbose), `
		span: foo
			tags: _verbose=1 tag1=val1 tag2=val2
	`))
	{
		require.Len(t, sr.Ended(), 1)
		otelSpan := sr.Ended()[0]
		exp := []attribute.KeyValue{
			{Key: "tag1", Value: attribute.StringValue("val1")},
			{Key: "tag2", Value: attribute.StringValue("val2")},
		}
		require.Equal(t, exp, otelSpan.Attributes())
	}

	RegisterTagRemapping("tag1", "one")
	RegisterTagRemapping("tag2", "two")

	sp2 := tr.StartSpan("bar", WithLogTags(l))
	sp2.SetRecordingType(RecordingVerbose)
	require.NoError(t, CheckRecordedSpans(sp2.FinishAndGetRecording(RecordingVerbose), `
		span: bar
			tags: _verbose=1 one=val1 two=val2
	`))

	{
		require.Len(t, sr.Ended(), 2)
		otelSpan := sr.Ended()[1]
		exp := []attribute.KeyValue{
			{Key: "one", Value: attribute.StringValue("val1")},
			{Key: "two", Value: attribute.StringValue("val2")},
		}
		require.Equal(t, exp, otelSpan.Attributes())
	}

	sp3 := tr.StartSpan("baz", WithLogTags(l))
	sp3.SetRecordingType(RecordingVerbose)
	require.NoError(t, CheckRecordedSpans(sp3.FinishAndGetRecording(RecordingVerbose), `
		span: baz
			tags: _verbose=1 one=val1 two=val2
	`))
	{
		require.Len(t, sr.Ended(), 3)
		otelSpan := sr.Ended()[2]
		exp := []attribute.KeyValue{
			{Key: "one", Value: attribute.StringValue("val1")},
			{Key: "two", Value: attribute.StringValue("val2")},
		}
		require.Equal(t, exp, otelSpan.Attributes())
	}
}
