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

	dev := makeDevCmd()

	if err := dev.cli.Execute(); err != nil {
		log.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}
