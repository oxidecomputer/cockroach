// Copyright 2017 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

// cockroach-short is an entry point for a CockroachDB binary that excludes
// certain components that are slow to build or have heavyweight dependencies.
// At present, the only excluded component is the web UI.
package main

import "github.com/cockroachdb/cockroach/pkg/cli"

func main() {
	cli.Main()
}
