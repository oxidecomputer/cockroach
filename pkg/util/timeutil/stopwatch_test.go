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

package timeutil_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/stretchr/testify/require"
)

// TestStopWatchStart makes sure that consequent calls to Start do not reset
// the internal startedAt time.
func TestStopWatchStart(t *testing.T) {
	timeSource := timeutil.NewTestTimeSource()
	w := timeutil.NewTestStopWatch(timeSource.Now)

	w.Start()
	timeSource.Advance()
	w.Start()
	timeSource.Advance()
	w.Stop()

	expected, actual := timeSource.Elapsed(), w.Elapsed()
	require.Equal(t, expected, actual)
}

// TestStopWatchStop makes sure that only the first call to Stop changes the
// state of the stop watch.
func TestStopWatchStop(t *testing.T) {
	timeSource := timeutil.NewTestTimeSource()
	w := timeutil.NewTestStopWatch(timeSource.Now)

	w.Start()
	timeSource.Advance()
	w.Stop()

	expected, actual := timeSource.Elapsed(), w.Elapsed()
	require.Equal(t, expected, actual)

	timeSource.Advance()
	w.Stop()
	require.Equal(t, actual, w.Elapsed(), "consequent call to StopWatch.Stop changed the elapsed time")
}

// TestStopWatchElapsed makes sure that the stop watch records the elapsed time
// correctly.
func TestStopWatchElapsed(t *testing.T) {
	timeSource := timeutil.NewTestTimeSource()
	w := timeutil.NewTestStopWatch(timeSource.Now)
	expected := time.Duration(10)

	w.Start()
	for i := int64(0); i < int64(expected); i++ {
		timeSource.Advance()
	}
	w.Stop()

	require.Equal(t, expected, w.Elapsed())
}

// TestStopWatchConcurrentUsage makes sure that the stop watch is safe for
// concurrent usage.
func TestStopWatchConcurrentUsage(t *testing.T) {
	defer leaktest.AfterTest(t)()

	const testDuration = time.Second
	const maxSleepTime = testDuration / 100
	const numGoroutines = 10

	// All operations that we can do on the stop watch.
	const (
		start int = iota
		stop
		elapsed
		numOperations
	)

	w := timeutil.NewStopWatch()
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		// Spin up multiple goroutines that will be using the stop watch
		// concurrently.
		go func() {
			defer wg.Done()
			rng, _ := randutil.NewTestRand()
			var timeSpent time.Duration
			for timeSpent < testDuration {
				// Sleep some random time, up to maxSleepTime.
				toSleep := time.Duration(float64(maxSleepTime) * rng.Float64())
				time.Sleep(toSleep)
				timeSpent += toSleep
				// Pick the operation randomly.
				switch operation := rng.Intn(numOperations); operation {
				case start:
					w.Start()
				case stop:
					w.Stop()
				case elapsed:
					_ = w.Elapsed()
				default:
					panic(fmt.Sprintf("unexpected operation %d", operation))
				}
			}
		}()
	}
	wg.Wait()
}
