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

//go:build !race
// +build !race

package tracing

import "runtime"

// setFinalizer registers a finalizer to run when the Span is garbage collected.
// Passing nil clears any previously registered finalizer.
func (sp *Span) setFinalizer(fn func(sp *Span)) {
	if fn == nil {
		// Avoid typed nil.
		runtime.SetFinalizer(sp, nil)
		return
	}
	runtime.SetFinalizer(sp, fn)
}
