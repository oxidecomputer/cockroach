// Copyright 2023 The Cockroach Authors.
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

package randgen

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/errors"
)

// TestSeedTypes verifies that at least one representative type is included into
// SeedTypes for all (with a few exceptions) type families.
func TestSeedTypes(t *testing.T) {
	defer leaktest.AfterTest(t)()

	noFamilyRepresentative := make(map[types.Family]struct{})
loop:
	for id := range types.Family_name {
		familyID := types.Family(id)
		switch familyID {
		case types.EnumFamily:
			// Enums need to created separately.
			continue loop
		case types.UnknownFamily, types.AnyFamily:
			// These are not included on purpose.
			continue loop
		}
		noFamilyRepresentative[familyID] = struct{}{}
	}
	for _, typ := range SeedTypes {
		delete(noFamilyRepresentative, typ.Family())
	}
	if len(noFamilyRepresentative) > 0 {
		s := "no representative for "
		for f := range noFamilyRepresentative {
			s += fmt.Sprintf("%s (%d) ", types.Family_name[int32(f)], f)
		}
		t.Fatal(errors.Errorf("%s", s))
	}
}
