source .github/buildomat/versions.sh

sudo apt-get install -y --no-install-recommends \
    build-essential autoconf cmake libedit-dev ncurses-dev

pushd /work
curl -sSfL --retry 10 -O "https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz"
curl -sSfL --retry 10 -O "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.xz"
curl -sSfL --retry 10 -O "https://github.com/yarnpkg/yarn/releases/download/v$YARN_VERSION/yarn-$YARN_VERSION.js"
sha256sum --ignore-missing -c "$OLDPWD/.github/buildomat/SHA256SUMS"
tar xf "go$GO_VERSION.linux-amd64.tar.gz"
tar xf "node-v$NODE_VERSION-linux-x64.tar.xz"
mv "yarn-$YARN_VERSION.js" "node-v$NODE_VERSION-linux-x64/bin/yarn"
chmod a+x "node-v$NODE_VERSION-linux-x64/bin/yarn"
popd
export PATH="/work/go/bin:/work/node-v$NODE_VERSION-linux-x64/bin:$PATH"
