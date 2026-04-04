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

package enginepb

// Empty returns whether a batch is empty.
func (b *RegistryUpdateBatch) Empty() bool {
	return len(b.Updates) == 0
}

// PutEntry adds an update to the batch corresponding to the addition of a new
// file entry to the registry. The entry should not be nil.
func (b *RegistryUpdateBatch) PutEntry(filename string, entry *FileEntry) {
	b.Updates = append(b.Updates, &RegistryUpdate{Filename: filename, Entry: entry})
}

// DeleteEntry adds an update to the batch corresponding to the deletion of a
// file entry from the registry.
func (b *RegistryUpdateBatch) DeleteEntry(filename string) {
	b.Updates = append(b.Updates, &RegistryUpdate{Filename: filename, Entry: nil})
}
