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

package tests

import (
	"context"
	"strconv"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/roachprod/install"
	"github.com/cockroachdb/cockroach/pkg/roachprod/logger"
)

// getDiskUsageInBytes does what's on the tin. nodeIdx starts at one.
func getDiskUsageInBytes(
	ctx context.Context, c cluster.Cluster, logger *logger.Logger, nodeIdx int,
) (int, error) {
	var result install.RunResultDetails
	for {
		var err error
		// `du` can warn if files get removed out from under it (which
		// happens during RocksDB compactions, for example). Discard its
		// stderr to avoid breaking Atoi later.
		// TODO(bdarnell): Refactor this stack to not combine stdout and
		// stderr so we don't need to do this (and the Warning check
		// below).
		result, err = c.RunWithDetailsSingleNode(
			ctx,
			logger,
			c.Node(nodeIdx),
			"du -sk {store-dir} 2>/dev/null | grep -oE '^[0-9]+'",
		)
		if err != nil {
			if ctx.Err() != nil {
				return 0, ctx.Err()
			}
			// If `du` fails, retry.
			// TODO(bdarnell): is this worth doing? It was originally added
			// because of the "files removed out from under it" problem, but
			// that doesn't result in a command failure, just a stderr
			// message.
			logger.Printf("retrying disk usage computation after spurious error: %s", err)
			continue
		}

		break
	}

	// We need this check because sometimes the first line of the roachprod output is a warning
	// about adding an ip to a list of known hosts.
	if strings.Contains(result.Stdout, "Warning") {
		result.Stdout = strings.Split(result.Stdout, "\n")[1]
	}

	size, err := strconv.Atoi(strings.TrimSpace(result.Stdout))
	if err != nil {
		return 0, err
	}

	return size * 1024, nil
}
