package chainspec

import (
	"testing"

	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	cosmosDockerRepository = "ghcr.io/strangelove-ventures/heighliner/gaia"
	cosmosDockerVersion    = "v15.1.0"
)

// GetCosmosHub returns a CosmosChain for the CosmosHub.
func GetCosmosHub(t *testing.T) *cosmos.CosmosChain {
	factory := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{cosmosHub})
	chains, err := factory.Chains(t.Name())
	require.NoError(t, err)
	return chains[0].(*cosmos.CosmosChain)
}

var cosmosHub = &interchaintest.ChainSpec{
	Name: "gaia",
	ChainConfig: ibc.ChainConfig{
		Type:                   "cosmos",
		Name:                   "gaia",
		ChainID:                "cosmoshub-4",
		Bin:                    "gaiad",
		Bech32Prefix:           "cosmos",
		Denom:                  "uatom",
		GasPrices:              "0.01uatom",
		GasAdjustment:          1.3,
		TrustingPeriod:         "504hours",
		NoHostMount:            false,
		Images:                 cosmosDockerImages(),
		UsingNewGenesisCommand: true,
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
