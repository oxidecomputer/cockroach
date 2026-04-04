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

package buildutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCrdbTestOn(t *testing.T) {
	// Sanity-check: make sure CrdbTestBuild is set. This should be true for
	// any test.
	require.True(t, CrdbTestBuild)
}
