// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package bootstrap

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/config/zonepb"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/datadriven"
	"github.com/stretchr/testify/require"
)

func TestInitialValuesToString(t *testing.T) {
	defer leaktest.AfterTest(t)()
	datadriven.Walk(t, testutils.TestDataPath(t), func(t *testing.T, path string) {
		datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
			codec := keys.SystemSQLCodec
			switch d.Cmd {
			case "system":
				break

			case "tenant":
				const dummyTenantID = 12345
				codec = keys.MakeSQLCodec(roachpb.MakeTenantID(dummyTenantID))

			default:
				t.Fatalf("unexpected command %q", d.Cmd)
			}
			var expectedHash string
			d.ScanArgs(t, "hash", &expectedHash)
			ms := MakeMetadataSchema(codec, zonepb.DefaultZoneConfigRef(), zonepb.DefaultSystemZoneConfigRef())
			result := InitialValuesToString(ms)
			h := sha256.Sum256([]byte(result))
			if actualHash := hex.EncodeToString(h[:]); expectedHash != actualHash {
				t.Errorf(`Unexpected hash value %s for %s.
If you're seeing this error message, this means that the bootstrapped system
schema has changed. Assuming that this is expected:
- If this occurred during development on the main branch, rewrite the expected
  test output and the hash value and move on.
- If this occurred during development of a patch for a release branch, make
  very sure that the underlying change really is expected and is backward-
  compatible and is absolutely necessary. If that's the case, then there are
  hardcoded literals in the main development branch as well as any subsequent
  release branches that need to be updated also.`, actualHash, d.Cmd)
			}
			return result
		})
	})
}

func TestRoundTripInitialValuesStringRepresentation(t *testing.T) {
	t.Run("system", func(t *testing.T) {
		roundTripInitialValuesStringRepresentation(t, 0 /* tenantID */)
	})
	t.Run("tenant", func(t *testing.T) {
		const dummyTenantID = 54321
		roundTripInitialValuesStringRepresentation(t, dummyTenantID)
	})
	t.Run("tenants", func(t *testing.T) {
		const dummyTenantID1, dummyTenantID2 = 54321, 12345
		require.Equal(t,
			InitialValuesToString(makeMetadataSchema(dummyTenantID1)),
			InitialValuesToString(makeMetadataSchema(dummyTenantID2)),
		)
	})
}

func roundTripInitialValuesStringRepresentation(t *testing.T, tenantID uint64) {
	ms := makeMetadataSchema(tenantID)
	expectedKVs, expectedSplits := ms.GetInitialValues()
	actualKVs, actualSplits, err := InitialValuesFromString(ms.codec, InitialValuesToString(ms))
	require.NoError(t, err)
	require.Len(t, actualKVs, len(expectedKVs))
	require.Len(t, actualSplits, len(expectedSplits))
	for i, actualKV := range actualKVs {
		expectedKV := expectedKVs[i]
		require.EqualValues(t, expectedKV, actualKV)
	}
	for i, actualSplit := range actualSplits {
		expectedSplit := expectedSplits[i]
		require.EqualValues(t, expectedSplit, actualSplit)
	}
}

func makeMetadataSchema(tenantID uint64) MetadataSchema {
	codec := keys.SystemSQLCodec
	if tenantID > 0 {
		codec = keys.MakeSQLCodec(roachpb.MakeTenantID(tenantID))
	}
	return MakeMetadataSchema(codec, zonepb.DefaultZoneConfigRef(), zonepb.DefaultSystemZoneConfigRef())
}
