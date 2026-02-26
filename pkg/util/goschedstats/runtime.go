// Copyright 2026 Oxide Computer Company
//
// Use of this software is governed by the Apache License, Version 2.0,
// included in the file licenses/APL.txt.

package goschedstats

import "runtime/metrics"

func numRunnableGoroutines() (numRunnable int, numProcs int) {
	sample := make([]metrics.Sample, 2)
	sample[0].Name = "/sched/goroutines/runnable:goroutines"
	sample[1].Name = "/sched/gomaxprocs:threads"
	metrics.Read(sample)
	return int(sample[0].Value.Uint64()), int(sample[1].Value.Uint64())
}
