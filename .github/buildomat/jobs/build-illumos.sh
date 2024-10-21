#!/bin/bash
#:
#: name = "build-illumos"
#: variety = "basic"
#: target = "helios-2.0"
#: output_rules = [
#:	"=/work/oxidecomputer/cockroach/cockroach.tgz",
#:	"=/work/oxidecomputer/cockroach/cockroach.tgz.sha256",
#: ]
#:
#: [[publish]]
#: series = "illumos-amd64"
#: name = "cockroach.tgz"
#: from_output = "/work/oxidecomputer/cockroach/cockroach.tgz"
#:
#: [[publish]]
#: series = "illumos-amd64"
#: name = "cockroach.tgz.sha256"
#: from_output = "/work/oxidecomputer/cockroach/cockroach.tgz.sha256"

set -o errexit
set -o pipefail
set -o xtrace

export BROWSERSLIST_IGNORE_OLD_DATA=1

source .github/buildomat/helios-setup.sh
gmake -j"$(nproc)" cockroach.tgz BUILDTYPE=release
