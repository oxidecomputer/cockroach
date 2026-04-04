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

package clisqlshell

import (
	"os"

	"github.com/cockroachdb/cockroach/pkg/server/pgurl"
)

// Shell represents an interactive shell
type Shell interface {
	// RunInteractive runs the shell.
	RunInteractive(cmdIn, cmdOut, cmdErr *os.File) (exitErr error)
}

// URLParser represents a function able to convert user-supplied
// strings to a URL object.
type URLParser = func(url string) (*pgurl.URL, error)
