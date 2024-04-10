package interchain

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/celestiaorg/celestia-app/test/interchain/chainspec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/testreporter"
	"github.com/strangelove-ventures/interchaintest/v6/testutil"
	"github.com/stretchr/testify/require"
)

// TestICA verifies that Interchain Accounts work as expected.
func TestICA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestICA in short mode.")
	}

	client, network := interchaintest.DockerSetup(t)
	celestia := chainspec.GetCelestia(t)
	cosmosHub := chainspec.GetCosmosHub(t)
	relayer := getRelayerFactory(t).Build(t, client, network)
	pathName := fmt.Sprintf("%s-to-%s", celestia.Config().ChainID, cosmosHub.Config().ChainID)
	interchain := interchaintest.NewInterchain().
		AddChain(celestia).
		AddChain(cosmosHub).
		AddRelayer(relayer, getRelayerName()).
		AddLink(interchaintest.InterchainLink{
			Chain1:  celestia,
			Chain2:  cosmosHub,
			Relayer: relayer,
			Path:    pathName,
		})

	ctx := context.Background()
	reporter := testreporter.NewNopReporter().RelayerExecReporter(t)
	err := interchain.Build(ctx, reporter, interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = interchain.Close() })

	err = relayer.StartRelayer(ctx, reporter, pathName)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 2, celestia, cosmosHub)
	require.NoError(t, err)

	celestiaConnections, err := relayer.GetConnections(ctx, reporter, celestia.Config().ChainID)
	require.NoError(t, err)
	require.Len(t, celestiaConnections, 1)

	cosmosConnections, err := relayer.GetConnections(ctx, reporter, cosmosHub.Config().ChainID)
	require.NoError(t, err)
	require.Len(t, cosmosConnections, 2) // 2 connections: the first is connection-0 and the second is connection-localhost.
	cosmosConnection := cosmosConnections[0]

	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), math.NewInt(10_000_000_000), celestia, cosmosHub)
	err = testutil.WaitForBlocks(ctx, 5, celestia, cosmosHub)
	require.NoError(t, err)

	celestiaUser, cosmosUser := users[0], users[1]
	celestiaAddr := celestiaUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(celestia.Config().Bech32Prefix)
	cosmosAddr := cosmosUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(cosmosHub.Config().Bech32Prefix)
	fmt.Printf("celestiaAddr: %s, cosmosAddr: %v\n", celestiaAddr, cosmosAddr)

	registerICA := []string{
		cosmosHub.Config().Bin, "tx", "interchain-accounts", "controller", "register", cosmosConnection.ID,
		"--chain-id", cosmosHub.Config().ChainID,
		"--home", cosmosHub.HomeDir(),
		"--node", cosmosHub.GetRPCAddress(),
		"--from", cosmosUser.KeyName(),
		"--keyring-backend", keyring.BackendTest,
		"--fees", "20000uatom",
		"--yes",
	}
	stdout, stderr, err := cosmosHub.Exec(ctx, registerICA, nil)
	require.NoError(t, err)
	require.Empty(t, stderr)
	t.Logf("stdout %v", string(stdout))

	err = testutil.WaitForBlocks(ctx, 5, celestia, cosmosHub)
	require.NoError(t, err)

	queryICA := []string{
		cosmosHub.Config().Bin, "query", "interchain-accounts", "controller", "interchain-account", cosmosAddr, cosmosConnection.ID,
		"--chain-id", cosmosHub.Config().ChainID,
		"--home", cosmosHub.HomeDir(),
		"--node", cosmosHub.GetRPCAddress(),
	}
	stdout, stderr, err = cosmosHub.Exec(ctx, queryICA, nil)
	t.Logf("stdout %v\n", string(stdout))
	t.Logf("stderr %v\n", string(stderr))
	t.Logf("err %v\n", err)
	// require.NoError(t, err)
	_ = testutil.WaitForBlocks(ctx, 100, celestia, cosmosHub)
}
