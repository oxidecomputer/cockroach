// Copyright 2016 The Cockroach Authors.
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

package envutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvOrDefault(t *testing.T) {
	const def = 123
	os.Clearenv()
	// These tests are mostly an excuse to exercise otherwise unused code.
	// TODO(knz): Test everything.
	if act := EnvOrDefaultBytes("COCKROACH_X", def); act != def {
		t.Errorf("expected %d, got %d", def, act)
	}
	if act := EnvOrDefaultInt("COCKROACH_X", def); act != def {
		t.Errorf("expected %d, got %d", def, act)
	}
}

func TestCheckVarName(t *testing.T) {
	t.Run("checkVarName", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			valid bool
		}{
			{
				name:  "abc123",
				valid: false,
			},
			{
				name:  "ABC 123",
				valid: false,
			},
			{
				name:  "@&) 123",
				valid: false,
			},
			{
				name:  "ABC123",
				valid: true,
			},
			{
				name:  "ABC_123",
				valid: true,
			},
		} {
			func() {
				defer func() {
					r := recover()
					if !tc.valid && r == nil {
						t.Errorf("expected panic for name %q, got none", tc.name)
					} else if tc.valid && r != nil {
						t.Errorf("unexpected panic for name %q, got %q", tc.name, r)
					}
				}()

				checkVarName(tc.name)
			}()
		}
	})

	t.Run("checkInternalVarName", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			valid bool
		}{
			{
				name:  "ABC_123",
				valid: false,
			},
			{
				name:  "COCKROACH_X",
				valid: true,
			},
		} {
			func() {
				defer func() {
					r := recover()
					if !tc.valid && r == nil {
						t.Errorf("expected panic for name %q, got none", tc.name)
					} else if tc.valid && r != nil {
						t.Errorf("unexpected panic for name %q, got %q", tc.name, r)
					}
				}()

				checkInternalVarName(tc.name)
			}()
		}
	})

	t.Run("checkExternalVarName", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			valid bool
		}{
			{
				name:  "COCKROACH_X",
				valid: false,
			},
			{
				name:  "ABC_123",
				valid: true,
			},
		} {
			func() {
				defer func() {
					r := recover()
					if !tc.valid && r == nil {
						t.Errorf("expected panic for name %q, got none", tc.name)
					} else if tc.valid && r != nil {
						t.Errorf("unexpected panic for name %q, got %q", tc.name, r)
					}
				}()

				checkExternalVarName(tc.name)
			}()
		}
	})
}

func TestTestSetEnvExists(t *testing.T) {
	key := "COCKROACH_ENVUTIL_TESTSETTING"
	require.NoError(t, os.Setenv(key, "before"))

	ClearEnvCache()
	value, ok := EnvString(key, 0)
	require.True(t, ok)
	require.Equal(t, value, "before")

	cleanup := TestSetEnv(t, key, "testvalue")
	value, ok = EnvString(key, 0)
	require.True(t, ok)
	require.Equal(t, value, "testvalue")

	cleanup()

	value, ok = EnvString(key, 0)
	require.True(t, ok)
	require.Equal(t, value, "before")
}

func TestTestSetEnvDoesNotExist(t *testing.T) {
	key := "COCKROACH_ENVUTIL_TESTSETTING"
	require.NoError(t, os.Unsetenv(key))

	ClearEnvCache()
	_, ok := EnvString(key, 0)
	require.False(t, ok)

	cleanup := TestSetEnv(t, key, "foo")
	value, ok := EnvString(key, 0)
	require.True(t, ok)
	require.Equal(t, value, "foo")

	cleanup()

	_, ok = EnvString(key, 0)
	require.False(t, ok)
}
