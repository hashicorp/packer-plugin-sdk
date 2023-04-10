 # Latest Release
 Please refer to [releases](https://github.com/hashicorp/packer-plugin-sdk/releases) for latest CHANGELOG information.

## 0.3.1 (July 28, 2022)

## 0.3.0 (June 09, 2022)

* multistep/commonsteps: Add default timeouts to the GitGetter, HgGetter,
    S3Getter, and GcsGetter getters to mitigate against resource exhaustion
    when calling out to external command line applications.
* multistep/commonsteps: Disable support for the `X-Terraform-Get` header to
    mitigate against protocol switching, endless redirect, and configuration
    bypass abuse of custom HTTP response header processing.
* multistep/commonsteps: Update settings for the default go-getter client to
    prevent arbitrary host access via go-getter's path traversal, symlink
    processing, and command injection flaws.
* sdk: Bump github.com/hashicorp/go-getter/v2, github.com/hashicorp/go-
    getter/gcs/v2, github.com/hashicorp/go-getter/s3/v2 to address a number of
    security vulnerabilities as defined in
    [HCSEC-2022-13](https://discuss.hashicorp.com/t/hcsec-2022-13-multiple-vulnerabilities-in-go-getter-library/39930)

## 0.2.13 (May 11, 2022)

* cmd/packer-sdc: Update golang.org/x/tools to fix internal package errors when
    running code generation commands with Go 1.18
    [GH-108](https://github.com/hashicorp/packer-plugin-sdk/pull/108)

## 0.2.12 (May 03, 2022)

* provisioner/shell-local: Add `env` argument to pass env vars through a key
    value store [GH-98](https://github.com/hashicorp/packer-plugin-sdk/pull/98)
* provisioner/shell: Add `env` argument to pass env vars through a key value
    store [GH-98](https://github.com/hashicorp/packer-plugin-sdk/pull/98)
* sdk: Bump github.com/hashicorp/go-getter/v2  to v2.0.2
    [GH-102](https://github.com/hashicorp/packer-plugin-sdk/pull/102)
* sdk: Bump github.com/hashicorp/hcl/v2  to v2.12.0
    [GH-106](https://github.com/hashicorp/packer-plugin-sdk/pull/106)
* sdk: Update crypto/ssh pkg used by SSH communicator The existing ssh client
    used by the SSH communicator was relying on legacy key algorithms and could
    not connect to recent versions of openssh, or servers with a limited set of
    fips approved algorithms. [GH-107](https://github.com/hashicorp/packer-plugin-sdk/pull/107)

## 0.2.11 (December 17, 2021)

* sdk: The SourceImageID field for registry/image.Image is now optional;
    calling Image#Validate will nolonger error if SourceImageID is empty.
    [GH-90](https://github.com/hashicorp/packer-plugin-sdk/pull/90)

## 0.2.10 (December 15, 2021)

* sdk: Security release with Go 1.7.5 to address [CVE-2021-44717](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-44717). [GH-93](https://github.com/hashicorp/packer-plugin-sdk/pull/93)

## 0.2.9 (November 11, 2021)

* sdk: packer-sdc fix plugin-check command for provisioner only plugins. [#88](https://github.com/hashicorp/packer-plugin-sdk/pull/88)

## 0.2.8 (November 10, 2021)

* sdk: packer-sdc add plugin-check command to check plugin validity, try `packer-sdc plugin-check -h` for help [#85](https://github.com/hashicorp/packer-plugin-sdk/pull/85)

## 0.2.7 (October 07, 2021)

* sdk: Make step_download store the used iso_url in state. [GH-84]
* sdk: Add `floppy_content` parameter in `FloppyConfig` [GH-82]

## 0.2.6 (September 30, 2021)

* sdk: Add SourceImageID field to registry/image object for HCP Packer registry ancestry support [GH-81]

## 0.2.5 (September 07, 2021)

* sdk: Bump SDK to Go 1.17
* sdk: Add `registry/image` package to support the creation of HCP Packer registry images from Packer Artifacts. [#76](https://github.com/hashicorp/packer-plugin-sdk/pull/76)

## 0.2.4 (August 31, 2021)

* sdk: Use xdg basedir spec on unix [#68](https://github.com/hashicorp/packer-plugin-sdk/pull/68) [#73](https://github.com/hashicorp/packer-plugin-sdk/pull/73)
* packer-sdc: add provisioner templates to for docs [#77](https://github.com/hashicorp/packer-plugin-sdk/pull/77)

## 0.2.3 (June 03, 2021)
* CDConfig: Add `cd_content` field for file templating for cd files [#61](https://github.com/hashicorp/packer-plugin-sdk/pull/61)

## 0.2.2 (May 14, 2021)
* Update masterzen/winrm dependency to allow NTLM support for winrm_no_proxy [GH-66]
* StepCreateCD: Clean up temporary directory and add more robust tests [GH-62]

## 0.2.1 (May 07, 2021)

* Update go-getter to v2.0.0 to fix godep compilation issue

## 0.2.0 (April 16, 2021)

* Add packer-sdc command that will help Plugin maintainers and Packer maintainers
    to generate the docs and the HCL2 glue code.

## 0.1.3 (April 07, 2021)

* Merge pull request #51 from hashicorp/cleanup_acctests
* Remove packer core dependencies

## 0.1.2 (April 01, 2021)

* core: Pin SDK to Golang 1.16 [[GH-48](https://github.com/hashicorp/packer-plugin-sdk/pull/48)]
* core: Update Packer to v1.7.1


## 0.1.1 (March 31, 2021)

### Notes

In release [v0.0.12](#0012-february-11-2021) a backwards incompatible change was introduced to the
    packer-plugin-sdk with the update to v1.2.4 for the `ugorji/go/codec`
    package. Plugins built with a version of the Packer SDK prior to v0.0.12
    are encouraged to update to the latest possible version of the SDK to
    prevent potential codec marshalling issues with Packer v1.7.0 and higher.

### Features

* commonsteps/http_config: Add `http_content` configuration option as an
    alternative method for serving static HTTP content. This option works
    similar to `http_directory` but has the ability to serve files that include
    Go templating variables that can be interpolated at runtime by Packer core.
    [[GH-43](https://github.com/hashicorp/packer-plugin-sdk/pull/43)]

### Improvements

* didyoumean: Add a "did you mean" package to help find a name from a set of
    predefined suggestions. [[GH-43](https://github.com/hashicorp/packer-plugin-sdk/pull/43)]

### Bugs fixes

* bootcommand: Fix pageUp and pageDown boot command usb key strokes.
    [[GH-46](https://github.com/hashicorp/packer-plugin-sdk/pull/46)]

## 0.1.0 (February 18, 2021)

* core: Update Packer to v1.7.0 [[GH-39](https://github.com/hashicorp/packer-plugin-sdk/pull/39)]

## 0.0.14 (February 17, 2021)

### Features

* plugin version validation: when no version is passed, default it to 0.0.0 for dev purposes [[GH-36](https://github.com/hashicorp/packer-plugin-sdk/pull/36)]
* update packer dependency [[GH-37](https://github.com/hashicorp/packer-plugin-sdk/pull/37)]

## 0.0.12 (February 11, 2021)

### Features
* core: Update ugorji/go/codec to v1.2.4 [[GH-31](https://github.com/hashicorp/packer-plugin-sdk/pull/31)]

## 0.0.11 (February 04, 2021)

### Features
* plugin: Introduce the concept of APIVersionMajor and APIVersionMinor to the
    RPC address to allow for backward compatible changes with the Packer API
    protocol. [GH-21]

## 0.0.10 (January 26, 2021)

### Improvements
* acctest/pluginacc:  Add a generic plugin acceptance test case structure [GH-28]
* packer/ui: Update UI mock to keep  output messages [GH-27]
* plugin: Add `api_version` to the plugin describe command output [GH-24]

## 0.0.9 (January 22, 2021)

### Bug fixes
* rpc/datasource: Fix error return on datasource RPC. [GH-23]

## 0.0.8 (January 21, 2021)

### Improvements
* acctest/datasource: Add Setup function for datasource acceptance testing.
    [GH-22]
* template/interpolate/aws: Add support for getting secrets of type number to
    secretsmanagers function. [GH-18]

## 0.0.7 (January 15, 2021)

### Features

* packer/datasource: Add support for new `datasource` plugin type. [GH-6]
    [GH-9] [GH-15]

### Improvements

* sdk/tests: Fix acceptance for various packages. [GH-10] [GH-13]

## 0.0.6 (January 7, 2021)

* sdk: Initial release


