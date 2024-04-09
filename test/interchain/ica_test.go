package interchain

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/celestiaorg/celestia-app/test/interchain/chainspec"
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
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
	pathName := fmt.Sprintf("%s-to-%s", celestia.Config().Name, cosmosHub.Config().Name)
	ic := interchaintest.NewInterchain().
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
	err := ic.Build(ctx, reporter, interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,
	})
	require.NoError(t, err)
	// t.Cleanup(func() { _ = ic.Close() })

	err = relayer.CreateClients(ctx, reporter, pathName, ibc.CreateClientOptions{TrustingPeriod: "330h"})
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 2, celestia, cosmosHub)
	require.NoError(t, err)

	err = relayer.CreateConnections(ctx, reporter, pathName)
	require.NoError(t, err)

	err = relayer.StartRelayer(ctx, reporter, pathName)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 2, celestia, cosmosHub)
	require.NoError(t, err)

	celestiaConnections, err := relayer.GetConnections(ctx, reporter, celestia.Config().ChainID)
	require.NoError(t, err)
	require.Len(t, celestiaConnections, 1)

	cosmosConnections, err := relayer.GetConnections(ctx, reporter, cosmosHub.Config().ChainID)
	require.NoError(t, err)
	require.Len(t, cosmosConnections, 1)
	cosmosConnection := cosmosConnections[0]

	amount := math.NewIntFromUint64(uint64(10_000_000_000))
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), amount, celestia, cosmosHub)

	celestiaUser, strideUser := users[0], users[1]
	celestiaAddr := celestiaUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(celestia.Config().Bech32Prefix)
	strideAddr := strideUser.(*cosmos.CosmosWallet).FormattedAddressWithPrefix(cosmosHub.Config().Bech32Prefix)
	fmt.Printf("celestiaAddr: %s, strideAddr: %v\n", celestiaAddr, strideAddr)

	registerICA := []string{
		cosmosHub.Config().Bin, "tx", "interchain-accounts", "controller", "register", cosmosConnection.ID,
		"--chain-id", cosmosHub.Config().ChainID,
		"--home", cosmosHub.HomeDir(),
		"--node", cosmosHub.GetRPCAddress(),
	}
	stdout, _, err := cosmosHub.Exec(ctx, registerICA, nil)
	require.NoError(t, err)
	t.Log(stdout)

	err = testutil.WaitForBlocks(ctx, 100, celestia, cosmosHub)
	require.NoError(t, err)
	// version := icatypes.NewDefaultMetadataString(ibctesting.FirstConnectionID, ibctesting.FirstConnectionID)
	// msgRegisterInterchainAccount := controllertypes.NewMsgRegisterInterchainAccount(ibctesting.FirstConnectionID, strideAddr, version)
	// txResp := BroadcastMessages(t, ctx, celestia, stride, strideUser, msgRegisterInterchainAccount)
	// fmt.Printf("txResp %v\n", txResp)

	// celestiaHeight, err := celestia.Height(ctx)
	// require.NoError(t, err)

	// isChannelOpen := func(found *chantypes.MsgChannelOpenConfirm) bool {
	// 	return found.PortId == "icahost"
	// }
	// _, err = cosmos.PollForMessage(ctx, celestia, cosmos.DefaultEncoding().InterfaceRegistry, celestiaHeight, celestiaHeight+30, isChannelOpen)
	// require.NoError(t, err)
}
