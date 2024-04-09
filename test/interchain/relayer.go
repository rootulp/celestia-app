package interchain

import (
	"testing"

	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/strangelove-ventures/interchaintest/v6/relayer"
	"go.uber.org/zap/zaptest"
)

func getRelayerName() string {
	return "hermes-relayer"
}

func getRelayerFactory(t *testing.T) interchaintest.RelayerFactory {
	return interchaintest.NewBuiltinRelayerFactory(
		ibc.Hermes,
		zaptest.NewLogger(t),
		relayer.CustomDockerImage("informalsystems/hermes", "v1.8.2", "1000:1000"),
	)
}
