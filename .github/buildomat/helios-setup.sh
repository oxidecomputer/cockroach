source .github/buildomat/versions.sh

pfexec pkg install \
    /developer/build-essential /ooce/developer/cmake "/ooce/developer/go-122@$GO_VERSION" "/ooce/runtime/node-20@$NODE_VERSION"

pushd /work
mkdir bin
curl -sSfL --retry 10 -O "https://github.com/yarnpkg/yarn/releases/download/v$YARN_VERSION/yarn-$YARN_VERSION.js"
sha256sum --ignore-missing -c "$OLDPWD/.github/buildomat/SHA256SUMS"
mv "yarn-$YARN_VERSION.js" bin/yarn
chmod a+x bin/yarn
popd
export PATH="/work/bin:/opt/ooce/go-1.22/bin:/opt/ooce/node-20/bin:$PATH"
