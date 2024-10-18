#!/bin/bash
#:
#: name = "lint"
#: variety = "basic"
#: target = "ubuntu-22.04"

set -o errexit
set -o pipefail
set -o xtrace

source .github/buildomat/linux-setup.sh
gmake -j"$(nproc)" lint

# verify go.mod is up to date
go mod tidy
git diff --exit-code
