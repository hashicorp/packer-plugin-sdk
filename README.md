# Packer Plugin SDK

This SDK enables building Packer plugins. This allows Packer's users to use both the officially-supported builders, provisioners, and post-processors, and custom in-house solutions.

Packer itself is a tool for building identical machine images for multiple platforms from a single source configuration. You can find more about Packer on its [website](https://www.packer.io) and [its GitHub repository](https://github.com/hashicorp/packer).

## Packer CLI Compatibility

Packer v1.7.0 or later is needed for this SDK. Versions of Packer prior to that release are still compatible with third-party plugins, but the plugins should use the plugin tooling from inside earlier versions of Packer to ensure complete API compatibility.

## Go Compatibility

The Packer Plugin SDK is built in Go, and uses the [support policy](https://golang.org/doc/devel/release.html#policy) of Go as its support policy. The two latest major releases of Go are supported by the SDK.

Currently, that means Go **1.14** or later must be used when building a provider with the SDK.

## Getting Started

See the [Extending Packer](https://www.packer.io/docs/extending) docs for a guided tour of plugin development.

## Documentation

See the [Extending Packer](https://www.packer.io/docs/extending) section on the Packer website.

## Packer Scope (Plugins VS Core)

### Packer Core

 - acts as an RPC _client_
 - interacts with the user
 - parses (HCL/JSON) configuration
 - manages build as whole, asks **plugin(s)** to manage the image lifecycle and modify the image being built.
 - discovers **plugin(s)** and their versions per configuration
 - manages **plugin** lifecycles (i.e. spins up & tears down plugin process)
 - passes relevant parts of parsed (valid JSON/HCL) and interpolated configuration to **plugin(s)**

### Packer Provider (via this SDK)

 - acts as RPC _server_
 - executes any domain-specific logic based on received parsed configuration. For builders this includes managing the vm lifecycle on a give hypervisor or cloud; for provisioners this involves calling the operation on the remote instance.
 - tests domain-specific logic via provided acceptance test framework
 - provides **Core** with template validation, artifact information, and information about whether the plugin process succeeded or failed.

## Migrating to SDK from built-in SDK

Migrating to the standalone SDK v1 is covered on the [Plugin SDK section](https://www.packer.io/docs/extend/plugin-sdk.html) of the website.

## Versioning

The Packer Plugin SDK is a [Go module](https://github.com/golang/go/wiki/Modules) versioned using [semantic versioning](https://semver.org/).

## Releasing

The Packer Plugin SDK is distributed as a Go module, so a minimal 'release' is just a git tag on main.

Releases can be triggered via CircleCI, or the release scripts run from one's own machine.

### Releasing via CircleCI


#### Changelog

`CHANGELOG.md` on the `main` branch must have a line of the following form:
```
# x.y.z (Upcoming)
```
where `x.y.z` is the SDK version you intend to release.


Underneath the `# x.y.z (Upcoming)` heading, please write human-readable entries describing the changes made in this version.

### Triggering a release

All commits to `main` trigger the `release` workflow on CircleCI, but manual approval is needed to proceed with the workflow.

#### Find all `main` workflows

Go to https://circleci.com/gh/hashicorp/workflows/packer-plugin-sdk/tree/main and click on the latest workflow.

It should show the latest `main` commit, which may be updates to the CHANGELOG.

The status should be `ON HOLD`.

#### Approve the `trigger-release` job

Find the `trigger-release` job on the workflow diagram, click on it, and click Approve.

#### Verify

The workflow will then run several other jobs to test and perform the release.

The main steps performed by the `release.sh` script are:
 - CHANGELOG.md: remove `(Upcoming)` and commit changes
 - tag `vx.y.z`
 - push to `main`

Verify that this has taken place correctly by inspecting the `main` branch once the release is complete.
### Releasing manually

### Prerequisites

Same as above in the CircleCI process.

### Releasing

Check out `main` locally and pull.

Run `make test` and any other steps you need to satisfy yourself that the release works.

Run `./scripts/release.sh`.

## Contributing

See [`.github/CONTRIBUTING.md`](https://github.com/hashicorp/packer-plugin-sdk/blob/master/.github/CONTRIBUTING.md)

## License

[Mozilla Public License v2.0](https://github.com/hashicorp/Packer-plugin-sdk/blob/master/LICENSE)
