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

package install

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/roachprod/logger"
	"github.com/cockroachdb/cockroach/pkg/util/retry"
	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/require"
)

// TestRoachprodEnv tests the roachprodEnvRegex and roachprodEnvValue methods.
func TestRoachprodEnv(t *testing.T) {
	cases := []struct {
		clusterName string
		node        Node
		tag         string
		value       string
		regex       string
	}{
		{
			clusterName: "a",
			node:        1,
			tag:         "",
			value:       "1",
			regex:       `ROACHPROD=1[ \/]`,
		},
		{
			clusterName: "local-foo",
			node:        2,
			tag:         "",
			value:       "local-foo/2",
			regex:       `ROACHPROD=local-foo\/2[ \/]`,
		},
		{
			clusterName: "a",
			node:        3,
			tag:         "foo",
			value:       "3/foo",
			regex:       `ROACHPROD=3\/foo[ \/]`,
		},
		{
			clusterName: "a",
			node:        4,
			tag:         "foo/bar",
			value:       "4/foo/bar",
			regex:       `ROACHPROD=4\/foo\/bar[ \/]`,
		},
		{
			clusterName: "local-foo",
			node:        5,
			tag:         "tag",
			value:       "local-foo/5/tag",
			regex:       `ROACHPROD=local-foo\/5\/tag[ \/]`,
		},
	}

	for idx, tc := range cases {
		t.Run(fmt.Sprintf("%d", idx+1), func(t *testing.T) {
			var c SyncedCluster
			c.Name = tc.clusterName
			c.Tag = tc.tag
			if value := c.roachprodEnvValue(tc.node); value != tc.value {
				t.Errorf("expected value `%s`, got `%s`", tc.value, value)
			}
			if regex := c.roachprodEnvRegex(tc.node); regex != tc.regex {
				t.Errorf("expected regex `%s`, got `%s`", tc.regex, regex)
			}
		})
	}
}

func TestRunWithMaybeRetry(t *testing.T) {
	var testRetryOpts = retry.Options{
		InitialBackoff: 10 * time.Millisecond,
		Multiplier:     2,
		MaxBackoff:     1 * time.Second,
		// This will run a total of 3 times `runWithMaybeRetry`
		MaxRetries: 2,
	}

	l := nilLogger()

	attempt := 0
	cases := []struct {
		f                func() (*RunResultDetails, error)
		shouldRetryFn    func(*RunResultDetails) bool
		nilRetryOpts     bool
		expectedAttempts int
		shouldError      bool
	}{
		{ // 1. Happy path: no error, no retry required
			f: func() (*RunResultDetails, error) {
				return newResult(0), nil
			},
			expectedAttempts: 1,
			shouldError:      false,
		},
		{ // 2. Error, but with no retries
			f: func() (*RunResultDetails, error) {
				return newResult(1), nil
			},
			shouldRetryFn: func(*RunResultDetails) bool {
				return false
			},
			expectedAttempts: 1,
			shouldError:      true,
		},
		{ // 3. Error, but no retry function specified
			f: func() (*RunResultDetails, error) {
				return newResult(1), nil
			},
			expectedAttempts: 3,
			shouldError:      true,
		},
		{ // 4. Error, with retries exhausted
			f: func() (*RunResultDetails, error) {
				return newResult(255), nil
			},
			shouldRetryFn:    func(d *RunResultDetails) bool { return d.RemoteExitStatus == 255 },
			expectedAttempts: 3,
			shouldError:      true,
		},
		{ // 5. Eventual success after retries
			f: func() (*RunResultDetails, error) {
				attempt++
				if attempt == 3 {
					return newResult(0), nil
				}
				return newResult(255), nil
			},
			shouldRetryFn:    func(d *RunResultDetails) bool { return d.RemoteExitStatus == 255 },
			expectedAttempts: 3,
			shouldError:      false,
		},
		{ // 6. Error, runs once because nil retryOpts
			f: func() (*RunResultDetails, error) {
				return newResult(255), nil
			},
			nilRetryOpts:     true,
			expectedAttempts: 1,
			shouldError:      true,
		},
	}

	for idx, tc := range cases {
		attempt = 0
		t.Run(fmt.Sprintf("%d", idx+1), func(t *testing.T) {
			var retryOpts *RunRetryOpts
			if !tc.nilRetryOpts {
				retryOpts = newRunRetryOpts(testRetryOpts, tc.shouldRetryFn)
			}
			res, _ := runWithMaybeRetry(l, retryOpts, tc.f)

			require.Equal(t, tc.shouldError, res.Err != nil)
			require.Equal(t, tc.expectedAttempts, res.Attempt)

			if tc.shouldError && tc.expectedAttempts == 3 {
				require.True(t, errors.Is(res.Err, ErrAfterRetry))
			}
		})
	}
}

func newResult(exitCode int) *RunResultDetails {
	var err error
	if exitCode != 0 {
		err = errors.Newf("Error with exit code %v", exitCode)
	}
	return &RunResultDetails{RemoteExitStatus: exitCode, Err: err}
}

func nilLogger() *logger.Logger {
	lcfg := logger.Config{
		Stdout: io.Discard,
		Stderr: io.Discard,
	}
	l, err := lcfg.NewLogger("" /* path */)
	if err != nil {
		panic(err)
	}
	return l
}

func TestGenFilenameFromArgs(t *testing.T) {
	const exp = "mkdir-p-logsredacted"
	require.Equal(t, exp, GenFilenameFromArgs(20, "mkdir -p logs/redacted && ./cockroach"))
	require.Equal(t, exp, GenFilenameFromArgs(20, "mkdir", "-p logs/redacted", "&& ./cockroach"))
	require.Equal(t, exp, GenFilenameFromArgs(20, "mkdir    -p logs/redacted && ./cockroach    "))
}
