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

package joberror

import (
	"fmt"
	"testing"

	circuitbreaker "github.com/cockroachdb/circuitbreaker"
	"github.com/cockroachdb/cockroach/pkg/util/circuit"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/redact"
	"github.com/stretchr/testify/require"
)

func TestErrBreakerOpenIsRetriable(t *testing.T) {
	br := circuit.NewBreaker(circuit.Options{
		Name: redact.Sprint("Breaker"),
		AsyncProbe: func(_ func(error), done func()) {
			done() // never untrip
		},
		EventHandler: &circuit.EventLogger{Log: func(redact.StringBuilder) {}},
	})
	br.Report(errors.New("test error"))
	utilBreakderErr := br.Signal().Err()
	// NB: This matches the error that dial produces.
	dialErr := errors.Wrapf(circuitbreaker.ErrBreakerOpen, "unable to dial n%d", 9)

	for _, e := range []error{
		utilBreakderErr,
		dialErr,
	} {
		t.Run(fmt.Sprintf("%s", e), func(t *testing.T) {
			require.False(t, IsPermanentBulkJobError(e))
		})
	}
}
