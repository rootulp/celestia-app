package chainspec

import (
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
)

const (
	cosmosDockerRepository = "ghcr.io/strangelove-ventures/heighliner/gaia"
	cosmosDockerVersion    = "v15.1.0"
)

var Cosmos = &interchaintest.ChainSpec{
	Name: "gaia",
	ChainConfig: ibc.ChainConfig{
		Type:                   "cosmos",
		Name:                   "gaia",
		ChainID:                "cosmoshub-4",
		Bin:                    "gaiad",
		Bech32Prefix:           "cosmos",
		Denom:                  "uatom",
		GasPrices:              "0.01uatom",
		GasAdjustment:          100.0,
		TrustingPeriod:         "504hours",
		NoHostMount:            false,
		UsingNewGenesisCommand: true,
		Images:                 cosmosDockerImages(),
	},
	NumValidators: numValidators(),
	NumFullNodes:  numFullNodes(),
}

func cosmosDockerImages() []ibc.DockerImage {
	return []ibc.DockerImage{
		{
			Repository: cosmosDockerRepository,
			Version:    cosmosDockerVersion,
			UidGid:     "1025:1025",
		},
	}
}
