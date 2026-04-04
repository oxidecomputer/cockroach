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

// Check that GitHub PR descriptions and commit messages contain the
// expected epic and issue references.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "docs-issue-generation",
	Short: "Generate a new set of release issues in the docs repo for a given commit.",
	Run: func(_ *cobra.Command, args []string) {
		params := defaultEnvParameters()
		docsIssueGeneration(params)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func defaultEnvParameters() parameters {
	const (
		githubAPITokenEnv = "GITHUB_API_TOKEN"
		buildVcsNumberEnv = "BUILD_VCS_NUMBER"
	)

	return parameters{
		Token: maybeEnv(githubAPITokenEnv, ""),
		Sha:   maybeEnv(buildVcsNumberEnv, "4dd8da9609adb3acce6795cea93b67ccacfc0270"),
	}
}

func maybeEnv(envKey, defaultValue string) string {
	v := os.Getenv(envKey)
	if v == "" {
		return defaultValue
	}
	return v
}
