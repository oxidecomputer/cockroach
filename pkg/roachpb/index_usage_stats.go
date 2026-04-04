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

package roachpb

// TableID is same as descpb.ID. We redefine it here to avoid importing descpb.
type TableID uint32

// IndexID is same as descpb.IndexID. We redefine it here to avoid importing
// descpb.
type IndexID uint32

// Add adds the fields from other IndexUsageStatistics.
func (m *IndexUsageStatistics) Add(other *IndexUsageStatistics) {
	m.TotalRowsRead += other.TotalRowsRead
	m.TotalRowsWritten += other.TotalRowsWritten

	m.TotalReadCount += other.TotalReadCount
	m.TotalWriteCount += other.TotalWriteCount

	if m.LastWrite.Before(other.LastWrite) {
		m.LastWrite = other.LastWrite
	}

	if m.LastRead.Before(other.LastRead) {
		m.LastRead = other.LastRead
	}
}
