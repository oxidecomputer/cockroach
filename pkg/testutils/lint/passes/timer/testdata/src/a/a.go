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

package a

import "github.com/cockroachdb/cockroach/pkg/util/timeutil"

func init() {
	timer := timeutil.NewTimer()
	for {
		timer.Reset(0)
		select {
		case <-timer.C:
			timer.Read = true
		}
	}
	for {
		timer.Reset(0)
		select {
		case <-timer.C: // want `must set timer.Read = true after reading from timer.C \(see timeutil/timer.go\)`
		}
	}
}
