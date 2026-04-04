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

package tree

// Persistence defines the persistence strategy for a given table.
type Persistence int

const (
	// PersistencePermanent indicates a permanent table.
	PersistencePermanent Persistence = iota
	// PersistenceTemporary indicates a temporary table.
	PersistenceTemporary
	// PersistenceUnlogged indicates an unlogged table.
	// Note this state is not persisted on disk and is used at parse time only.
	PersistenceUnlogged
)

// IsTemporary returns whether the Persistence value is Temporary.
func (p Persistence) IsTemporary() bool {
	return p == PersistenceTemporary
}

// IsUnlogged returns whether the Persistence value is Unlogged.
func (p Persistence) IsUnlogged() bool {
	return p == PersistenceUnlogged
}
