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

package rttanalysis

import (
	gosql "database/sql"
	"sync"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/util/tracing"
)

// ClusterConstructor is used to construct a Cluster for an individual case run.
type ClusterConstructor func(testing.TB) *Cluster

// MakeClusterConstructor creates a new ClusterConstructor using the provided
// function. The intention is that the caller will use the provided knobs when
// constructing the cluster and will return a handle to the SQL database.
func MakeClusterConstructor(
	f func(testing.TB, base.TestingKnobs) (_ *gosql.DB, cleanup func()),
) ClusterConstructor {
	return func(t testing.TB) *Cluster {
		c := &Cluster{}
		beforePlan := func(trace tracing.Recording, stmt string) {
			if _, ok := c.stmtToKVBatchRequests.Load(stmt); ok {
				c.stmtToKVBatchRequests.Store(stmt, trace)
			}
		}
		c.sql, c.cleanup = f(t, base.TestingKnobs{
			SQLExecutor: &sql.ExecutorTestingKnobs{
				WithStatementTrace: beforePlan,
			},
		})
		return c
	}
}

// Cluster abstracts a cockroach cluster for use in rttanalysis benchmarks.
type Cluster struct {
	stmtToKVBatchRequests sync.Map
	cleanup               func()
	sql                   *gosql.DB
}

func (c *Cluster) conn() *gosql.DB {
	return c.sql
}

func (c *Cluster) clearStatementTrace(stmt string) {
	c.stmtToKVBatchRequests.Store(stmt, nil)
}

func (c *Cluster) getStatementTrace(stmt string) (tracing.Recording, bool) {
	out, _ := c.stmtToKVBatchRequests.Load(stmt)
	r, ok := out.(tracing.Recording)
	return r, ok
}

func (c *Cluster) close() {
	c.cleanup()
}
