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

package tpcc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewResult(t *testing.T) {
	// Ensure you don't get panics when calling common methods
	// on a trivial Result that doesn't have any data attached.
	res := NewResult(1000, 0, 0, nil)
	require.Error(t, res.FailureError())
	require.Zero(t, res.Efficiency())
	require.Zero(t, res.TpmC())
}
