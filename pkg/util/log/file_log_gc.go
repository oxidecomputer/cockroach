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

package log

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync/atomic"
)

// gcDaemon runs the GC loop for the given logger.
func (l *fileSink) gcDaemon(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-l.gcNotify:
		}

		logging.mu.Lock()
		doGC := !logging.mu.disableDaemons
		logging.mu.Unlock()

		if doGC {
			l.gcOldFiles()
		}
	}
}

// gcOldFiles removes the "old" files that do not match
// the configured size and number threshold.
func (l *fileSink) gcOldFiles() {
	// This only lists the log files for the current logger (sharing the
	// prefix).
	dir, allFiles, err := l.listLogFiles()
	if err != nil {
		fmt.Fprintf(OrigStderr, "unable to GC log files: %s\n", err)
		return
	}

	if len(allFiles) == 0 {
		// Nothing to do.
		return
	}

	logFilesCombinedMaxSize := atomic.LoadInt64(&l.logFilesCombinedMaxSize)
	if logFilesCombinedMaxSize == 0 {
		// Nothing to do.
		return
	}

	files := selectFilesInGroup(allFiles, math.MaxInt64)
	if len(files) == 0 {
		return
	}
	// files is sorted with the newest log files first (which we want
	// to keep). Note that we always keep the most recent log file.
	sum := files[0].SizeBytes
	for _, f := range files[1:] {
		sum += f.SizeBytes
		if sum < logFilesCombinedMaxSize {
			continue
		}
		path := filepath.Join(dir, f.Name)
		if err := os.Remove(path); err != nil {
			fmt.Fprintln(OrigStderr, err)
		}
	}
}
