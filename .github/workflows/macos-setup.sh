SOURCE_DIR=$PWD
source .github/buildomat/versions.sh

brew install autoconf coreutils make

mkdir "$HOME/toolchain"
pushd "$HOME/toolchain"
curl -sSfL --retry 10 -O "https://go.dev/dl/go$GO_VERSION.darwin-arm64.tar.gz"
curl -sSfL --retry 10 -O "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-darwin-arm64.tar.xz"
curl -sSfL --retry 10 -O "https://github.com/yarnpkg/yarn/releases/download/v$YARN_VERSION/yarn-$YARN_VERSION.js"
sha256sum --ignore-missing -c "$SOURCE_DIR/.github/buildomat/SHA256SUMS"
tar xf "go$GO_VERSION.darwin-arm64.tar.gz"
tar xf "node-v$NODE_VERSION-darwin-arm64.tar.xz"
mv "yarn-$YARN_VERSION.js" "node-v$NODE_VERSION-darwin-arm64/bin/yarn"
chmod a+x "node-v$NODE_VERSION-darwin-arm64/bin/yarn"
export PATH="$PWD/go/bin:$PWD/node-v$NODE_VERSION-darwin-arm64/bin:$PATH"

# Apply patch to fix golang/go#53000
pushd go/src
patch -p2 <"$SOURCE_DIR/.github/workflows/e66f895667cd51d0d28c42d369a803c12db8bb35.patch"
go build cmd/cgo
popd
mv go/src/cgo go/pkg/tool/darwin_arm64/cgo

popd
