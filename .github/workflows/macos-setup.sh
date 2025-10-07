SOURCE_DIR=$PWD
source .github/buildomat/versions.sh

# Pin cmake to <4.0. This is staving off the eventual need to either patch,
# update, or remove GEOS. https://github.com/oxidecomputer/cockroach/issues/18
#
# Method loosely based on https://github.com/actions/runner-images/pull/12791
brew uninstall cmake
cmake_commit="b4e46db74e74a8c1650b38b1da222284ce1ec5ce"
tap_name="local/pinned"
brew tap-new --no-git "$tap_name" >/dev/null
cmake_formula_dir="$(brew --repo "$tap_name")/Formula"
mkdir -p "$cmake_formula_dir"
curl -fsSL "https://raw.githubusercontent.com/Homebrew/homebrew-core/$cmake_commit/Formula/c/cmake.rb" -o "$cmake_formula_dir/cmake.rb"
brew install "$tap_name/cmake"

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
popd
