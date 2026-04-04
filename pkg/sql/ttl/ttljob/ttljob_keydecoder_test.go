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
	"fmt"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/encoding"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

func TestKeyToDatums(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	const tenantID = 111

	testCases := []struct {
		desc                 string
		keyBytes             []byte
		errorFmt             string
		expectedErrorMessage string
		expectedDatums       tree.Datums
	}{
		{
			desc:                 "StripTenantPrefix error",
			keyBytes:             []byte{1, 2, 3},
			errorFmt:             stripTenantPrefixErrorFmt,
			expectedErrorMessage: `error decoding tenant prefix of 010203: invalid tenant id prefix: /Local/"` + "\u0002\u0003" + `"`,
		},
		{
			desc:                 "DecodePartialTableIDIndexID error",
			keyBytes:             []byte{254, 246, tenantID},
			errorFmt:             decodePartialTableIDIndexIDFmt,
			expectedErrorMessage: `error decoding table/index ID of fef66f: insufficient bytes to decode uvarint value`,
		},
		{
			desc:                 "EncDatumFromBuffer error",
			keyBytes:             []byte{254, 246, tenantID, 1, 1, 5},
			errorFmt:             encDatumFromBufferFmt,
			expectedErrorMessage: `error decoding EncDatum of fef66f010105: slice too short for float (1)`,
		},
		{
			desc:                 "EnsureDecoded error",
			keyBytes:             []byte{254, 246, tenantID, 1, 1, 1},
			errorFmt:             ensureDecodedFmt,
			expectedErrorMessage: `error ensuring encoding of fef66f010101: error decoding 1 bytes: insufficient bytes to decode varint value: ""`,
		},
		{
			desc:           "success",
			keyBytes:       encoding.EncodeVarintAscending([]byte{254, 246, tenantID, 1, 1}, 100),
			expectedDatums: []tree.Datum{tree.NewDInt(100)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tenantID := roachpb.MakeTenantID(tenantID)
			codec := keys.MakeSQLCodec(tenantID)
			keyBytes := tc.keyBytes
			var alloc tree.DatumAlloc
			datums, err := keyToDatums(keyBytes, codec, []*types.T{types.Int}, &alloc)
			expectedErrorMessage := tc.expectedErrorMessage
			if expectedErrorMessage != "" {
				require.Error(t, err)
				actualErrorMessage := err.Error()
				require.Equal(t, expectedErrorMessage, actualErrorMessage)
				parts := strings.Split(actualErrorMessage, ":")
				// Verify that the hex encoded key from the error message matches the original key.
				var errorKeyBytes []byte
				_, err := fmt.Sscanf(parts[0], tc.errorFmt, &errorKeyBytes)
				require.NoError(t, err)
				require.Equal(t, keyBytes, errorKeyBytes)
			}
			expectedDatums := tc.expectedDatums
			if expectedDatums != nil {
				require.Equal(t, expectedDatums, datums)
			}
		})
	}
}
