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

package rttanalysis

import (
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/log"
)

type testingB interface {
	testing.TB
	N() int
	ResetTimer()
	StopTimer()
	StartTimer()
	ReportMetric(float64, string)
	Run(string, func(testingB))

	// logScope is used to wrap log.Scope and make it available only in
	// appropriate contexts.
	logScope() (getDirectory func() string, close func())
}

type bShim struct {
	*testing.B
}

var _ testingB = bShim{}

func (b bShim) logScope() (getDirectory func() string, close func()) {
	sc := log.Scope(b)
	return sc.GetDirectory, func() { sc.Close(b) }
}
func (b bShim) N() int { return b.B.N }
func (b bShim) Run(name string, f func(b testingB)) {
	b.B.Run(name, func(b *testing.B) {
		f(bShim{b})
	})
}

// tShim is used by the expectation test's testing.T to appear as though
// it is a testing.B and can be used capture the output. The object also
// suppresses creation of a new log.Scope in order to permit parallel
// execution.
type tShim struct {
	*testing.T
	scope   *log.TestLogScope
	results *resultSet
}

var _ testingB = tShim{}

func (ts tShim) logScope() (getDirectory func() string, close func()) {
	return ts.scope.GetDirectory, func() {}
}
func (ts tShim) GetDirectory() string {
	return ts.scope.GetDirectory()
}
func (ts tShim) N() int      { return 2 }
func (ts tShim) ResetTimer() {}
func (ts tShim) StopTimer()  {}
func (ts tShim) StartTimer() {}
func (ts tShim) ReportMetric(f float64, s string) {
	if s == roundTripsMetric {
		ts.results.add(benchmarkResult{
			name:   ts.Name(),
			result: int(f),
		})
	}
}
func (ts tShim) Name() string {
	// Trim the name of the outermost test off the beginning of the name.
	tn := ts.T.Name()
	return tn[strings.Index(tn, "/")+1:]
}
func (ts tShim) Run(s string, f func(testingB)) {
	ts.T.Run(s, func(t *testing.T) {
		f(tShim{results: ts.results, T: t, scope: ts.scope})
	})
}
