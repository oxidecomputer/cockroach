// Copyright 2017 The Cockroach Authors.
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

package sql

import (
	"context"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/server/telemetry"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/roleoption"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqltelemetry"
	"github.com/cockroachdb/cockroach/pkg/util/errorutil/unimplemented"
	"github.com/cockroachdb/cockroach/pkg/util/log/eventpb"
	"github.com/cockroachdb/errors"
)

type createDatabaseNode struct {
	n *tree.CreateDatabase
}

// CreateDatabase creates a database.
// Privileges: superuser or CREATEDB
func (p *planner) CreateDatabase(ctx context.Context, n *tree.CreateDatabase) (planNode, error) {
	if err := checkSchemaChangeEnabled(
		ctx,
		p.ExecCfg(),
		"CREATE DATABASE",
	); err != nil {
		return nil, err
	}

	if n.Name == "" {
		return nil, errEmptyDatabaseName
	}

	if tmpl := n.Template; tmpl != "" {
		// See https://www.postgresql.org/docs/current/static/manage-ag-templatedbs.html
		if !strings.EqualFold(tmpl, "template0") {
			return nil, unimplemented.NewWithIssuef(10151,
				"unsupported template: %s", tmpl)
		}
	}

	if enc := n.Encoding; enc != "" {
		// We only support UTF8 (and aliases for UTF8).
		if !(strings.EqualFold(enc, "UTF8") ||
			strings.EqualFold(enc, "UTF-8") ||
			strings.EqualFold(enc, "UNICODE")) {
			return nil, unimplemented.NewWithIssueDetailf(35882, "create.db.encoding",
				"unsupported encoding: %s", enc)
		}
	}

	if col := n.Collate; col != "" {
		// We only support C and C.UTF-8.
		if col != "C" && col != "C.UTF-8" {
			return nil, unimplemented.NewWithIssueDetailf(16618, "create.db.collation",
				"unsupported collation: %s", col)
		}
	}

	if ctype := n.CType; ctype != "" {
		// We only support C and C.UTF-8.
		if ctype != "C" && ctype != "C.UTF-8" {
			return nil, unimplemented.NewWithIssueDetailf(35882, "create.db.classification",
				"unsupported character classification: %s", ctype)
		}
	}

	if n.ConnectionLimit != -1 {
		return nil, unimplemented.NewWithIssueDetailf(
			54241,
			"create.db.connection_limit",
			"only connection limit -1 is supported, got: %d",
			n.ConnectionLimit,
		)
	}

	if n.SurvivalGoal != tree.SurvivalGoalDefault &&
		n.PrimaryRegion == tree.PrimaryRegionNotSpecifiedName {
		return nil, pgerror.New(
			pgcode.InvalidDatabaseDefinition,
			"PRIMARY REGION must be specified when using SURVIVE",
		)
	}

	if n.Placement != tree.DataPlacementUnspecified {
		if !p.EvalContext().SessionData().PlacementEnabled {
			return nil, errors.WithHint(pgerror.New(
				pgcode.ExperimentalFeature,
				"PLACEMENT requires that the session setting enable_multiregion_placement_policy "+
					"is enabled",
			),
				"to use PLACEMENT, enable the session setting with SET"+
					" enable_multiregion_placement_policy = true or enable the cluster setting"+
					" sql.defaults.multiregion_placement_policy.enabled",
			)
		}

		if n.PrimaryRegion == tree.PrimaryRegionNotSpecifiedName {
			return nil, pgerror.New(
				pgcode.InvalidDatabaseDefinition,
				"PRIMARY REGION must be specified when using PLACEMENT",
			)
		}

		if n.Placement == tree.DataPlacementRestricted &&
			n.SurvivalGoal == tree.SurvivalGoalRegionFailure {
			return nil, pgerror.New(
				pgcode.InvalidDatabaseDefinition,
				"PLACEMENT RESTRICTED can only be used with SURVIVE ZONE FAILURE",
			)
		}
	}

	hasCreateDB, err := p.HasRoleOption(ctx, roleoption.CREATEDB)
	if err != nil {
		return nil, err
	}
	if !hasCreateDB {
		return nil, pgerror.New(
			pgcode.InsufficientPrivilege,
			"permission denied to create database",
		)
	}

	return &createDatabaseNode{n: n}, nil
}

func (n *createDatabaseNode) startExec(params runParams) error {
	telemetry.Inc(sqltelemetry.SchemaChangeCreateCounter("database"))

	desc, created, err := params.p.createDatabase(
		params.ctx, n.n, tree.AsStringWithFQNames(n.n, params.Ann()))
	if err != nil {
		return err
	}
	if created {
		// Log Create Database event. This is an auditable log event and is
		// recorded in the same transaction as the table descriptor update.
		if err := params.p.logEvent(params.ctx, desc.GetID(),
			&eventpb.CreateDatabase{
				DatabaseName: n.n.Name.String(),
			}); err != nil {
			return err
		}
	}
	return nil
}

func (*createDatabaseNode) Next(runParams) (bool, error) { return false, nil }
func (*createDatabaseNode) Values() tree.Datums          { return tree.Datums{} }
func (*createDatabaseNode) Close(context.Context)        {}

// ReadingOwnWrites implements the planNodeReadingOwnWrites Interface. This is
// required because we create a type descriptor for multi-region databases,
// which must be read during validation. We also call CONFIGURE ZONE which
// perms multiple KV operations on descriptors and expects to see its own writes.
func (*createDatabaseNode) ReadingOwnWrites() {}
