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

package config

import (
	"sort"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
)

// SystemConfigMask is a mask that can be applied to a set of system config
// entries to filter out unwanted entries.
type SystemConfigMask struct {
	allowed []roachpb.Key
}

// MakeSystemConfigMask constructs a new SystemConfigMask that passes through
// only the specified keys when applied.
func MakeSystemConfigMask(allowed ...roachpb.Key) SystemConfigMask {
	sort.Slice(allowed, func(i, j int) bool {
		return allowed[i].Compare(allowed[j]) < 0
	})
	return SystemConfigMask{allowed: allowed}
}

// Apply applies the mask to the provided set of system config entries,
// returning a filtered down set of entries.
func (m SystemConfigMask) Apply(entries SystemConfigEntries) SystemConfigEntries {
	var res SystemConfigEntries
	for _, key := range m.allowed {
		i := sort.Search(len(entries.Values), func(i int) bool {
			return entries.Values[i].Key.Compare(key) >= 0
		})
		if i < len(entries.Values) && entries.Values[i].Key.Equal(key) {
			res.Values = append(res.Values, entries.Values[i])
		}
	}
	return res
}
