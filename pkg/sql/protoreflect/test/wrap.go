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

package protoreflecttest

import (
	"encoding/json"

	"github.com/cockroachdb/cockroach/pkg/sql/protoreflect"
	"github.com/gogo/protobuf/jsonpb"
)

// SecretMessage is a message which should be redacted.
const SecretMessage = "secret message"

// RedactedMessage is the string the SecretMessage should be redacted to.
const RedactedMessage = "nothing to see here"

// MarshalJSONPB implements jsonpb.JSONPBMarshaler interface.
func (m Inner) MarshalJSONPB(marshaller *jsonpb.Marshaler) ([]byte, error) {
	if protoreflect.ShouldRedact(marshaller) && m.Value == SecretMessage {
		m.Value = RedactedMessage
	}
	return json.Marshal(m)
}
