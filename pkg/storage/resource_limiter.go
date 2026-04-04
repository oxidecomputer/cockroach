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

package storage

import (
	"time"

	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
)

// ResourceLimiterOptions is defining limits for resource limiter to restrict number
// of iterations.
type ResourceLimiterOptions struct {
	MaxRunTime time.Duration
}

// ResourceLimitReached indicates which resource threshold is exceeded.
// Soft threshold is advisory as to stop building results at convenient
// point, hard threshold is when iteration should stop straight away.
type ResourceLimitReached int

//go:generate stringer -type ResourceLimitReached
const (
	ResourceLimitNotReached  ResourceLimitReached = 0
	ResourceLimitReachedSoft ResourceLimitReached = 1
	ResourceLimitReachedHard ResourceLimitReached = 2
)

// ResourceLimiter provides a facility to stop cooperative long-running operation.
type ResourceLimiter interface {
	// IsExhausted returns true when limited resource is exhausted. Iterator is
	// checking the exhaustion status of resource limiter every time it advances
	// to the next underlying key value pair.
	IsExhausted() ResourceLimitReached
}

// TimeResourceLimiter provides limiter based on wall clock time.
type TimeResourceLimiter struct {
	softMaxRunTime time.Duration
	hardMaxRunTime time.Duration
	startTime      time.Time
	ts             timeutil.TimeSource
}

var _ ResourceLimiter = &TimeResourceLimiter{}

// NewResourceLimiter create new default resource limiter. Current implementation is wall clock time based.
// Timer starts as soon as limiter is created.
// If no limits are specified in opts nil is returned.
func NewResourceLimiter(opts ResourceLimiterOptions, ts timeutil.TimeSource) ResourceLimiter {
	if opts.MaxRunTime == 0 {
		return nil
	}
	softTimeLimit := time.Duration(float64(opts.MaxRunTime) * 0.9)
	return &TimeResourceLimiter{hardMaxRunTime: opts.MaxRunTime, softMaxRunTime: softTimeLimit, startTime: ts.Now(), ts: ts}
}

// IsExhausted implements ResourceLimiter interface.
func (l *TimeResourceLimiter) IsExhausted() ResourceLimitReached {
	timePassed := l.ts.Since(l.startTime)
	if timePassed < l.softMaxRunTime {
		return ResourceLimitNotReached
	}
	if timePassed < l.hardMaxRunTime {
		return ResourceLimitReachedSoft
	}
	return ResourceLimitReachedHard
}
