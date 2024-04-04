# Interchaintest

This folder contains tests that use [interchaintest](https://github.com/strangelove-ventures/interchaintest) to assert IBC features work as expected. Candidates for testing include:

1. Interchain Accounts (ICA)
1. Packet Forward Middleware (PFM)
1. Relayer Incentivization Middleware

## Usage

Run the tests via

```bash
make test-interchain
```

## Contributing

If you have local modifications that you would like to test via interchaintest, you'll need to create a new Celestia Docker image with your modifications. CI should automatically create a Docker image and publish it to GHCR if you create a PR against celestia-app. If that doesn't work, you can manually create an image via:

```shell
# make local modifications and commit them
make build-ghcr-docker
make publish-ghcr-docker
```

After you have a new Docker image with your modifications, you must update the test to reference the new image.
