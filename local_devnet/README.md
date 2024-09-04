# Local devnet

This directory contains a local devnet that contains Prometheus, Grafana, OTEL collector, and 4 validators. The devnet uses a Docker image that is built locally using the code in the root directory. If you make any code changes, you must [build the images](#build-the-images) again so that the change is reflected in the local devnet.

Note: If you edit the `go.mod` to contain a local package, for example, using a local version of celestia-core:

```go
replace github.com/celestiaorg/celestia-core => ../celestia-core
```

You will need to push that version to GitHub, then take the commit and change it in the `go.mod` to:

```go
replace github.com/celestiaorg/celestia-core => github.com/fork_name/celestia-core <commit>
```

To have those changes reflected when you [build the images](#build-the-images). Otherwise, the build will fail because the directory containing the changes is not part of the Docker context.

## Usage

### Prerequisites

1. Install [Docker](https://docs.docker.com/get-docker/).
1. Verify that Docker is up and running:

```shell
docker run hello-world
```

should return:

```text
Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:
 1. The Docker client contacted the Docker daemon.
 2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (arm64v8)
 3. The Docker daemon created a new container from that image which runs the
    executable that produces the output you are currently reading.
 4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:
 $ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker ID:
 https://hub.docker.com/

For more examples and ideas, visit:
 https://docs.docker.com/get-started/
```

### Build the images

```shell
docker compose build
```

### Start the devnet

To run all the validators and the telemetry stack:

```shell
docker compose up --detach
```

If you want to run just a specific instance, you can specify its name: `core0`, `core1`, `core2`, `core3`, which are the validators. Then, `prometheus`, `grafana`, `otel-collector` for the telemetry.

Example:

```shell
docker-compose up --detach core0 core1
```

Will run only two validators: `core0` and `core1`, and no telemetry.

Note: Starting `core0` is always needed because it is the only genesis validator. If you don't start it, then the network won't start. For the rest of the workloads, they're optional and any combination of them is allowed.

### Stop the devnet

```shell
docker-compose stop
```

### Delete the devnet

```shell
docker-compose down
```

## Monitoring

Monitoring is preconfigured in the `celestia-app/config.toml`. This means that you will have access to the metrics the moment you spin up the devnet along with telemetry.

### Accessing grafana

Grafana is exposed in `localhost:3000`. The default credentials are `admin:admin`. Then, you will find predefined data sources to get the metrics from.

For the dashboards, if you create fresh ones and save them, they will be saved on your machine under `telemetry/grafana`.

## Updating the configuration

The four validators use the same genesis, the same comet's config `celestia-app/config.toml`, and the same app config `celestia-app/app.toml`. If you make a change on one of them, you will make the change on all four validators.

Note: if the network is already running, and you make a change in the files, the changes will not be reflected until you [stop the devnet](#stop-the-devnet), then [start](#start-the-devnet) it again. Also, they will be reflected if you [delete it](#delete-the-devnet) and [start](#start-the-devnet) it again.

## Sending transactions

By default, `core0` is exposed to the host network. This means that the RPC is exposed under `localhost:26657` and the gRPC under `localhost:9090` and can be used to query the network. For example, if you have `celestia-appd` already installed:

```shell
celestia-appd query staking validators
```

will return the four validators.

To send transactions, you will need to first send some tokens to your account:

```shell
celestia-appd tx bank send core0 <your_address> <amount> --fees 21000utia --chain-id local_devnet --keyring-backend test --keyring-dir local_devnet/celestia-app/core0 --broadcast-mode block
```

For example, `<amount>` can be `100utia`.

Then, you can start sending any transaction you want.
