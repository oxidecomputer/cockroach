// Copyright 2017 The Cockroach Authors.
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

//go:build !linux
// +build !linux

package hlc

import (
	"context"

	"github.com/cockroachdb/errors"
)

// ClockSource contains the handle of the clock device as well as the
// clock id.
type ClockSource struct {
}

// UnixNano is not used on platforms other than Linux
func (p ClockSource) UnixNano() int64 {
	panic(errors.New("clock device not supported on this platform"))
}

// MakeClockSource us not used on platforms other than Linux
func MakeClockSource(_ context.Context, _ string) (ClockSource, error) {
	return ClockSource{}, errors.New("clock device not supported on this platform")
}
