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

package testutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var someFunc = func(n int) int {
	return n - 1
}

type someType struct {
	funcField func(int) int
}

func TestTestingHookWithGlobal(t *testing.T) {
	require.Equal(t, 9, someFunc(10))
	restoreHook := TestingHook(&someFunc, func(n int) int {
		return n + 1
	})
	require.Equal(t, 11, someFunc(10))
	restoreHook()
	require.Equal(t, 9, someFunc(10))
}

func TestTestingHookWithStruct(t *testing.T) {
	s := someType{
		funcField: func(n int) int {
			return n - 1
		},
	}
	require.Equal(t, 9, s.funcField(10))
	restoreHook := TestingHook(&s.funcField, func(n int) int {
		return n + 1
	})
	require.Equal(t, 11, s.funcField(10))
	restoreHook()
	require.Equal(t, 9, s.funcField(10))
}
