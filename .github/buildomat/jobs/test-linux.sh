#!/bin/bash
#:
#: name = "test-linux"
#: variety = "basic"
#: target = "ubuntu-22.04"

set -o errexit
set -o pipefail
set -o xtrace

source .github/buildomat/linux-setup.sh
gmake -j"$(nproc)" test
