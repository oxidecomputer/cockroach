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
//
// bazci is glue code to make debugging Bazel builds and tests in Teamcity as
// painless as possible.
//
// bazci [build|test] \
//     --artifacts_dir=$ARTIFACTS_DIR targets... -- [command-line options]
//
// bazci will invoke a `bazel build` or `bazel test` of all the given targets
// and stage the resultant build/test artifacts in the given `artifacts_dir`.
// The build/test artifacts are munged slightly such that TC can easily parse
// them.
package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("")

	if _, err := exec.LookPath("bazel"); err != nil {
		log.Printf("ERROR: bazel not found in $PATH")
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}
