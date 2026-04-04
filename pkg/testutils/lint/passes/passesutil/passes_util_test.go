// Copyright 2020 The Cockroach Authors.
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

package passesutil_test

import (
	"path/filepath"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/build/bazel"
	"github.com/cockroachdb/cockroach/pkg/testutils/lint/passes/forbiddenmethod"
	"github.com/cockroachdb/cockroach/pkg/testutils/lint/passes/unconvert"
	"github.com/cockroachdb/cockroach/pkg/testutils/skip"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func init() {
	if bazel.BuiltWithBazel() {
		bazel.SetGoEnv()
	}
}

// Use tests from other packages to also test this package. This ensures
// that if that code changes, somebody will look here. Also it allows for
// coverage checking here.

func requireNotEmpty(t *testing.T, path string) {
	t.Helper()
	files, err := filepath.Glob(path)
	require.NoError(t, err)
	require.NotEmpty(t, files)
}

func getTestdataForPackage(t *testing.T, pkg string) string {
	if bazel.BuiltWithBazel() {
		runfiles, err := bazel.RunfilesPath()
		require.NoError(t, err)
		return filepath.Join(runfiles, "pkg", "testutils", "lint", "passes", pkg, "testdata")
	}
	return filepath.Join("..", pkg, "testdata")
}

func TestDescriptorMarshal(t *testing.T) {
	skip.UnderStress(t)
	testdata, err := filepath.Abs(getTestdataForPackage(t, "forbiddenmethod"))
	require.NoError(t, err)
	requireNotEmpty(t, testdata)
	analysistest.TestData = func() string { return testdata }
	analysistest.Run(t, testdata, forbiddenmethod.DescriptorMarshalAnalyzer, "descmarshaltest")
}

func TestUnconvert(t *testing.T) {
	skip.UnderStress(t)
	testdata, err := filepath.Abs(getTestdataForPackage(t, "unconvert"))
	require.NoError(t, err)
	requireNotEmpty(t, testdata)
	analysistest.TestData = func() string { return testdata }
	analysistest.Run(t, testdata, unconvert.Analyzer, "a")
}
