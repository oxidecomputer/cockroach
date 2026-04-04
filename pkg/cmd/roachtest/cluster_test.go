// Copyright 2018 The Cockroach Authors.
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

package main

import (
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/option"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/spec"
	test2 "github.com/cockroachdb/cockroach/pkg/cmd/roachtest/test"
	"github.com/cockroachdb/cockroach/pkg/roachprod/logger"
	"github.com/cockroachdb/cockroach/pkg/util/version"
	"github.com/stretchr/testify/assert"
)

func TestClusterNodes(t *testing.T) {
	c := &clusterImpl{spec: spec.MakeClusterSpec(spec.GCE, "", 10)}
	opts := func(opts ...option.Option) []option.Option {
		return opts
	}
	testCases := []struct {
		opts     []option.Option
		expected string
	}{
		{opts(), ""},
		{opts(c.All()), ":1-10"},
		{opts(c.Range(1, 2)), ":1-2"},
		{opts(c.Range(2, 5)), ":2-5"},
		{opts(c.All(), c.Range(2, 5)), ":1-10"},
		{opts(c.Range(2, 5), c.Range(7, 9)), ":2-5,7-9"},
		{opts(c.Range(2, 5), c.Range(6, 8)), ":2-8"},
		{opts(c.Node(2), c.Node(4), c.Node(6)), ":2,4,6"},
		{opts(c.Node(2), c.Node(3), c.Node(4)), ":2-4"},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			nodes := c.MakeNodes(tc.opts...)
			if tc.expected != nodes {
				t.Fatalf("expected %s, but found %s", tc.expected, nodes)
			}
		})
	}
}

type testWrapper struct {
	*testing.T
	l *logger.Logger
}

func (t testWrapper) BuildVersion() *version.Version {
	panic("implement me")
}

func (t testWrapper) Cockroach() string {
	return "./dummy-path/to/cockroach"
}

func (t testWrapper) DeprecatedWorkload() string {
	return "./dummy-path/to/workload"
}

func (t testWrapper) IsBuildVersion(s string) bool {
	panic("implement me")
}

func (t testWrapper) Spec() interface{} {
	panic("implement me")
}

func (t testWrapper) VersionsBinaryOverride() map[string]string {
	panic("implement me")
}

func (t testWrapper) Progress(f float64) {
	panic("implement me")
}

func (t testWrapper) WorkerStatus(args ...interface{}) {
}

func (t testWrapper) WorkerProgress(f float64) {
	panic("implement me")
}

var _ test2.Test = testWrapper{}

// ArtifactsDir is part of the test.Test interface.
func (t testWrapper) ArtifactsDir() string {
	return ""
}

// PerfArtifactsDir is part of the test.Test interface.
func (t testWrapper) PerfArtifactsDir() string {
	return ""
}

// logger is part of the testI interface.
func (t testWrapper) L() *logger.Logger {
	return t.l
}

// Status is part of the testI interface.
func (t testWrapper) Status(args ...interface{}) {}

func TestClusterMachineType(t *testing.T) {
	testCases := []struct {
		machineType      string
		expectedCPUCount int
	}{
		// AWS machine types
		{"m5.large", 2},
		{"m5.xlarge", 4},
		{"m5.2xlarge", 8},
		{"m5.4xlarge", 16},
		{"m5.12xlarge", 48},
		{"m5.24xlarge", 96},
		{"m5d.large", 2},
		{"m5d.xlarge", 4},
		{"m5d.2xlarge", 8},
		{"m5d.4xlarge", 16},
		{"m5d.12xlarge", 48},
		{"m5d.24xlarge", 96},
		{"c5d.large", 2},
		{"c5d.xlarge", 4},
		{"c5d.2xlarge", 8},
		{"c5d.4xlarge", 16},
		{"c5d.9xlarge", 36},
		{"c5d.18xlarge", 72},
		// GCE machine types
		{"n1-standard-1", 1},
		{"n1-standard-2", 2},
		{"n1-standard-4", 4},
		{"n1-standard-8", 8},
		{"n1-standard-16", 16},
		{"n1-standard-32", 32},
		{"n1-standard-64", 64},
		{"n1-standard-96", 96},
	}
	for _, tc := range testCases {
		t.Run(tc.machineType, func(t *testing.T) {
			cpuCount := MachineTypeToCPUs(tc.machineType)
			if tc.expectedCPUCount != cpuCount {
				t.Fatalf("expected %d CPUs, but found %d", tc.expectedCPUCount, cpuCount)
			}
		})
	}
}

func TestCmdLogFileName(t *testing.T) {
	ts := time.Date(2000, 1, 1, 15, 4, 12, 0, time.Local)

	const exp = `run_150412.000000000_n1,3-4,9_cockroach-bla-foo-ba`
	nodes := option.NodeListOption{1, 3, 4, 9}
	assert.Equal(t,
		exp,
		cmdLogFileName(ts, nodes, "./cockroach", "bla", "--foo", "bar"),
	)
	assert.Equal(t,
		exp,
		cmdLogFileName(ts, nodes, "./cockroach bla --foo bar"),
	)
}
