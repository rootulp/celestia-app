package interchain

import (
	"testing"

	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/strangelove-ventures/interchaintest/v6/relayer"
	"go.uber.org/zap/zaptest"
)

func getRelayerName() string {
	return "cosmos-relayer"
}

const (
	dockerRepo    = "ghcr.io/cosmos/relayer"
	dockerVersion = "v2.4.1"
	uidGid        = "100:1000"
)

func getRelayerFactory(t *testing.T) interchaintest.RelayerFactory {
	return interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
		relayer.CustomDockerImage(dockerRepo, dockerVersion, uidGid),
	)
}
