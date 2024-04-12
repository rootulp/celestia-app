package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/celestiaorg/celestia-app/v2/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/v2/test/util/testnode"
	"github.com/stretchr/testify/require"
)

func TestE2EThroughput(t *testing.T) {
	if os.Getenv("KNUU_NAMESPACE") != "test" {
		t.Skip("skipping e2e throughput test")
	}

	if os.Getenv("E2E_LATEST_VERSION") != "" {
		latestVersion = os.Getenv("E2E_LATEST_VERSION")
		_, isSemVer := ParseVersion(latestVersion)
		switch {
		case isSemVer:
		case latestVersion == "latest":
		case len(latestVersion) == 7:
		case len(latestVersion) >= 8:
			// assume this is a git commit hash (we need to trim the last digit to match the docker image tag)
			latestVersion = latestVersion[:7]
		default:
			t.Fatalf("unrecognised version: %s", latestVersion)
		}
	}

	t.Log("Running throughput test", "version", latestVersion)

	// create a new testnet
	testnet, err := New(t.Name(), seed, GetGrafanaInfoFromEnvVar())
	require.NoError(t, err)
	t.Cleanup(func() {
		t.Log("Cleaning up testnet")
		testnet.Cleanup()
	})

	// add 2 validators
	require.NoError(t, testnet.CreateGenesisNodes(2, latestVersion, 10000000,
		0, defaultResources))

	// obtain the GRPC endpoints of the validators
	gRPCEndpoints, err := testnet.RemoteGRPCEndpoints()
	require.NoError(t, err)
	t.Log("validators GRPC endpoints", gRPCEndpoints)

	// create txsim nodes and point them to the validators
	t.Log("Creating txsim nodes")
	// version of the txsim docker image to be used
	txsimVersion := "a92de72"

	err = testnet.CreateTxClients(txsimVersion, 1, "10000-10000", defaultResources, gRPCEndpoints)
	require.NoError(t, err)

	// start the testnet
	t.Log("Setting up testnet")
	require.NoError(t, testnet.Setup()) // configs, genesis files, etc.
	t.Log("Starting testnet")
	require.NoError(t, testnet.Start())

	// once the testnet is up, start the txsim
	t.Log("Starting txsim nodes")
	err = testnet.StartTxClients()
	require.NoError(t, err)

	// wait some time for the txsim to submit transactions
	time.Sleep(1 * time.Minute)

	t.Log("Reading blockchain")
	blockchain, err := testnode.ReadBlockchain(context.Background(), testnet.Node(0).AddressRPC())
	require.NoError(t, err)

	totalTxs := 0
	for _, block := range blockchain {
		require.Equal(t, appconsts.LatestVersion, block.Version.App)
		totalTxs += len(block.Data.Txs)
	}
	require.Greater(t, totalTxs, 10)
}
