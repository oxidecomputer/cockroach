SOURCE_DIR=$PWD
source .github/buildomat/versions.sh

brew install autoconf coreutils make

mkdir "$HOME/toolchain"
pushd "$HOME/toolchain"
curl -sSfL --retry 10 -O "https://go.dev/dl/go$GO_VERSION.darwin-arm64.tar.gz"
curl -sSfL --retry 10 -O "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-darwin-arm64.tar.xz"
curl -sSfL --retry 10 -O "https://github.com/yarnpkg/yarn/releases/download/v$YARN_VERSION/yarn-$YARN_VERSION.js"
sha256sum -c "$SOURCE_DIR/.github/buildomat/SHA256SUMS.macos"
tar xf "go$GO_VERSION.darwin-arm64.tar.gz"
tar xf "node-v$NODE_VERSION-darwin-arm64.tar.xz"
mv "yarn-$YARN_VERSION.js" "node-v$NODE_VERSION-darwin-arm64/bin/yarn"
chmod a+x "node-v$NODE_VERSION-darwin-arm64/bin/yarn"
export PATH="$PWD/go/bin:$PWD/node-v$NODE_VERSION-darwin-arm64/bin:$PATH"
popd
