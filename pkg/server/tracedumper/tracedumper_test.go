// Copyright 2021 The Cockroach Authors.
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

package tracedumper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/server/dumpstore"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlutil"
	"github.com/cockroachdb/cockroach/pkg/sql/tests"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

func TestTraceDumperZipCreation(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	baseDir, dirCleanupFn := testutils.TempDir(t)
	defer dirCleanupFn()
	traceDir := filepath.Join(baseDir, "trace_dir")
	require.NoError(t, os.Mkdir(traceDir, 0755))

	baseTime := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	td := TraceDumper{
		currentTime: func() time.Time {
			return baseTime
		},
		store: dumpstore.NewStore(traceDir, nil, nil),
	}
	ctx := context.Background()
	params, _ := tests.CreateTestServerParams()
	s, _, _ := serverutils.StartServer(t, params)
	defer s.Stopper().Stop(ctx)

	filename := "foo"
	td.Dump(ctx, filename, 123, s.InternalExecutor().(sqlutil.InternalExecutor))
	expectedFilename := fmt.Sprintf("%s.%s.%s.zip", jobTraceDumpPrefix, baseTime.Format(timeFormat),
		filename)
	fullpath := td.store.GetFullPath(expectedFilename)
	_, err := os.Stat(fullpath)
	require.NoError(t, err)
}
