package chainspec

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// GetStride returns a CosmosChain for Stride.
func GetStride(t *testing.T) *cosmos.CosmosChain {
	factory := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{stride})
	chains, err := factory.Chains(t.Name())
	require.NoError(t, err)
	return chains[0].(*cosmos.CosmosChain)
}

var stride = &interchaintest.ChainSpec{
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
		GasPrices:     "0.1ustrd",
		GasAdjustment: 1.1,
		ModifyGenesis: ModifyGenesisStride(),
	},
	NumFullNodes:  numFullNodes(),
	NumValidators: numValidators(),
}

const (
	StrideAdminAccount  = "admin"
	StrideAdminMnemonic = "tone cause tribe this switch near host damage idle fragile antique tail soda alien depth write wool they rapid unfold body scan pledge soft"
)

// ModifyGenesisStride assumes there is only 1 validator.
// Confirmed this works per https://www.diffchecker.com/hX3DbmSx/
func ModifyGenesisStride() func(ibc.ChainConfig, []byte) ([]byte, error) {
	return func(cfg ibc.ChainConfig, input []byte) ([]byte, error) {
		genesis := make(map[string]interface{})
		if err := json.Unmarshal(input, &genesis); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genesis file: %w", err)
		}

		// Example of how to override some config
		// DayEpochIndex := 1
		// DayEpochLen := "100s"
		// if err := dyno.Set(genesis, DayEpochLen, "app_state", "epochs", "epochs", DayEpochIndex, "duration"); err != nil {
		// 	return nil, err
		// }

		// TODO: override the genesis file based on how stride/cmd/consumer.go does it.
		// The add consumer section makes this diff to the genesis.json: https://www.diffchecker.com/OUcyGRYB/

		result, err := json.Marshal(genesis)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal genesis bytes to json: %w", err)
		}
		return result, nil
	}
}
