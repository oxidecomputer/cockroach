#!/bin/bash
#:
#: name = "build-linux"
#: variety = "basic"
#: target = "ubuntu-22.04"
#: output_rules = [
#:	"=/work/oxidecomputer/cockroach/cockroach.tgz",
#:	"=/work/oxidecomputer/cockroach/cockroach.tgz.sha256",
#: ]
#:
#: [[publish]]
#: series = "linux-amd64"
#: name = "cockroach.tgz"
#: from_output = "=/work/oxidecomputer/cockroach/cockroach.tgz"
#:
#: [[publish]]
#: series = "linux-amd64"
#: name = "cockroach.tgz.sha256"
#: from_output = "=/work/oxidecomputer/cockroach/cockroach.tgz.sha256"

set -o errexit
set -o pipefail
set -o xtrace

export BROWSERSLIST_IGNORE_OLD_DATA=1

source .github/buildomat/linux-setup.sh
gmake -j"$(nproc)" cockroach.tgz BUILDTYPE=release
