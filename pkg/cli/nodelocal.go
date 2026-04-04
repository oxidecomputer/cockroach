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

package cli

import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/cockroachdb/cockroach/pkg/cli/clierrorplus"
	"github.com/cockroachdb/cockroach/pkg/cli/clisqlclient"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

const (
	chunkSize = 4 * 1024
)

var nodeLocalUploadCmd = &cobra.Command{
	Use:   "upload <source> <destination>",
	Short: "Upload file from source to destination",
	Long: `
Uploads a file to a gateway node's local file system using a SQL connection.
`,
	Args: cobra.MinimumNArgs(2),
	RunE: clierrorplus.MaybeShoutError(runUpload),
}

func runUpload(cmd *cobra.Command, args []string) (resErr error) {
	conn, err := makeSQLClient("cockroach nodelocal", useSystemDb)
	if err != nil {
		return err
	}
	defer func() { resErr = errors.CombineErrors(resErr, conn.Close()) }()

	source := args[0]
	destination := args[1]
	reader, err := openSourceFile(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	return uploadFile(context.Background(), conn, reader, destination)
}

func openSourceFile(source string) (io.ReadCloser, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get source file stats for %s", source)
	}
	if stat.IsDir() {
		return nil, fmt.Errorf("source file %s is a directory, not a file", source)
	}
	return f, nil
}

func uploadFile(
	ctx context.Context, conn clisqlclient.Conn, reader io.Reader, destination string,
) error {
	if err := conn.EnsureConn(); err != nil {
		return err
	}

	ex := conn.GetDriverConn()
	if _, err := ex.ExecContext(ctx, `BEGIN`, nil); err != nil {
		return err
	}

	// Construct the nodelocal URI as the destination for the CopyIn stmt.
	nodelocalURL := url.URL{
		Scheme: "nodelocal",
		Host:   "self",
		Path:   destination,
	}
	stmt, err := conn.GetDriverConn().Prepare(sql.CopyInFileStmt(nodelocalURL.String(), sql.CrdbInternalName,
		sql.NodelocalFileUploadTable))
	if err != nil {
		return err
	}

	defer func() {
		if stmt != nil {
			_ = stmt.Close()
			_, _ = ex.ExecContext(ctx, `ROLLBACK`, nil)
		}
	}()

	send := make([]byte, chunkSize)
	for {
		n, err := reader.Read(send)
		if n > 0 {
			// TODO(adityamaru): Switch to StmtExecContext once the copyin driver
			// supports it.
			//lint:ignore SA1019 DriverConn doesn't support go 1.8 API
			_, err = stmt.Exec([]driver.Value{string(send[:n])})
			if err != nil {
				return err
			}
		} else if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	if err := stmt.Close(); err != nil {
		return err
	}
	stmt = nil

	if _, err := ex.ExecContext(ctx, `COMMIT`, nil); err != nil {
		return err
	}

	nodeID, _, _, err := conn.GetServerMetadata(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get node id")
	}
	fmt.Printf("successfully uploaded to nodelocal://%s\n", filepath.Join(roachpb.NodeID(nodeID).String(), destination))
	return nil
}

var nodeLocalCmds = []*cobra.Command{
	nodeLocalUploadCmd,
}

var nodeLocalCmd = &cobra.Command{
	Use:   "nodelocal [command]",
	Short: "upload and delete nodelocal files",
	Long:  "Upload and delete files on the gateway node's local file system.",
	RunE:  UsageAndErr,
}

func init() {
	nodeLocalCmd.AddCommand(nodeLocalCmds...)
}
