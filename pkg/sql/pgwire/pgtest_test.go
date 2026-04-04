// Copyright 2019 The Cockroach Authors.
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

package pgwire

import (
	"context"
	"flag"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	_ "github.com/cockroachdb/cockroach/pkg/cloud/impl" // register cloud storage providers
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/pgtest"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

var (
	flagAddr = flag.String("addr", "", "pass a custom postgres address to TestWalk instead of starting an in-memory node")
	flagUser = flag.String("user", "postgres", "username used if -addr is specified")
)

func TestPGTest(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	if *flagAddr == "" {
		newServer := func() (addr, user string, cleanup func()) {
			ctx := context.Background()
			s, db, _ := serverutils.StartServer(t, base.TestServerArgs{
				Insecure: true,
			})
			cleanup = func() {
				s.Stopper().Stop(ctx)
			}
			addr = s.ServingSQLAddr()
			user = security.RootUser
			// None of the tests read that much data, so we hardcode the max message
			// size to something small. This lets us test the handling of large
			// query inputs. See the large_input test.
			_, _ = db.ExecContext(ctx, "SET CLUSTER SETTING sql.conn.max_read_buffer_message_size = '32 KiB'")
			return addr, user, cleanup
		}
		pgtest.WalkWithNewServer(t, testutils.TestDataPath(t, "pgtest"), newServer)
	} else {
		pgtest.WalkWithRunningServer(t, testutils.TestDataPath(t, "pgtest"), *flagAddr, *flagUser)
	}
}
