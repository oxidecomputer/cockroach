#!/bin/bash
#:
#: name = "lint"
#: variety = "basic"
#: target = "ubuntu-22.04"

set -o errexit
set -o pipefail
set -o xtrace

source .github/buildomat/linux-setup.sh

failed=0
gmake -j"$(nproc)" lint || ((++failed))
gmake -j"$(nproc)" generate || ((++failed))
go mod tidy || ((++failed))
git diff --exit-code || ((++failed))
((!failed)) || exit
