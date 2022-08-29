# celestia-app

[![Go Reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/celestiaorg/celestia-app)
[![GitHub Release](https://img.shields.io/github/v/release/celestiaorg/celestia-app)](https://github.com/celestiaorg/celestia-app/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/celestiaorg/celestia-app)](https://goreportcard.com/report/github.com/celestiaorg/celestia-app)
[![Lint](https://github.com/celestiaorg/celestia-app/actions/workflows/lint.yml/badge.svg)](https://github.com/celestiaorg/celestia-app/actions/workflows/lint.yml)
[![Tests / Code Coverage](https://github.com/celestiaorg/celestia-app/actions/workflows/test.yml/badge.svg)](https://github.com/celestiaorg/celestia-app/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/celestiaorg/celestia-app/branch/main/graph/badge.svg?token=CWGA4RLDS9)](https://codecov.io/gh/celestiaorg/celestia-app)

**celestia-app** is a blockchain application built using Cosmos SDK and [celestia-core](https://github.com/celestiaorg/celestia-core) in place of Tendermint.

## Diagram

```ascii
                ^  +-------------------------------+  ^
                |  |                               |  |
                |  |  State-machine = Application  |  |
                |  |                               |  |   celestia-app (built with Cosmos SDK)
                |  |            ^      +           |  |
                |  +----------- | ABCI | ----------+  v
Celestia        |  |            +      v           |  ^
validator or    |  |                               |  |
full consensus  |  |           Consensus           |  |
node            |  |                               |  |
                |  +-------------------------------+  |   celestia-core (fork of Tendermint Core)
                |  |                               |  |
                |  |           Networking          |  |
                |  |                               |  |
                v  +-------------------------------+  v
```

## Install

1. [Install Go](https://go.dev/doc/install) 1.18+
1. Clone this repo
1. Install the celestia-app CLI

    ```shell
    make install
    ```

## Usage

```sh
# Print help message
celestia-appd --help

# Create your own single node devnet
celestia-appd init test --chain-id test
celestia-appd keys add user1
celestia-appd add-genesis-account <address from above command> 10000000utia,1000token
celestia-appd gentx user1 1000000utia --chain-id test
celestia-appd collect-gentxs
celestia-appd start

# Post data to the local devnet
celestia-appd tx payment payForData [hexNamespace] [hexMessage] [flags]
```

See <https://docs.celestia.org/category/celestia-app> for more information

## Contributing

### Tools

1. Install [golangci-lint](https://golangci-lint.run/usage/install/)
1. Install [markdownlint](https://github.com/DavidAnson/markdownlint)

### Helpful Commands

```sh
# Build a new celestia-app binary and output to build/celestia-appd
make build

# Run tests
make test

# Format code with linters (this assumes golangci-lint and markdownlint are installed)
make fmt
```

## Careers

We are hiring Go engineers! Join us in building the future of blockchain scaling and interoperability. [Apply here](https://jobs.lever.co/celestia).
