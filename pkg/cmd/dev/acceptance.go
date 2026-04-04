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

func makeAcceptanceCmd(runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	acceptanceCmd := &cobra.Command{
		Use:     "acceptance",
		Short:   "Run acceptance tests",
		Long:    "Run acceptance tests.",
		Example: "dev acceptance",
		Args:    cobra.ExactArgs(0),
		RunE:    runE,
	}
	addCommonBuildFlags(acceptanceCmd)
	addCommonTestFlags(acceptanceCmd)
	acceptanceCmd.Flags().String(volumeFlag, "bzlhome", "the Docker volume to use as the container home directory (only used for cross builds)")
	return acceptanceCmd
}

func (d *dev) acceptance(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	var (
		filter  = mustGetFlagString(cmd, filterFlag)
		short   = mustGetFlagBool(cmd, shortFlag)
		timeout = mustGetFlagDuration(cmd, timeoutFlag)
	)

	// First we have to build cockroach.
	crossArgs, targets, err := d.getBasicBuildArgs(ctx, []string{"//pkg/cmd/cockroach-short:cockroach-short"})
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
	logDir := filepath.Join(workspace, "artifacts", "acceptance")
	err = d.os.MkdirAll(logDir)
	if err != nil {
		return err
	}
	cockroachBin := filepath.Join(workspace, "artifacts", "cockroach-short")

	var args []string
	args = append(args, "run", "//pkg/acceptance:acceptance_test", "--config=test")
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
	args = append(args, fmt.Sprintf("--test_arg=-l=%s", logDir))
	args = append(args, "--test_env=TZ=America/New_York")
	args = append(args, fmt.Sprintf("--test_arg=-b=%s", cockroachBin))

	logCommand("bazel", args...)
	return d.exec.CommandContextInheritingStdStreams(ctx, "bazel", args...)
}
