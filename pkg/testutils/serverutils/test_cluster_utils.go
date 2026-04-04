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

package serverutils

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cockroachdb/cockroach/pkg/testutils"
)

// SetClusterSetting executes set cluster settings statement, and then ensures that
// all nodes in the test cluster see that setting update.
func SetClusterSetting(t testutils.TB, c TestClusterInterface, name string, value interface{}) {
	t.Helper()
	strVal := func() string {
		switch v := value.(type) {
		case string:
			return v
		case int, int32, int64:
			return fmt.Sprintf("%d", v)
		case bool:
			return strconv.FormatBool(v)
		case float32, float64:
			return fmt.Sprintf("%f", v)
		case fmt.Stringer:
			return v.String()
		default:
			return fmt.Sprintf("%v", value)
		}
	}()
	query := fmt.Sprintf("SET CLUSTER SETTING %s='%s'", name, strVal)
	// Set cluster setting statement ensures the setting is propagated to the local registry.
	// So, just execute the query against each node in the cluster.
	for i := 0; i < c.NumServers(); i++ {
		_, err := c.ServerConn(i).ExecContext(context.Background(), query)
		if err != nil {
			t.Fatal(err)
		}
	}
}
