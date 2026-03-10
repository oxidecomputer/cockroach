source .github/buildomat/versions.sh

GO_MAJOR=${GO_VERSION%.*}
NODE_MAJOR=${NODE_VERSION%%.*}

pfexec pkg install \
    /developer/build-essential /ooce/developer/cmake "/ooce/developer/go-${GO_MAJOR/./}@$GO_VERSION" "/ooce/runtime/node-$NODE_MAJOR@$NODE_VERSION"

pushd /work
mkdir bin
curl -sSfL --retry 10 -O "https://github.com/yarnpkg/yarn/releases/download/v$YARN_VERSION/yarn-$YARN_VERSION.js"
sha256sum -c "$OLDPWD/.github/buildomat/SHA256SUMS.helios"
mv "yarn-$YARN_VERSION.js" bin/yarn
chmod a+x bin/yarn
popd
export PATH="/work/bin:/opt/ooce/go-$GO_MAJOR/bin:/opt/ooce/node-$NODE_MAJOR/bin:$PATH"
