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

package persistedsqlstats_test

import (
	"context"
	gosql "database/sql"
	"fmt"
	"sync"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/skip"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
)

func BenchmarkConcurrentSelect1(b *testing.B) {
	skip.UnderShort(b)
	defer log.Scope(b).Close(b)
	ctx := context.Background()

	for _, numOfConcurrentConn := range []int{24, 48, 64} {
		b.Run(fmt.Sprintf("concurrentConn=%d", numOfConcurrentConn), func(b *testing.B) {
			s, db, _ := serverutils.StartServer(b, base.TestServerArgs{})
			sqlServer := s.SQLServer().(*sql.Server)
			defer s.Stopper().Stop(ctx)

			starter := make(chan struct{})
			latencyChan := make(chan float64, numOfConcurrentConn)
			defer close(latencyChan)

			var wg sync.WaitGroup
			for connIdx := 0; connIdx < numOfConcurrentConn; connIdx++ {
				sqlConn, err := db.Conn(ctx)
				if err != nil {
					b.Fatalf("unexpected error creating db conn: %s", err)
				}
				wg.Add(1)

				go func(conn *gosql.Conn, idx int) {
					defer wg.Done()
					runner := sqlutils.MakeSQLRunner(conn)
					<-starter

					start := timeutil.Now()
					for i := 0; i < b.N; i++ {
						runner.Exec(b, "SELECT 1")
					}
					duration := timeutil.Since(start)
					latencyChan <- float64(duration.Milliseconds()) / float64(b.N)
				}(sqlConn, connIdx)
			}

			close(starter)
			wg.Wait()

			var totalLat float64
			for i := 0; i < numOfConcurrentConn; i++ {
				totalLat += <-latencyChan
			}
			b.ReportMetric(
				sqlServer.ServerMetrics.
					StatsMetrics.
					SQLTxnStatsCollectionOverhead.
					Snapshot().Mean(),
				"overhead(ns/op)")
		})
	}
}
