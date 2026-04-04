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

package cli

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver"
	"github.com/cockroachdb/cockroach/pkg/server"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/errors"
)

// runInitialSQL concerns itself with running "initial SQL" code when
// a cluster is started for the first time.
//
// The "startSingleNode" argument is true for `start-single-node`,
// and `cockroach demo` with 2 nodes or fewer.
// If adminUser is non-empty, an admin user with that name is
// created upon initialization. Its password is then also returned.
func runInitialSQL(
	ctx context.Context, s *server.Server, startSingleNode bool, adminUser, adminPassword string,
) error {
	newCluster := s.InitialStart() && s.NodeID() == kvserver.FirstNodeID
	if !newCluster {
		// The initial SQL code only runs the first time the cluster is initialized.
		return nil
	}

	if startSingleNode {
		// For start-single-node, set the default replication factor to
		// 1 so as to avoid warning messages and unnecessary rebalance
		// churn.
		if err := cliDisableReplication(ctx, s); err != nil {
			log.Ops.Errorf(ctx, "could not disable replication: %v", err)
			return err
		}
		log.Ops.Infof(ctx, "Replication was disabled for this cluster.\n"+
			"When/if adding nodes in the future, update zone configurations to increase the replication factor.")
	}

	if adminUser != "" && !s.Insecure() {
		if err := createAdminUser(ctx, s, adminUser, adminPassword); err != nil {
			return err
		}
	}

	return nil
}

// createAdminUser creates an admin user with the given name.
func createAdminUser(ctx context.Context, s *server.Server, adminUser, adminPassword string) error {
	return s.RunLocalSQL(ctx,
		func(ctx context.Context, ie *sql.InternalExecutor) error {
			_, err := ie.Exec(
				ctx, "admin-user", nil,
				fmt.Sprintf("CREATE USER %s WITH PASSWORD $1", adminUser),
				adminPassword,
			)
			if err != nil {
				return err
			}
			// TODO(knz): Demote the admin user to an operator privilege with fewer options.
			_, err = ie.Exec(ctx, "admin-user", nil, fmt.Sprintf("GRANT admin TO %s", tree.Name(adminUser)))
			return err
		})
}

// cliDisableReplication changes the replication factor on
// all defined zones to become 1. This is used by start-single-node
// and demo to define single-node clusters, so as to avoid
// churn in the log files.
//
// The change is effected using the internal SQL interface of the
// given server object.
func cliDisableReplication(ctx context.Context, s *server.Server) error {
	return s.RunLocalSQL(ctx,
		func(ctx context.Context, ie *sql.InternalExecutor) (retErr error) {
			it, err := ie.QueryIterator(ctx, "get-zones", nil,
				"SELECT target FROM crdb_internal.zones")
			if err != nil {
				return err
			}
			// We have to make sure to close the iterator since we might return
			// from the for loop early (before Next() returns false).
			defer func() { retErr = errors.CombineErrors(retErr, it.Close()) }()

			var ok bool
			for ok, err = it.Next(ctx); ok; ok, err = it.Next(ctx) {
				zone := string(*it.Cur()[0].(*tree.DString))
				if _, err := ie.Exec(ctx, "set-zone", nil,
					fmt.Sprintf("ALTER %s CONFIGURE ZONE USING num_replicas = 1", zone)); err != nil {
					return err
				}
			}
			return err
		})
}
