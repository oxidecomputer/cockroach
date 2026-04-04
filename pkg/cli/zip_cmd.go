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

package cli

import (
	"github.com/cockroachdb/cockroach/pkg/cli/clierrorplus"
	"github.com/spf13/cobra"
)

var debugZipCmd = &cobra.Command{
	Use:   "zip <file>",
	Short: "gather cluster debug data into a zip file",
	Long: `

Gather cluster debug data into a zip file. Data includes cluster events, node
liveness, node status, range status, node stack traces, node engine stats, log
files, and SQL schema.

Retrieval of per-node details (status, stack traces, range status, engine stats)
requires the node to be live and operating properly. Retrieval of SQL data
requires the cluster to be live.
`,
	Args: cobra.ExactArgs(1),
	RunE: clierrorplus.MaybeDecorateError(runDebugZip),
}
