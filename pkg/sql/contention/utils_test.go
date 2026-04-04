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

package contention

import "github.com/cockroachdb/cockroach/pkg/util/cache"

// SetSizeConstants updates the constants for the sizes of caches of the
// registries for tests. If any of the passed-in arguments is not positive, it
// is ignored. A cleanup function is returned to restore the original values.
func SetSizeConstants(indexMap, orderedKeyMap, numTxns int) func() {
	oldIndexMapMaxSize := indexMapMaxSize
	oldOrderedKeyMapMaxSize := orderedKeyMapMaxSize
	oldMaxNumTxns := maxNumTxns
	if indexMap > 0 {
		indexMapMaxSize = indexMap
	}
	if orderedKeyMap > 0 {
		orderedKeyMapMaxSize = orderedKeyMap
	}
	if numTxns > 0 {
		maxNumTxns = numTxns
	}
	return func() {
		indexMapMaxSize = oldIndexMapMaxSize
		orderedKeyMapMaxSize = oldOrderedKeyMapMaxSize
		maxNumTxns = oldMaxNumTxns
	}
}

// CalculateTotalNumContentionEvents returns the total number of contention
// events that r knows about.
func CalculateTotalNumContentionEvents(r *Registry) uint64 {
	numContentionEvents := uint64(0)
	r.indexMap.internalCache.Do(func(e *cache.Entry) {
		v := e.Value.(*indexMapValue)
		numContentionEvents += v.numContentionEvents
	})
	return numContentionEvents
}
