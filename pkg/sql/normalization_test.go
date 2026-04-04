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

package sql_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

func TestNFCNormalization(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	s, db, _ := serverutils.StartServer(t, base.TestServerArgs{Insecure: true})
	defer s.Stopper().Stop(context.Background())
	defer db.Close()

	_, err := db.Exec("CREATE TABLE café (a INT)")
	require.NoError(t, err)

	_, err = db.Exec("CREATE TABLE Cafe\u0301 (a INT)")
	require.Errorf(t, err, "The tables should be considered duplicates when normalized")
	require.True(t, strings.Contains(err.Error(), "already exists"))

	_, err = db.Exec("CREATE TABLE cafe\u0301 (a INT)")
	require.Errorf(t, err, "The tables should be considered duplicates when normalized")
	require.True(t, strings.Contains(err.Error(), "already exists"))

	_, err = db.Exec("CREATE TABLE caf\u00E9 (a INT)")
	require.Errorf(t, err, "The tables should be considered duplicates when normalized")
	require.True(t, strings.Contains(err.Error(), "already exists"))

	_, err = db.Exec("CREATE TABLE \"caf\u00E9\" (a INT)")
	require.Errorf(t, err, "The tables should be considered duplicates when normalized")
	require.True(t, strings.Contains(err.Error(), "already exists"))

	_, err = db.Exec("CREATE TABLE \"cafe\u0301\" (a INT)")
	require.Errorf(t, err, "The tables should be considered duplicates when normalized")
	require.True(t, strings.Contains(err.Error(), "already exists"))

	_, err = db.Exec(`CREATE TABLE "Café" (a INT)`)
	require.NoError(t, err)
	//Ensure normal strings are not normalized like double quoted strings
	var b bool
	err = db.QueryRow("SELECT 'caf\u00E9' = 'cafe\u0301'").Scan(&b)
	require.NoError(t, err)
	require.False(t, b)

}
