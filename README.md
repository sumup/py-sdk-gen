<div align="center">

# py-sdk-gen

[![Stars](https://img.shields.io/github/stars/sumup/py-sdk-gen?style=social)](https://github.com/sumup/py-sdk-gen/)
[![CI Status](https://github.com/sumup/py-sdk-gen/workflows/CI/badge.svg)](https://github.com/sumup/py-sdk-gen/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/sumup/py-sdk-gen)](./LICENSE)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.1%20adopted-ff69b4.svg)](https://github.com/sumup/py-sdk-gen/tree/main/CODE_OF_CONDUCT.md)

`py-sdk-gen` is a highly opinionated OpenAPI specs to SDK generator for [Python](https://www.python.org/).

</div>

## Quickstart

Install the latest version of `py-sdk-gen` using

```sh
go install github.com/sumup/py-sdk-gen/cmd/py-sdk-gen@latest
```

And generate your SDK:

```sh
py-sdk-gen --mod github.com/me/mypackage --package mypackage --name 'My API' ./openapi.yaml
```

## Overview

`py-sdk-gen` generates structured SDK that is easy to navigate. Operations are grouped under tags and py-sdk-gen works under the assumption that each operation has one tag and one tag only.

When bootstrapping new project py-sdk-gen will generate all the necessary code for a valid SDK. On following runs it will update only code related to your OpenAPI specs but won't touch the client implementation and other files. This leaves you with the option to customize the client and add other features as necessary. You can opt out of this behavior using the `--force` flag.

## Usage

As a bade minimum, you will need to provide full path of your module (if you are bootstrapping new SDK), package name, and the source OpenAPI specs:

```sh
py-sdk-gen generate --mod github.com/me/mypackage --package mypackage --name 'My API' ./openapi.yaml
```

For further options see

```sh
py-sdk-gen help
```

### Maintainers

- [Matous Dzivjak](mailto:matous.dzivjak@sumup.com)
