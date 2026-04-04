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

package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{Use: "release"}

const (
	envSMTPUser     = "SMTP_USER"
	envSMTPPassword = "SMTP_PASSWORD"
	envGithubToken  = "GITHUB_TOKEN"
	releaseSeries   = "release-series"
	templatesDir    = "template-dir"
	smtpUser        = "smtp-user"
	smtpHost        = "smtp-host"
	smtpPort        = "smtp-port"
	emailAddresses  = "to"
	dryRun          = "dry-run"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(postReleaseSeriesBlockersCmd)
	rootCmd.AddCommand(setOrchestrationVersionCmd)
}
