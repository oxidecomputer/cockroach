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

package main

import "github.com/spf13/cobra"

func makeGoCmd(runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	return &cobra.Command{
		Use:     "go <arguments>",
		Short:   "Run `go` with the given arguments",
		Long:    "Run `go` with the given arguments",
		Example: "dev go mod tidy",
		Args:    cobra.MinimumNArgs(0),
		RunE:    runE,
	}
}

func (d *dev) gocmd(cmd *cobra.Command, commandLine []string) error {
	beforeDash, afterDash := splitArgsAtDash(cmd, commandLine)
	ctx := cmd.Context()
	args := []string{"run", "@go_sdk//:bin/go", "--ui_event_filters=-DEBUG,-info,-stdout,-stderr", "--noshow_progress", "--"}
	args = append(args, beforeDash...)
	args = append(args, afterDash...)
	return d.exec.CommandContextInheritingStdStreams(ctx, "bazel", args...)
}
