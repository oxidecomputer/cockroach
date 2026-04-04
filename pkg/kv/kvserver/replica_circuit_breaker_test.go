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

package kvserver

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/liveness"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/echotest"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/redact"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/raft/v3"
)

func TestReplicaUnavailableError(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	var repls roachpb.ReplicaSet
	repls.AddReplica(roachpb.ReplicaDescriptor{NodeID: 1, StoreID: 10, ReplicaID: 100})
	repls.AddReplica(roachpb.ReplicaDescriptor{NodeID: 2, StoreID: 20, ReplicaID: 200})
	desc := roachpb.NewRangeDescriptor(10, roachpb.RKey("a"), roachpb.RKey("z"), repls)
	var ba roachpb.BatchRequest
	ba.Add(&roachpb.RequestLeaseRequest{})
	lm := liveness.IsLiveMap{
		1: liveness.IsLiveMapEntry{IsLive: true},
	}
	ts, err := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
	require.NoError(t, err)
	wrappedErr := errors.New("probe failed")
	rs := raft.Status{}
	ctx := context.Background()
	err = errors.DecodeError(ctx, errors.EncodeError(ctx, replicaUnavailableError(
		wrappedErr, desc, desc.Replicas().AsProto()[0], lm, &rs, hlc.Timestamp{WallTime: ts.UnixNano()}),
	))
	require.True(t, errors.Is(err, wrappedErr), "%+v", err)
	echotest.Require(t, string(redact.Sprint(err)), testutils.TestDataPath(t, "replica_unavailable_error.txt"))
}
