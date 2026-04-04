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

package tracing

import (
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
	"github.com/cockroachdb/cockroach/pkg/util/tracing/tracingpb"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/redact"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

// TraceToJSON returns the string representation of the trace in JSON format.
//
// TraceToJSON assumes that the first span in the recording contains all the
// other spans.
func TraceToJSON(trace Recording) (string, error) {
	root := normalizeSpan(trace[0], trace)
	marshaller := jsonpb.Marshaler{
		Indent: "\t",
	}
	str, err := marshaller.MarshalToString(&root)
	if err != nil {
		return "", err
	}
	return str, nil
}

func normalizeSpan(s tracingpb.RecordedSpan, trace Recording) tracingpb.NormalizedSpan {
	var n tracingpb.NormalizedSpan
	n.Operation = s.Operation
	n.StartTime = s.StartTime
	n.Duration = s.Duration
	n.Tags = s.Tags
	n.Logs = s.Logs
	n.StructuredRecords = s.StructuredRecords

	for _, ss := range trace {
		if ss.ParentSpanID != s.SpanID {
			continue
		}
		n.Children = append(n.Children, normalizeSpan(ss, trace))
	}
	return n
}

// MessageToJSONString converts a protocol message into a JSON string. The
// emitDefaults flag dictates whether fields with zero values are rendered or
// not.
//
// TODO(andrei): It'd be nice if this function dealt with redactable vs safe
// fields, like EventPayload.AppendJSONFields does.
func MessageToJSONString(msg protoutil.Message, emitDefaults bool) (string, error) {
	// Convert to json.
	jsonEncoder := jsonpb.Marshaler{EmitDefaults: emitDefaults}
	msgJSON, err := jsonEncoder.MarshalToString(msg)
	if err != nil {
		return "", errors.Newf("error when converting %s to JSON string", proto.MessageName(msg))
	}

	return msgJSON, nil
}

// RedactAndTruncateError redacts the error and truncates the string
// representation of the error to a fixed length.
func RedactAndTruncateError(err error) string {
	maxErrLength := 250
	redactedErr := string(redact.Sprintf("%v", err))
	if len(redactedErr) < maxErrLength {
		maxErrLength = len(redactedErr)
	}
	return redactedErr[:maxErrLength]
}
