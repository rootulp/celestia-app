package chainspec

import (
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
)

// ICAD is the chain spec for the Inter-Chain Accounts Demo chain.
var ICAD = &interchaintest.ChainSpec{
	Name: "icad",
	ChainConfig: ibc.ChainConfig{
		Images:                 []ibc.DockerImage{{Repository: "ghcr.io/cosmos/ibc-go-icad", Version: "v0.5.0"}},
		UsingNewGenesisCommand: true,
	},
	NumValidators: numValidators(),
	NumFullNodes:  numFullNodes(),
}
