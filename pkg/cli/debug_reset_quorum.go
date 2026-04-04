// Copyright 2020 The Cockroach Authors.
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
	"fmt"
	"strconv"

	"github.com/cockroachdb/cockroach/pkg/cli/clierrorplus"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/spf13/cobra"
)

var debugResetQuorumCmd = &cobra.Command{
	Use: "reset-quorum [range ID]",
	Short: "Reset quorum on the given range" +
		" by designating the target node as the sole voter.",
	Long: `
Reset quorum on the given range by designating the current node as 
the sole voter. Any existing data for the range is discarded. 

This command is UNSAFE and should only be used with the supervision 
of Cockroach Labs support. It is a last-resort option to recover a 
specified range after multiple node failures and loss of quorum.

Data on any surviving replicas will not be used to restore quorum. 
Instead, these replicas will be removed irrevocably.
`,
	Args: cobra.ExactArgs(1),
	RunE: clierrorplus.MaybeDecorateError(runDebugResetQuorum),
}

func runDebugResetQuorum(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rangeID, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		return err
	}

	// Set up GRPC Connection for running ResetQuorum.
	cc, _, finish, err := getClientGRPCConn(ctx, serverCfg)
	if err != nil {
		log.Errorf(ctx, "connection to server failed: %v", err)
		return err
	}
	defer finish()

	// Call ResetQuorum to reset quorum for given range on target node.
	_, err = roachpb.NewInternalClient(cc).ResetQuorum(ctx, &roachpb.ResetQuorumRequest{
		RangeID: int32(rangeID),
	})
	if err != nil {
		return err
	}

	fmt.Printf("ok; please verify https://<console>/#/reports/range/%d", rangeID)

	return nil
}
