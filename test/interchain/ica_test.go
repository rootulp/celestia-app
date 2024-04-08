package interchain

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/strangelove-ventures/interchaintest/v6/testreporter"
	"github.com/strangelove-ventures/interchaintest/v6/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	relayerName     = "relayerName"
	path            = "path"
	DefaultGasValue = 500_000_0000
)

// TestICA verifies that Interchain Accounts work as expected.
func TestICA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestICA in short mode.")
	}

	client, network := interchaintest.DockerSetup(t)
	celestia, icad := getChains(t)

	relayer := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
	).Build(t, client, network)
	ic := interchaintest.NewInterchain().
		AddChain(celestia).
		AddChain(icad).
		AddRelayer(relayer, relayerName).
		AddLink(interchaintest.InterchainLink{
			Chain1:  celestia,
			Chain2:  icad,
			Relayer: relayer,
			Path:    path,
		})

	ctx := context.Background()
	reporter := testreporter.NewNopReporter().RelayerExecReporter(t)
	err := ic.Build(ctx, reporter, interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = ic.Close() })

	err = relayer.CreateClients(ctx, reporter, path, ibc.CreateClientOptions{TrustingPeriod: "330h"})
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 2, celestia, icad)
	require.NoError(t, err)

	err = relayer.CreateConnections(ctx, reporter, path)
	require.NoError(t, err)

	err = relayer.StartRelayer(ctx, reporter, path)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 2, celestia, icad)
	require.NoError(t, err)

	connections, err := relayer.GetConnections(ctx, reporter, celestia.Config().ChainID)
	require.NoError(t, err)
	require.Len(t, connections, 1)

	connections, err = relayer.GetConnections(ctx, reporter, icad.Config().ChainID)
	require.NoError(t, err)
	require.Len(t, connections, 1)

	amount := math.NewIntFromUint64(uint64(10_000_000_000))
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), amount, celestia, icad)

	celestiaUser, icadUser := users[0], users[1]
	celestiaAddr := celestiaUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(celestia.Config().Bech32Prefix)
	icadAddr := icadUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(icad.Config().Bech32Prefix)

	fmt.Printf("celestiaAddr: %s, icadAddr: %v\n", celestiaAddr, icadAddr)

	registerICA := []string{
		icad.Config().Bin, "tx", "intertx", "register",
		"--from", icadAddr,
		"--connection-id", connections[0].ID,
		"--chain-id", icad.Config().ChainID,
		"--home", icad.HomeDir(),
		"--node", icad.GetRPCAddress(),
		"--keyring-backend", keyring.BackendTest,
		"-y",
	}
	stdout, stderr, err := icad.Exec(ctx, registerICA, nil)
	require.NoError(t, err)
	require.Empty(t, stderr)
	t.Log(string(stdout))

	// celestiaHeight, err := celestia.Height(ctx)
	// require.NoError(t, err)
	// // Wait for channel open confirm
	// isChannelFound := func(found *chantypes.MsgChannelOpenConfirm) bool { return found.PortId == "icahost" }
	// _, err = cosmos.PollForMessage(ctx, celestia, cosmos.DefaultEncoding().InterfaceRegistry, celestiaHeight, celestiaHeight+30, isChannelFound)
	// require.NoError(t, err)
	err = testutil.WaitForBlocks(ctx, 10, celestia, icad)
	require.NoError(t, err)

	queryICA := []string{
		icad.Config().Bin, "query", "intertx", "interchainaccounts", connections[0].ID, icadAddr,
		"--chain-id", icad.Config().ChainID,
		"--home", icad.HomeDir(),
		"--node", icad.GetRPCAddress(),
	}
	stdout, stderr, err = icad.Exec(ctx, queryICA, nil)
	require.NoError(t, err)
	require.Empty(t, stderr)
	t.Log(string(stdout))
}
