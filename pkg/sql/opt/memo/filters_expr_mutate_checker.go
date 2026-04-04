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

package memo

import (
	"github.com/cockroachdb/cockroach/pkg/util/buildutil"
	"github.com/cockroachdb/errors"
)

// FiltersExprMutateChecker is used to check if a FiltersExpr has been
// erroneously mutated. This code is called in crdb_test builds so that the
// check is run for tests, but the overhead is not incurred for non-test builds.
type FiltersExprMutateChecker struct {
	hasher hasher
	hash   internHash
}

// Init initializes a FiltersExprMutateChecker with the original filters.
func (fmc *FiltersExprMutateChecker) Init(filters FiltersExpr) {
	if !buildutil.CrdbTestBuild {
		return
	}
	// This initialization pattern ensures that fields are not unwittingly
	// reused. Field reuse must be explicit.
	*fmc = FiltersExprMutateChecker{}
	fmc.hasher.Init()
	fmc.hasher.HashFiltersExpr(filters)
	fmc.hash = fmc.hasher.hash
}

// CheckForMutation panics if the given filters are not equal to the filters
// passed for the previous Init function call.
func (fmc *FiltersExprMutateChecker) CheckForMutation(filters FiltersExpr) {
	if !buildutil.CrdbTestBuild {
		return
	}
	fmc.hasher.Init()
	fmc.hasher.HashFiltersExpr(filters)
	if fmc.hash != fmc.hasher.hash {
		panic(errors.AssertionFailedf("filters should not be mutated"))
	}
}
