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

// Package oidext contains oids that are not in `github.com/lib/pq/oid`
// as they are not shipped by default with postgres.
// As CRDB does not support extensions, we'll need to automatically assign
// a few OIDs of our own.
package oidext

import (
	"fmt"
	"testing"

	"github.com/lib/pq/oid"
	"github.com/stretchr/testify/require"
)

func TestTypeName(t *testing.T) {
	testCases := []struct {
		o            oid.Oid
		expectedName string
		expectedOk   bool
	}{
		{oid.T_int4, "INT4", true},
		{T_geometry, "GEOMETRY", true},
		{oid.Oid(99988199), "", false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("oid:%d", tc.o), func(t *testing.T) {
			name, ok := TypeName(tc.o)
			require.Equal(t, tc.expectedName, name)
			require.Equal(t, tc.expectedOk, ok)
		})
	}
}

// TestOIDsMaxValue ensures that any predefined OID values are less than
// the set maximum for OIDs.
func TestOIDsMaxValue(t *testing.T) {
	for oid, name := range oid.TypeName {
		if oid >= CockroachPredefinedOIDMax {
			t.Fatalf("oid %d for type %s greater than max %d", oid, name, CockroachPredefinedOIDMax)
		}
	}
	for oid, name := range ExtensionTypeName {
		if oid >= CockroachPredefinedOIDMax {
			t.Fatalf("oid %d for type %s greater than max %d", oid, name, CockroachPredefinedOIDMax)
		}
	}
}
