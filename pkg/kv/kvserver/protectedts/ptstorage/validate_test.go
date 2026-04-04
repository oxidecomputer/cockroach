// Copyright 2019 The Cockroach Authors.
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

package ptstorage

import (
	"context"
	"strconv"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/protectedts"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/protectedts/ptpb"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateRecordForProtect(t *testing.T) {
	target := ptpb.MakeClusterTarget()
	for i, tc := range []struct {
		r   *ptpb.Record
		err error
	}{
		{
			r: &ptpb.Record{
				ID:        uuid.MakeV4().GetBytes(),
				Timestamp: hlc.Timestamp{WallTime: 1, Logical: 1},
				MetaType:  "job",
				Meta:      []byte("junk"),
				Target:    target,
			},
			err: nil,
		},
		{
			r: &ptpb.Record{
				Timestamp: hlc.Timestamp{WallTime: 1, Logical: 1},
				MetaType:  "job",
				Meta:      []byte("junk"),
				Target:    target,
			},
			err: errZeroID,
		},
		{
			r: &ptpb.Record{
				ID:       uuid.MakeV4().GetBytes(),
				MetaType: "job",
				Meta:     []byte("junk"),
				Target:   target,
			},
			err: errZeroTimestamp,
		},
		{
			r: &ptpb.Record{
				ID:        uuid.MakeV4().GetBytes(),
				Timestamp: hlc.Timestamp{WallTime: 1, Logical: 1},
				Meta:      []byte("junk"),
				Target:    target,
			},
			err: errInvalidMeta,
		},
		{
			r: &ptpb.Record{
				ID:        uuid.MakeV4().GetBytes(),
				Timestamp: hlc.Timestamp{WallTime: 1, Logical: 1},
				MetaType:  "job",
				Meta:      []byte("junk"),
			},
			err: errNilTarget,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			st := cluster.MakeTestingClusterSettings()
			require.Equal(t, validateRecordForProtect(context.Background(), tc.r, st,
				&protectedts.TestingKnobs{}), tc.err)
		})

		// Test that prior to the `AlterSystemProtectedTimestampAddColumn` migration
		// we validate that records have a non-nil `Spans` field.
		t.Run("errEmptySpans", func(t *testing.T) {
			r := &ptpb.Record{
				ID:        uuid.MakeV4().GetBytes(),
				Timestamp: hlc.Timestamp{WallTime: 1, Logical: 1},
				MetaType:  "job",
				Meta:      []byte("junk"),
				Target:    target,
			}
			st := cluster.MakeTestingClusterSettings()
			require.Equal(t, validateRecordForProtect(context.Background(), r, st,
				&protectedts.TestingKnobs{DisableProtectedTimestampForMultiTenant: true}), errEmptySpans)
		})
	}
}
