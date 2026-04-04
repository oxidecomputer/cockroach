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
	"os"
	"path/filepath"
)

// whereis is a helper executable that is basically just `realpath`. It's meant
// to be used like:
//     bazel run ... --run_under //build/bazelutil/whereis
// ... which will print the location of the binary you're running. Useful
// because Bazel can be a little unclear about where exactly to find any given
// executable.
func main() {
	if len(os.Args) != 2 {
		panic("expected a single argument")
	}
	abs, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", abs)
}
