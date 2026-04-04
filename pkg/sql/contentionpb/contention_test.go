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

package contentionpb

import (
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/util"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/cockroach/pkg/util/uuid"
	"github.com/stretchr/testify/require"
)

func TestExtendedContentionEventHash(t *testing.T) {
	event1 := ExtendedContentionEvent{}
	event1.BlockingEvent.TxnMeta.ID = uuid.FastMakeV4()
	event1.WaitingTxnID = uuid.FastMakeV4()
	event1.CollectionTs = timeutil.Now()

	eventWithDifferentBlockingTxnID := event1
	eventWithDifferentBlockingTxnID.BlockingEvent.TxnMeta.ID = uuid.FastMakeV4()

	require.NotEqual(t, eventWithDifferentBlockingTxnID.Hash(), event1.Hash())

	eventWithDifferentWaitingTxnID := event1
	eventWithDifferentWaitingTxnID.WaitingTxnID = uuid.FastMakeV4()
	require.NotEqual(t, eventWithDifferentWaitingTxnID.Hash(), event1.Hash())

	eventWithDifferentCollectionTs := event1
	eventWithDifferentCollectionTs.CollectionTs = event1.CollectionTs.Add(time.Second)
	require.NotEqual(t, eventWithDifferentCollectionTs.Hash(), event1.Hash())
}

func TestHashingUUID(t *testing.T) {
	// Ensure that if two UUIDs are only different in the first or last 8 bytes,
	// they still produces different hash.
	uuid1 := uuid.UUID{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	}
	fnv1 := util.MakeFNV64()
	hashUUID(uuid1, &fnv1)

	uuid2 := uuid.UUID{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 17,
	}
	fnv2 := util.MakeFNV64()
	hashUUID(uuid2, &fnv2)

	uuid3 := uuid.UUID{
		0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	}
	fnv3 := util.MakeFNV64()
	hashUUID(uuid3, &fnv3)

	require.NotEqual(t, fnv1.Sum(), fnv2.Sum())
	require.NotEqual(t, fnv1.Sum(), fnv3.Sum())
	require.NotEqual(t, fnv2.Sum(), fnv3.Sum())
}
