// Copyright 2022 The Cockroach Authors.
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

//go:build nightly
// +build nightly

package metamorphic

func TestPebbleEquivalenceNightly(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	skip.UnderRace(t)
	if *opCount < 1000000 {
		oldOpCount := *opCount
		// Override number of operations to at least 1 million.
		*opCount = 1000000

		defer func() {
			*opCount = oldOpCount
		}()
	}

	runPebbleEquivalenceTest(t)
}
