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

package sessiondatapb

// SequenceCache stores sequence values that have already been created in KV
// and are available to be given out as sequence numbers. Values for sequences
// are keyed by the descpb.ID of each sequence. These IDs are represented as
// uint32 to prevent an import cycle with the descpb package. The cache should
// only be accessed using the provided API.
//
// The cache ensures that values are invalidated when new descriptor versions are seen. Note that
// new descriptor versions may not monotonically increase. For example, the sequence schema
// may be altered in a txn, so the cache sees a new version V and invalidates/repopulates itself. Then,
// the txn may get rolled back, so the cache will see version V-1 and invalidate/repopulate itself again.
type SequenceCache map[uint32]*SequenceCacheEntry

// NextValue fetches the next value in the sequence cache. If the values in the cache have all been
// given out or if the descriptor version has changed, then fetchNextValues() is used to repopulate the cache.
func (sc SequenceCache) NextValue(
	seqID uint32, clientVersion uint32, fetchNextValues func() (int64, int64, int64, error),
) (int64, error) {
	// Create entry for this sequence ID if there are no existing entries.
	if _, found := sc[seqID]; !found {
		sc[seqID] = &SequenceCacheEntry{}
	}
	cacheEntry := sc[seqID]

	if cacheEntry.NumValues > 0 && cacheEntry.CachedVersion == clientVersion {
		cacheEntry.CurrentValue += cacheEntry.Increment
		cacheEntry.NumValues--
		return cacheEntry.CurrentValue - cacheEntry.Increment, nil
	}

	currentValue, increment, numValues, err := fetchNextValues()
	if err != nil {
		return 0, err
	}

	// One value must be returned, and the rest of the values are stored.
	val := currentValue
	cacheEntry.CurrentValue = currentValue + increment
	cacheEntry.Increment = increment
	cacheEntry.NumValues = numValues - 1
	cacheEntry.CachedVersion = clientVersion
	return val, nil
}
