// Copyright 2017 The Cockroach Authors.
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

//go:build !deadlock && race
// +build !deadlock,race

package syncutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAssertHeld(t *testing.T) {
	type mutex interface {
		Lock()
		Unlock()
		AssertHeld()
	}

	testCases := []struct {
		m mutex
	}{
		{&Mutex{}},
		{&RWMutex{}},
	}
	for _, c := range testCases {
		// The normal, successful case.
		c.m.Lock()
		c.m.AssertHeld()
		c.m.Unlock()

		// The unsuccessful case.
		require.PanicsWithValue(t, "mutex is not write locked", c.m.AssertHeld)
	}
}

func TestAssertRHeld(t *testing.T) {
	var m RWMutex

	// The normal, successful case.
	m.RLock()
	m.AssertRHeld()
	m.RUnlock()

	// The normal case with two readers.
	m.RLock()
	m.RLock()
	m.AssertRHeld()
	m.RUnlock()
	m.RUnlock()

	// The case where a write lock is held.
	m.Lock()
	m.AssertRHeld()
	m.Unlock()

	// The unsuccessful case with no readers.
	require.PanicsWithValue(t, "mutex is not read locked", m.AssertRHeld)
}
