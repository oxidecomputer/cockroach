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
	"context"
	gosql "database/sql"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
)

var reg = NewRegistry(1 /* numNodes */, MakeClusterConstructor(func(
	t testing.TB, knobs base.TestingKnobs,
) (_ *gosql.DB, cleanup func()) {
	s, sql, _ := serverutils.StartServer(t, base.TestServerArgs{
		UseDatabase: "bench",
		Knobs:       knobs,
	})
	return sql, func() { s.Stopper().Stop(context.Background()) }
}))
