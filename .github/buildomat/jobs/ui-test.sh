#!/bin/bash
#:
#: name = "ui-test"
#: variety = "basic"
#: target = "ubuntu-22.04"

set -o errexit
set -o pipefail
set -o xtrace

source .github/buildomat/linux-setup.sh
install_output=$(npx @puppeteer/browsers install chrome@stable)
echo "${install_output}"
CHROME_BIN=$(awk '{ print $2 }' <<<"${install_output}")
export CHROME_BIN

sudo apt-get satisfy -y --no-install-recommends \
    "$(tr '\n' , <"$(dirname "${CHROME_BIN}")"/deb.deps)"

gmake -j"$(nproc)" ui-test || error=$?
gmake -j"$(nproc)" ui-lint
exit "${error:-0}"
