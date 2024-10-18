source .github/buildomat/versions.sh

pfexec pkg install \
    /developer/build-essential /ooce/developer/cmake /ooce/developer/go-117 /ooce/runtime/node-16

pushd /work
mkdir bin
curl -sSfL --retry 10 -O "https://github.com/yarnpkg/yarn/releases/download/v$YARN_VERSION/yarn-$YARN_VERSION.js"
sha256sum --ignore-missing -c "$OLDPWD/.github/buildomat/SHA256SUMS"
mv "yarn-$YARN_VERSION.js" bin/yarn
chmod a+x bin/yarn
popd
export PATH="/work/bin:/opt/ooce/go-1.17/bin:$PATH"
