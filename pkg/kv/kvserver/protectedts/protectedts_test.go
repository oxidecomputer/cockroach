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

package protectedts

import "testing"

// TestProtectedTimestamps exists mostly to defeat the unused linter.
func TestProtectedTimestamps(t *testing.T) {
	var (
		_ Provider
		_ Cache
		_ Storage
		_ = EmptyCache(nil)
		_ = ErrNotExists
		_ = ErrExists
		_ = PollInterval
		_ = MaxBytes
		_ = MaxSpans
	)
}
