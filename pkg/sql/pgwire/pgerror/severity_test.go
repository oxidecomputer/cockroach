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

package pgerror

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/require"
)

func TestSeverity(t *testing.T) {
	testCases := []struct {
		err              error
		expectedSeverity string
	}{
		{WithSeverity(fmt.Errorf("notice me"), "NOTICE ME"), "NOTICE ME"},
		{WithSeverity(WithSeverity(fmt.Errorf("notice me"), "IGNORE ME"), "NOTICE ME"), "NOTICE ME"},
		{WithSeverity(WithCandidateCode(fmt.Errorf("notice me"), pgcode.FeatureNotSupported), "NOTICE ME"), "NOTICE ME"},
		{New(pgcode.Uncategorized, "i am an error"), "ERROR"},
		{WithCandidateCode(WithSeverity(errors.Newf("i am not an error"), "NOT AN ERROR"), pgcode.System), "NOT AN ERROR"},
		{fmt.Errorf("something else"), "ERROR"},
	}

	for _, tc := range testCases {
		t.Run(tc.err.Error(), func(t *testing.T) {
			severity := GetSeverity(tc.err)
			require.Equal(t, tc.expectedSeverity, severity)
		})
	}
}
