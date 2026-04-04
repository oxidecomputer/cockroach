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

package keys

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/stretchr/testify/require"
)

func TestRewriteKeyToTenantPrefix(t *testing.T) {
	for _, tc := range []struct {
		oldTenant, newTenant uint64
		suffix               string
	}{
		{1, 1, "abc"},
		{1, 2, "acc"},
		{2, 1, "abc"},
		{1, 1 << 15, "abc"},
		{1 << 15, 1, "abc"},
		{2, 1 << 15, "abc"},
		{1 << 15, 2, "abc"},
		{1 << 15, 1 << 15, "abc"},
		{1 << 13, 1 << 15, "abc"},
		{1 << 15, 1 << 13, "abc"},
	} {
		t.Run(fmt.Sprintf("%d to %d %s", tc.oldTenant, tc.newTenant, tc.suffix), func(t *testing.T) {
			old := append(MakeSQLCodec(roachpb.MakeTenantID(tc.oldTenant)).TablePrefix(5), tc.suffix...)
			new := MakeTenantPrefix(roachpb.MakeTenantID(tc.newTenant))
			got, err := rewriteKeyToTenantPrefix(old, new)
			require.NoError(t, err)

			expect := append(MakeSQLCodec(roachpb.MakeTenantID(tc.newTenant)).TablePrefix(5), tc.suffix...)
			require.Equal(t, expect, got)
		})
	}
}
