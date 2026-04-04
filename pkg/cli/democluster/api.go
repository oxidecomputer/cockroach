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

package democluster

import (
	"context"
	gosql "database/sql"

	democlusterapi "github.com/cockroachdb/cockroach/pkg/cli/democluster/api"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/server"
)

// DemoCluster represents a demo cluster.
type DemoCluster interface {
	democlusterapi.DemoCluster

	// Start starts up the demo cluster.
	// The runInitialSQL function argument is applied to the first server
	// before the initialization completes.
	Start(
		ctx context.Context,
		runInitialSQL func(ctx context.Context, s *server.Server, startSingleNode bool, adminUser, adminPassword string) error,
	) error

	// GetConnURL retrieves the connection URL to the first node.
	GetConnURL() string

	// GetSQLCredentials retrieves the authentication credentials to
	// establish SQL connections to the demo cluster.
	// (These are already embedded in the connection URL produced
	// by GetConnURL() however a client may wish to have them
	// available as discrete values.)
	GetSQLCredentials() (adminUser security.SQLUsername, adminPassword, certsDir string)

	// Close shuts down the demo cluster.
	Close(ctx context.Context)

	// EnableEnterprise enables enterprise features for this demo,
	// if available in this build. The returned callback should be called
	// before terminating the demo.
	EnableEnterprise(ctx context.Context) (func(), error)

	// SetupWorkload initializes the workload generator if defined.
	SetupWorkload(ctx context.Context) error

	// SetClusterSetting overrides a default cluster setting at system level
	// and for all tenants.
	SetClusterSetting(ctx context.Context, setting string, value interface{}) error
}

// EnableEnterprise is not implemented here in order to keep OSS builds successful.
// The cliccl package sets this function if enterprise features are available to demo.
var EnableEnterprise func(db *gosql.DB, org string) (func(), error)
