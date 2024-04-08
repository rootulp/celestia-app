package chainspec

import (
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
)

var Stride = &interchaintest.ChainSpec{
	Name: "stride",
	ChainConfig: ibc.ChainConfig{
		Type:    "cosmos",
		Name:    "stride",
		ChainID: "stride-1",
		Images: []ibc.DockerImage{{
			Repository: "ghcr.io/strangelove-ventures/heighliner/stride",
			Version:    "v21.0.0",
			UidGid:     "1025:1025",
		}},
		Bin:           "strided",
		Bech32Prefix:  "stride",
		Denom:         "ustrd",
		GasPrices:     "0.0ustrd",
		GasAdjustment: 1.1,
	},
	NumFullNodes:  numFullNodes(),
	NumValidators: numValidators(),
}
