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

package slstorage

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/sqlliveness"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
	"github.com/cockroachdb/errors"
)

// FakeStorage implements the sqlliveness.Storage interface.
type FakeStorage struct {
	mu struct {
		syncutil.Mutex
		sessions map[sqlliveness.SessionID]hlc.Timestamp
	}
}

// NewFakeStorage creates a new FakeStorage.
func NewFakeStorage() *FakeStorage {
	fs := &FakeStorage{}
	fs.mu.sessions = make(map[sqlliveness.SessionID]hlc.Timestamp)
	return fs
}

// IsAlive implements the sqlliveness.Storage interface.
func (s *FakeStorage) IsAlive(
	_ context.Context, sid sqlliveness.SessionID,
) (alive bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.mu.sessions[sid]
	return ok, nil
}

// Insert implements the sqlliveness.Storage interface.
func (s *FakeStorage) Insert(
	_ context.Context, sid sqlliveness.SessionID, expiration hlc.Timestamp,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.mu.sessions[sid]; ok {
		return errors.Errorf("session %s already exists", sid)
	}
	s.mu.sessions[sid] = expiration
	return nil
}

// Update implements the sqlliveness.Storage interface.
func (s *FakeStorage) Update(
	_ context.Context, sid sqlliveness.SessionID, expiration hlc.Timestamp,
) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.mu.sessions[sid]; !ok {
		return false, nil
	}
	s.mu.sessions[sid] = expiration
	return true, nil
}

// Delete is needed to manually delete a session for testing purposes.
func (s *FakeStorage) Delete(_ context.Context, sid sqlliveness.SessionID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.mu.sessions, sid)
	return nil
}
