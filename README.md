This is Oxide's long-term maintenance branch/fork of [CockroachDB 22.1](https://github.com/cockroachdb/cockroach/tree/release-22.1).

Oxide uses CockroachDB for control plane data storage on the Oxide Cloud Computer, which uses [illumos](https://illumos.org) (specifically [Helios](https://github.com/oxidecomputer/helios)) as the underlying operating system. We launched our product with CockroachDB 22.1. After Cockroach Labs' announcement that they will change to a strictly proprietary (source-available) model, we made the decision to continue self-supporting on this BSL-licensed version for the foreseeable future. For more context, see [RFD 110](https://rfd.shared.oxide.computer/rfd/110) and [RFD 508](https://rfd.shared.oxide.computer/rfd/508).

The primary goal of this branch is to keep the wheels of building and testing CockroachDB rolling smoothly to enable our ability to self-support. Our product runs illumos, but we also support development of our product on Linux and macOS, so it's important to us that our bug fixes that go into the illumos build also end up on developer machines too. Currently we have builds for illumos, Linux, and macOS; and run the test suite on illumos and Linux.

You're welcome to use this branch under [the same terms we are](./licenses/BSL.txt), but note that we're unable to provide any support outside the context of its use in our product.

## Binary downloads

Binaries for each commit are available via Buildomat's published artefact APIs:

```
https://buildomat.eng.oxide.computer/public/file/oxidecomputer/cockroach/$SERIES/$COMMIT/cockroach.tgz
```

where:

- `$SERIES` is one of `illumos-amd64`, `linux-amd64`, or `darwin-amd64` (or `darwin-arm64`, both refer to the same universal binary)
- `$COMMIT` is the full 40-character commit hash

A SHA-256 checksum is also available at `cockroach.tgz.sha256`.

## Major changes from upstream

- We've removed all CCL-licensed code. The "Cockroach Community License" primarily covers enterprise features of the database. The split between BSL-licensed and CCL-licensed code was fairly clean, but we're not planning to maintain it, so we shouldn't keep it around.
- The repository no longer needs to be cloned into `$GOPATH/github.com/cockroachdb/cockroach`; everything works without setting `$GOPATH`.
- We're in the process of removing code specific to Cockroach Labs development processes (such as their tools for automated GitHub issues and pull requests). We're also not using the newer Bazel-based build system, and are planning on removing the related files eventually.
- Our Linux builds are built with an Ubuntu 22.04 toolchain, so they won't run on systems with glibc < 2.35 or GCC < 11.
