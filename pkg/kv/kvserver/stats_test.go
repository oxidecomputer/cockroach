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

package kvserver

import (
	"context"
	"reflect"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/stateloader"
	"github.com/cockroachdb/cockroach/pkg/storage/enginepb"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/stop"
)

func TestRangeStatsInit(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	ctx := context.Background()
	tc := testContext{}
	stopper := stop.NewStopper()
	defer stopper.Stop(ctx)
	tc.Start(ctx, t, stopper)
	ms := enginepb.MVCCStats{
		LiveBytes:       1,
		KeyBytes:        2,
		ValBytes:        3,
		IntentBytes:     4,
		LiveCount:       5,
		KeyCount:        6,
		ValCount:        7,
		IntentCount:     8,
		IntentAge:       9,
		GCBytesAge:      10,
		LastUpdateNanos: 11,
	}
	rsl := stateloader.Make(tc.repl.RangeID)
	if err := rsl.SetMVCCStats(ctx, tc.engine, &ms); err != nil {
		t.Fatal(err)
	}
	loadMS, err := rsl.LoadMVCCStats(ctx, tc.engine)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ms, loadMS) {
		t.Errorf("mvcc stats mismatch %+v != %+v", ms, loadMS)
	}
}
