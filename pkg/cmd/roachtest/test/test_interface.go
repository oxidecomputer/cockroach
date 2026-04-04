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

package test

import (
	"github.com/cockroachdb/cockroach/pkg/roachprod/logger"
	"github.com/cockroachdb/cockroach/pkg/util/version"
)

// Test is the interface through which roachtests interact with the
// test harness.
type Test interface {
	Cockroach() string // path to main cockroach binary
	Name() string
	BuildVersion() *version.Version
	IsBuildVersion(string) bool // "vXX.YY"
	Helper()
	// Spec() returns the *registry.TestSpec as an interface{}.
	//
	// TODO(tbg): cleaning this up is mildly tricky. TestSpec has the Run field
	// which depends both on `test` (and `cluster`, though this matters less), so
	// we get cyclic imports. We should split up `Run` off of `TestSpec` and have
	// `TestSpec` live in `spec` to avoid this problem, but this requires a pass
	// through all registered roachtests to change how they register the test.
	Spec() interface{}
	VersionsBinaryOverride() map[string]string
	Skip(args ...interface{})
	Skipf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(string, ...interface{})
	FailNow()
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Failed() bool
	ArtifactsDir() string
	// PerfArtifactsDir is the directory on cluster nodes in which perf artifacts
	// reside. Upon success this directory is copied into test's ArtifactsDir from
	// each node in the cluster.
	PerfArtifactsDir() string
	L() *logger.Logger
	Progress(float64)
	Status(args ...interface{})
	WorkerStatus(args ...interface{})
	WorkerProgress(float64)

	// DeprecatedWorkload returns the path to the workload binary.
	// Don't use this, invoke `./cockroach workload` instead.
	DeprecatedWorkload() string
}
