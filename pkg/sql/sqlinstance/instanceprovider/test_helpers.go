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

package instanceprovider

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/sqlinstance"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlinstance/instancestorage"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlliveness"
	"github.com/cockroachdb/cockroach/pkg/util/stop"
)

// TestInstanceProvider exposes ShutdownSQLInstanceForTest
// method for testing purposes.
type TestInstanceProvider interface {
	sqlinstance.Provider
	ShutdownSQLInstanceForTest(context.Context)
}

// NewTestInstanceProvider initializes a instanceprovider.provider
// for test purposes
func NewTestInstanceProvider(
	stopper *stop.Stopper, session sqlliveness.Instance, addr string,
) TestInstanceProvider {
	storage := instancestorage.NewFakeStorage()
	p := &provider{
		storage:      storage,
		stopper:      stopper,
		session:      session,
		instanceAddr: addr,
		initialized:  make(chan struct{}),
	}
	p.mu.started = true
	return p
}

// ShutdownSQLInstanceForTest explicitly calls shutdownSQLInstance for testing purposes.
func (p *provider) ShutdownSQLInstanceForTest(ctx context.Context) {
	p.shutdownSQLInstance(ctx)
}
