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

package settings

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/redact"
	"github.com/stretchr/testify/require"
)

func TestEncodedValueSafeFormat(t *testing.T) {
	defer leaktest.AfterTest(t)()

	for _, tc := range []struct {
		rv       EncodedValue
		redacted string
		regular  string
	}{
		{
			rv: EncodedValue{
				Value: "asdf",
				Type:  "b",
			},

			regular:  `"asdf" (b)`,
			redacted: `‹"asdf"› (b)`,
		},
	} {
		t.Run(tc.regular, func(t *testing.T) {
			require.Equal(t, tc.regular, tc.rv.String())
			require.Equal(t, tc.redacted, string(redact.Sprint(tc.rv)))
		})
	}
}
