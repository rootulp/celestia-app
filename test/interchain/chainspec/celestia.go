package chainspec

import (
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/strangelove-ventures/interchaintest/v6/testutil"
)

const (
	celestiaDockerRepository = "ghcr.io/celestiaorg/celestia-app"
	// celestiaDockerTag is the docker tag for a Celestia image with ICA enabled.
	// https://github.com/celestiaorg/celestia-app/commit/21f53b3f7ea118a03d5c2928cd586cc4c62751ad
	// https://github.com/celestiaorg/celestia-app/pkgs/container/celestia-app/200172087?tag=21f53b3
	celestiaDockerTag = "21f53b3"
)

var Celestia = &interchaintest.ChainSpec{
	Name: "celestia",
	ChainConfig: ibc.ChainConfig{
		Type:                "cosmos",
		Name:                "celestia-app",
		ChainID:             "celestia",
		Bin:                 "celestia-appd",
		Bech32Prefix:        "celestia",
		Denom:               "utia",
		GasPrices:           "0.002utia",
		GasAdjustment:       1.5,
		TrustingPeriod:      "336hours",
		Images:              celestiaDockerImages(),
		ConfigFileOverrides: celestiaConfigFileOverrides(),
	},
	NumValidators: numValidators(),
	NumFullNodes:  numFullNodes(),
}

func celestiaDockerImages() []ibc.DockerImage {
	return []ibc.DockerImage{
		{
			Repository: celestiaDockerRepository,
			Version:    celestiaDockerTag,
			UidGid:     "10001:10001",
		},
	}
}

func celestiaConfigFileOverrides() map[string]any {
	txIndex := make(testutil.Toml)
	txIndex["indexer"] = "kv"

	storage := make(testutil.Toml)
	storage["discard_abci_responses"] = false

	configToml := make(testutil.Toml)
	configToml["tx_index"] = txIndex
	configToml["storage"] = storage

	result := make(map[string]any)
	result["config/config.toml"] = configToml
	return result
}
