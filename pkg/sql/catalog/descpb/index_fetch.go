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

package descpb

// IndexFetchSpecVersionInitial is the initial IndexFetchSpec version.
const IndexFetchSpecVersionInitial = 1

// KeyColumns returns the key columns in the index, excluding any key suffix
// columns.
func (s *IndexFetchSpec) KeyColumns() []IndexFetchSpec_KeyColumn {
	return s.KeyAndSuffixColumns[:len(s.KeyAndSuffixColumns)-int(s.NumKeySuffixColumns)]
}

// KeyFullColumns returns the key columns in the index, plus all key suffix
// columns if that index is not a unique index. It parallels
// TableDescriptor.IndexFullColumns.
func (s *IndexFetchSpec) KeyFullColumns() []IndexFetchSpec_KeyColumn {
	if s.IsUniqueIndex {
		// For unique indexes, the suffix columns are not part of the key (except
		// when the key columns contain a NULL).
		return s.KeyColumns()
	}
	return s.KeyAndSuffixColumns
}

// KeySuffixColumns returns the key suffix columns.
func (s *IndexFetchSpec) KeySuffixColumns() []IndexFetchSpec_KeyColumn {
	return s.KeyAndSuffixColumns[len(s.KeyAndSuffixColumns)-int(s.NumKeySuffixColumns):]
}
