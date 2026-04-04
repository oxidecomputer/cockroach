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

package tests

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/option"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/test"
)

const goPath = `/mnt/data1/go`

// installGolang installs a specific version of Go on all nodes in
// "node".
func installGolang(
	ctx context.Context, t test.Test, c cluster.Cluster, node option.NodeListOption,
) {
	if err := repeatRunE(
		ctx, t, c, node, "update apt-get", `sudo apt-get -qq update`,
	); err != nil {
		t.Fatal(err)
	}

	if err := repeatRunE(
		ctx,
		t,
		c,
		node,
		"install dependencies (go uses C bindings)",
		`sudo apt-get -qq install build-essential`,
	); err != nil {
		t.Fatal(err)
	}

	if err := repeatRunE(
		ctx, t, c, node, "download go", `curl -fsSL https://dl.google.com/go/go1.17.11.linux-amd64.tar.gz > /tmp/go.tgz`,
	); err != nil {
		t.Fatal(err)
	}
	if err := repeatRunE(
		ctx, t, c, node, "verify tarball", `sha256sum -c - <<EOF
d69a4fe2694f795d8e525c72b497ededc209cb7185f4c3b62d7a98dd6227b3fe /tmp/go.tgz
EOF`,
	); err != nil {
		t.Fatal(err)
	}
	if err := repeatRunE(
		ctx, t, c, node, "extract go", `sudo tar -C /usr/local -zxf /tmp/go.tgz && rm /tmp/go.tgz`,
	); err != nil {
		t.Fatal(err)
	}
	if err := repeatRunE(
		ctx, t, c, node, "force symlink go", "sudo ln -sf /usr/local/go/bin/go /usr/bin",
	); err != nil {
		t.Fatal(err)
	}
}
