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
	"fmt"
	"os"

	"github.com/cockroachdb/cockroach/pkg/cli/clierrorplus"
	"github.com/cockroachdb/cockroach/pkg/cli/clisqlshell"
	"github.com/cockroachdb/cockroach/pkg/server/pgurl"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catconstants"
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

// sqlShellCmd opens a sql shell.
var sqlShellCmd = &cobra.Command{
	Use:   "sql [options]",
	Short: "open a sql shell",
	Long: `
Open a sql shell running against a cockroach database.
`,
	Args: cobra.NoArgs,
	RunE: clierrorplus.MaybeDecorateError(runTerm),
}

func runTerm(cmd *cobra.Command, args []string) (resErr error) {
	closeFn, err := sqlCtx.Open(os.Stdin)
	if err != nil {
		return err
	}
	defer closeFn()

	if cliCtx.IsInteractive {
		// The user only gets to see the welcome message on interactive sessions.
		// Refer to README.md to understand the general design guidelines for
		// help texts.
		const welcomeMessage = `#
# Welcome to the CockroachDB SQL shell.
# All statements must be terminated by a semicolon.
# To exit, type: \q.
#
`
		fmt.Print(welcomeMessage)
	}

	conn, err := makeSQLClient(catconstants.InternalSQLAppName, useDefaultDb)
	if err != nil {
		return err
	}
	defer func() { resErr = errors.CombineErrors(resErr, conn.Close()) }()

	sqlCtx.ShellCtx.ParseURL = makeURLParser(cmd)
	return sqlCtx.Run(conn)
}

func makeURLParser(cmd *cobra.Command) clisqlshell.URLParser {
	return func(url string) (*pgurl.URL, error) {
		// Parse it as if --url was specified.
		up := urlParser{cmd: cmd, cliCtx: &cliCtx}
		if err := up.setInternal(url, false /* warn */); err != nil {
			return nil, err
		}
		return cliCtx.sqlConnURL, nil
	}
}
