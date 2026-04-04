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

package sqllivenesstestutils

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/sqlliveness"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
)

type alwaysAliveSession string

// NewAlwaysAliveSession constructs and returns a session that is forever alive
// for testing purposes.
func NewAlwaysAliveSession(name string) sqlliveness.Session {
	return alwaysAliveSession(name)
}

// ID implements the sqlliveness.Session interface.
func (f alwaysAliveSession) ID() sqlliveness.SessionID {
	return sqlliveness.SessionID(f)
}

// Expiration implements the sqlliveness.Session interface.
func (f alwaysAliveSession) Expiration() hlc.Timestamp { return hlc.MaxTimestamp }

// Start implements the sqlliveness.Session interface.
func (f alwaysAliveSession) Start() hlc.Timestamp { return hlc.MinTimestamp }

// RegisterCallbackForSessionExpiry implements the sqlliveness.Session interface.
func (f alwaysAliveSession) RegisterCallbackForSessionExpiry(func(context.Context)) {}
