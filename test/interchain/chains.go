package interchain

import (
	"testing"

	"github.com/celestiaorg/celestia-app/test/interchain/chainspec"
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// getChains returns two chains for testing: celestia and stride.
func getChains(t *testing.T) (celestia *cosmos.CosmosChain, stride *cosmos.CosmosChain) {
	chainSpecs := []*interchaintest.ChainSpec{chainspec.Celestia, chainspec.Stride}
	chainFactory := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), chainSpecs)
	chains, err := chainFactory.Chains(t.Name())
	require.NoError(t, err)
	return chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)
}