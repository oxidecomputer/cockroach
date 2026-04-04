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

package ptutil

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/spanconfig"
	"github.com/cockroachdb/cockroach/pkg/spanconfig/spanconfigptsreader"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/errors"
)

// TestingVerifyProtectionTimestampExistsOnSpans refreshes the PTS state in KV and
// ensures a protection at the given protectionTimestamp exists for all the
// supplied spans.
func TestingVerifyProtectionTimestampExistsOnSpans(
	ctx context.Context,
	t *testing.T,
	srv serverutils.TestServerInterface,
	ptsReader spanconfig.ProtectedTSReader,
	protectionTimestamp hlc.Timestamp,
	spans roachpb.Spans,
) error {
	testutils.SucceedsSoon(t, func() error {
		if err := spanconfigptsreader.TestingRefreshPTSState(
			ctx, t, ptsReader, srv.Clock().Now(),
		); err != nil {
			return err
		}
		for _, sp := range spans {
			timestamps, _, err := ptsReader.GetProtectionTimestamps(ctx, sp)
			if err != nil {
				return err
			}
			found := false
			for _, ts := range timestamps {
				if ts.Equal(protectionTimestamp) {
					found = true
					break
				}
			}
			if !found {
				return errors.Newf("protection timestamp %s does not exist on span %s", protectionTimestamp, sp)
			}
		}
		return nil
	})
	return nil
}
