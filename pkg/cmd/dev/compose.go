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

package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func makeComposeCmd(runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	composeCmd := &cobra.Command{
		Use:     "compose",
		Short:   "Run compose tests",
		Long:    "Run compose tests.",
		Example: "dev compose",
		Args:    cobra.ExactArgs(0),
		RunE:    runE,
	}
	addCommonBuildFlags(composeCmd)
	addCommonTestFlags(composeCmd)
	composeCmd.Flags().String(volumeFlag, "bzlhome", "the Docker volume to use as the container home directory (only used for cross builds)")
	return composeCmd
}

func (d *dev) compose(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	var (
		filter  = mustGetFlagString(cmd, filterFlag)
		short   = mustGetFlagBool(cmd, shortFlag)
		timeout = mustGetFlagDuration(cmd, timeoutFlag)
	)

	crossArgs, targets, err := d.getBasicBuildArgs(ctx, []string{"//pkg/cmd/cockroach:cockroach", "//pkg/compose/compare/compare:compare_test"})
	if err != nil {
		return err
	}
	volume := mustGetFlagString(cmd, volumeFlag)
	err = d.crossBuild(ctx, crossArgs, targets, "crosslinux", volume)
	if err != nil {
		return err
	}

	workspace, err := d.getWorkspace(ctx)
	if err != nil {
		return err
	}
	cockroachBin := filepath.Join(workspace, "artifacts", "cockroach")
	compareBin := filepath.Join(workspace, "artifacts", "compare_test")

	var args []string
	args = append(args, "run", "//pkg/compose:compose_test", "--config=test")
	if numCPUs != 0 {
		args = append(args, fmt.Sprintf("--local_cpu_resources=%d", numCPUs))
	}
	if filter != "" {
		args = append(args, fmt.Sprintf("--test_filter=%s", filter))
	}
	if short {
		args = append(args, "--test_arg", "-test.short")
	}
	if timeout > 0 {
		args = append(args, fmt.Sprintf("--test_timeout=%d", int(timeout.Seconds())))
	}

	args = append(args, "--test_arg", "-cockroach", "--test_arg", cockroachBin)
	args = append(args, "--test_arg", "-compare", "--test_arg", compareBin)

	logCommand("bazel", args...)
	return d.exec.CommandContextInheritingStdStreams(ctx, "bazel", args...)
}
