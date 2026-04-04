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

package spanconfigsqlwatcher

import (
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/protectedts/ptpb"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/systemschema"
	"github.com/cockroachdb/cockroach/pkg/sql/rowenc/valueside"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
	"github.com/cockroachdb/errors"
)

// protectedTimestampDecoder decodes rows from the system.protected_ts_records
// table.
type protectedTimestampDecoder struct {
	alloc   tree.DatumAlloc
	decoder valueside.Decoder
}

// newProtectedTimestampDecoder instantiates a protectedTimestampDecoder.
func newProtectedTimestampDecoder() *protectedTimestampDecoder {
	columns := systemschema.ProtectedTimestampsRecordsTable.PublicColumns()
	return &protectedTimestampDecoder{
		decoder: valueside.MakeDecoder(columns),
	}
}

// DecodeRow decodes a row of the system.protected_ts_records table.
func (d *protectedTimestampDecoder) decode(kv roachpb.KeyValue) (target ptpb.Target, _ error) {
	if !kv.Value.IsPresent() {
		return ptpb.Target{},
			errors.AssertionFailedf("missing value for key in system.protected_ts_records: %v", kv)
	}

	// The columns after the `id` field are stored as a family.
	bytes, err := kv.Value.GetTuple()
	if err != nil {
		return ptpb.Target{}, err
	}

	datums, err := d.decoder.Decode(&d.alloc, bytes)
	if err != nil {
		return ptpb.Target{}, err
	}

	if len(datums) != 8 {
		return ptpb.Target{}, errors.AssertionFailedf("expected len(datums) == 8, but found %d", len(datums))
	}

	if t := datums[7]; t != tree.DNull {
		targetBytes := tree.MustBeDBytes(t)
		if err := protoutil.Unmarshal([]byte(targetBytes), &target); err != nil {
			return ptpb.Target{}, errors.Wrap(err, "failed to unmarshal target")
		}
	}

	return target, nil
}

// TestingProtectedTimestampDecoderFn constructs a protectedTimestampDecoder and
// exposes its decode method.
func TestingProtectedTimestampDecoderFn() func(roachpb.KeyValue) (ptpb.Target, error) {
	return newProtectedTimestampDecoder().decode
}
