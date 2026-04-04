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

package rpc

import (
	"github.com/gogo/protobuf/proto"
	// Used instead of gogo/protobuf/proto for the fallback case
	// to match the behavior of the upstream codec in
	// google.golang.org/grpc/encoding/proto that we are
	// replacing:
	//
	//   https://github.com/grpc/grpc-go/blob/7b167fd6eca1ab8f05ec14085d63197cacd41438/encoding/proto/proto.go
	//
	gproto "github.com/golang/protobuf/proto"
	"google.golang.org/grpc/encoding"
)

const name = "proto"

type codec struct{}

var _ encoding.Codec = codec{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	if pm, ok := v.(proto.Marshaler); ok {
		return pm.Marshal()
	}
	return gproto.Marshal(v.(gproto.Message))
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	if pm, ok := v.(proto.Unmarshaler); ok {
		return pm.Unmarshal(data)
	}
	return gproto.Unmarshal(data, v.(gproto.Message))
}

func (codec) Name() string {
	return name
}
