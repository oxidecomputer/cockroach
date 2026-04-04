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

package colexecerror_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/colexecerror"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

// TestCatchVectorizedRuntimeError verifies that the panic-catcher doesn't catch
// panics that originate outside of the vectorized engine and correctly
// annotates errors that are propagated via
// colexecerror.(Internal|Expected)Error methods.
func TestCatchVectorizedRuntimeError(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	// Setup multiple levels of catchers to ensure that the panic-catcher
	// doesn't fool itself into catching panics that the inner catcher emitted.
	require.Panics(t, func() {
		require.NoError(t, colexecerror.CatchVectorizedRuntimeError(func() {
			require.NoError(t, colexecerror.CatchVectorizedRuntimeError(func() {
				colexecerror.NonCatchablePanic(errors.New("should not be caught"))
			}))
		}))
	})

	const shouldBeCaughtText = "should be caught"
	shouldBeCaughtErr := errors.New(shouldBeCaughtText)
	const annotationText = "unexpected error from the vectorized engine"

	// Propagate an error as an internal one (this should add annotations to the
	// returned error).
	annotatedErr := colexecerror.CatchVectorizedRuntimeError(func() {
		colexecerror.InternalError(shouldBeCaughtErr)
	})
	require.NotNil(t, annotatedErr)
	require.True(t, strings.Contains(annotatedErr.Error(), shouldBeCaughtText))
	require.True(t, strings.Contains(annotatedErr.Error(), annotationText))

	// Propagate an error as an expected one (this should *not* add annotations
	// to the returned error).
	notAnnotatedErr := colexecerror.CatchVectorizedRuntimeError(func() {
		colexecerror.ExpectedError(shouldBeCaughtErr)
	})
	require.NotNil(t, notAnnotatedErr)
	require.True(t, strings.Contains(notAnnotatedErr.Error(), shouldBeCaughtText))
	require.False(t, strings.Contains(notAnnotatedErr.Error(), annotationText))
}

// TestNonCatchablePanicIsNotCaught verifies that panics emitted via
// NonCatchablePanic() method are not caught by the catcher.
func TestNonCatchablePanicIsNotCaught(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	require.Panics(t, func() {
		require.NoError(t, colexecerror.CatchVectorizedRuntimeError(func() {
			colexecerror.NonCatchablePanic("should panic")
		}))
	})
}
