// Copyright 2020 The Cockroach Authors.
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

package optional_test

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/optional"
	"github.com/stretchr/testify/require"
)

func TestUint(t *testing.T) {
	var v optional.Uint
	require.False(t, v.HasValue())
	require.Equal(t, uint64(0), v.Value())
	require.Equal(t, v.String(), "<unset>")

	v.Set(0)
	require.True(t, v.HasValue())
	require.Equal(t, uint64(0), v.Value())
	require.Equal(t, v.String(), "0")

	v.Set(10)
	require.True(t, v.HasValue())
	require.Equal(t, uint64(10), v.Value())
	require.Equal(t, v.String(), "10")

	v.Add(100)
	require.True(t, v.HasValue())
	require.Equal(t, uint64(110), v.Value())
	require.Equal(t, v.String(), "110")

	v.Clear()
	require.False(t, v.HasValue())
	require.Equal(t, uint64(0), v.Value())
	require.Equal(t, v.String(), "<unset>")

	var other optional.Uint

	v.MaybeAdd(other)
	require.False(t, v.HasValue())
	require.Equal(t, uint64(0), v.Value())
	require.Equal(t, v.String(), "<unset>")

	other.Set(100)
	v.MaybeAdd(other)
	require.True(t, v.HasValue())
	require.Equal(t, uint64(100), v.Value())
	require.Equal(t, v.String(), "100")
}
