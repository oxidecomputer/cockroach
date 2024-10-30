#!/bin/bash
#:
#: name = "artifacts-macos"
#: variety = "basic"
#: target = "ubuntu-22.04"
#: output_rules = [
#:	"=/work/oxidecomputer/cockroach/cockroach.tgz",
#:	"=/work/oxidecomputer/cockroach/cockroach.tgz.sha256",
#: ]
#:
#: [[publish]]
#: series = "darwin-amd64"
#: name = "cockroach.tgz"
#: from_output = "/work/oxidecomputer/cockroach/cockroach.tgz"
#:
#: [[publish]]
#: series = "darwin-amd64"
#: name = "cockroach.tgz.sha256"
#: from_output = "/work/oxidecomputer/cockroach/cockroach.tgz.sha256"
#:
#: [[publish]]
#: series = "darwin-arm64"
#: name = "cockroach.tgz"
#: from_output = "/work/oxidecomputer/cockroach/cockroach.tgz"
#:
#: [[publish]]
#: series = "darwin-arm64"
#: name = "cockroach.tgz.sha256"
#: from_output = "/work/oxidecomputer/cockroach/cockroach.tgz.sha256"

set -o errexit
set -o pipefail
set -o xtrace

sudo apt-get install -y jq unzip
timeout 30m .github/buildomat/fetch-gh-artifacts.sh build-macos
unzip build.zip
