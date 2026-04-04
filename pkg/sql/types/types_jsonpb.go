// Copyright 2015 The Cockroach Authors.
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

package types

import (
	"bytes"

	"github.com/gogo/protobuf/jsonpb"
)

// This file contains logic to allow the *types.T to properly marshal to json.
// It is a separate file to make it straightforward to defeat the linter that
// refuses to allow one to call a method Marshal unless it's protoutil.Marshal.

// MarshalJSONPB marshals the T to json. This is necessary as otherwise
// this field will be lost to the crdb_internal.pb_to_json and the likes.
func (t *T) MarshalJSONPB(marshaler *jsonpb.Marshaler) ([]byte, error) {
	temp := *t
	if err := temp.downgradeType(); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := marshaler.Marshal(&buf, &temp.InternalType); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalJSONPB unmarshals the T to json. This is necessary as otherwise
// this field will be lost to the crdb_internal.json_to_pb and the likes.
func (t *T) UnmarshalJSONPB(unmarshaler *jsonpb.Unmarshaler, data []byte) error {
	if err := unmarshaler.Unmarshal(bytes.NewReader(data), &t.InternalType); err != nil {
		return err
	}
	return t.upgradeType()
}
