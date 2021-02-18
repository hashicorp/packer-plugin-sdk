## 0.1.0 (Upcoming)

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


