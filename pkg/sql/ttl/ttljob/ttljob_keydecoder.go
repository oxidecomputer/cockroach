// Copyright 2022 The Cockroach Authors.
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

package ttljob

import (
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/rowenc"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/errors"
)

const (
	stripTenantPrefixErrorFmt      = "error decoding tenant prefix of %x"
	decodePartialTableIDIndexIDFmt = "error decoding table/index ID of %x"
	encDatumFromBufferFmt          = "error decoding EncDatum of %x"
	ensureDecodedFmt               = "error ensuring encoding of %x"
)

// keyToDatums translates a Key on a span for a table to the appropriate datums.
func keyToDatums(
	key roachpb.Key, codec keys.SQLCodec, pkTypes []*types.T, alloc *tree.DatumAlloc,
) (tree.Datums, error) {

	// Decode the datums ourselves, instead of using rowenc.DecodeKeyVals.
	// We cannot use rowenc.DecodeKeyVals because we may not have the entire PK
	// as the key for the span (e.g. a PK (a, b) may only be split on (a)).
	partialKey, err := codec.StripTenantPrefix(key)
	if err != nil {
		// Convert key to []byte to prevent hex encoding output of Key.String().
		return nil, errors.Wrapf(err, stripTenantPrefixErrorFmt, []byte(key))
	}
	partialKey, _, _, err = rowenc.DecodePartialTableIDIndexID(partialKey)
	if err != nil {
		// Convert key to []byte to prevent hex encoding output of Key.String().
		return nil, errors.Wrapf(err, decodePartialTableIDIndexIDFmt, []byte(key))
	}
	encDatums := make([]rowenc.EncDatum, 0, len(pkTypes))
	for len(partialKey) > 0 && len(encDatums) < len(pkTypes) {
		i := len(encDatums)
		// We currently assume all PRIMARY KEY columns are ascending, and block
		// creation otherwise.
		enc := descpb.DatumEncoding_ASCENDING_KEY
		var val rowenc.EncDatum
		val, partialKey, err = rowenc.EncDatumFromBuffer(pkTypes[i], enc, partialKey)
		if err != nil {
			// Convert key to []byte to prevent hex encoding output of Key.String().
			return nil, errors.Wrapf(err, encDatumFromBufferFmt, []byte(key))
		}
		encDatums = append(encDatums, val)
	}

	datums := make(tree.Datums, len(encDatums))
	for i, encDatum := range encDatums {
		if err := encDatum.EnsureDecoded(pkTypes[i], alloc); err != nil {
			// Convert key to []byte to prevent hex encoding output of Key.String().
			return nil, errors.Wrapf(err, ensureDecodedFmt, []byte(key))
		}
		datums[i] = encDatum.Datum
	}
	return datums, nil
}
