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

package geos

import (
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/require"
)

func TestInitGEOS(t *testing.T) {
	t.Run("test no initGEOS paths", func(t *testing.T) {
		_, _, err := initGEOS([]string{})
		require.Error(t, err)
		require.Regexp(t, "Ensure you have the spatial libraries installed as per the instructions in .*install-cockroachdb-", strings.Join(errors.GetAllHints(err), "\n"))
	})

	t.Run("test invalid initGEOS paths", func(t *testing.T) {
		_, _, err := initGEOS([]string{"/invalid/path"})
		require.Error(t, err)
		require.Regexp(t, "Ensure you have the spatial libraries installed as per the instructions in .*install-cockroachdb-", strings.Join(errors.GetAllHints(err), "\n"))
	})

	t.Run("test valid initGEOS paths", func(t *testing.T) {
		ret, loc, err := initGEOS(findLibraryDirectories("", ""))
		require.NoError(t, err)
		require.NotEmpty(t, loc)
		require.NotNil(t, ret)
	})
}

func TestEnsureInit(t *testing.T) {
	// Fetch at least once.
	_, err := ensureInit(EnsureInitErrorDisplayPublic, "", "")
	require.NoError(t, err)

	fakeErr := errors.Newf("contain path info do not display me")
	defer func() { geosOnce.err = nil }()

	geosOnce.err = fakeErr
	_, err = ensureInit(EnsureInitErrorDisplayPrivate, "", "")
	require.Contains(t, err.Error(), fakeErr.Error())

	_, err = ensureInit(EnsureInitErrorDisplayPublic, "", "")
	require.Equal(t, errors.Newf("geos: this operation is not available").Error(), err.Error())
}
