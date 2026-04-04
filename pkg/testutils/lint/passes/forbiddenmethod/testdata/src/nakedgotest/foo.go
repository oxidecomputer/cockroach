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

package nakedgotest

func toot() {
	// The nolint comment below should have no effect.
	// For some reason though it shows up in the CommentMap
	// for the *ast.GoStmt, though. No idea why.

	//nolint:nakedgo should not help anyone
}

func A() {
	//nolint: I'm a noop
	go func() {}()                              // want `Use of go keyword not allowed, use a Stopper instead`
	go toot()                                   // want `Use of go keyword not allowed, use a Stopper instead`
	go /* nolint: nakedgo not helping */ toot() // want `Use of go keyword not allowed, use a Stopper instead`
	/* nolint: nakedgo nope */ go toot() // want `Use of go keyword not allowed, use a Stopper instead`
	//nolint:nakedgo nope, this one neither

	go func() {}() // want `Use of go keyword not allowed, use a Stopper instead`

	go func() {}() //nolint:nakedgo

	go toot() //nolint:nakedgo

	go func() { /* want `Use of go keyword not allowed, use a Stopper instead` */ //nolint:nakedgo
		_ = 0
	}()

	go func() { // want `Use of go keyword not allowed, use a Stopper instead`
		_ = 0 //nolint:nakedgo
	}()

	// Finally, doing it right!

	go func() {
		_ = 0
	}() //nolint:nakedgo
}
