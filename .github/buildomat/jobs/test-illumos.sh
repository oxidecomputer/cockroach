#!/bin/bash
#:
#: name = "test-illumos"
#: variety = "basic"
#: target = "helios-2.0"

set -o errexit
set -o pipefail
set -o xtrace

source .github/buildomat/helios-setup.sh
gmake -j"$(nproc)" test
