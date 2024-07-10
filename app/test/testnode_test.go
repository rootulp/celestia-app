package app_test

import (
	"context"
	"testing"

	"github.com/celestiaorg/celestia-app/v2/pkg/user"
	"github.com/celestiaorg/celestia-app/v2/test/util/genesis"
	"github.com/celestiaorg/celestia-app/v2/test/util/testnode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_testnode(t *testing.T) {
	t.Run("testnode can start a network with default chain ID", func(t *testing.T) {
		testnode.NewNetwork(t, testnode.DefaultConfig())
	})
	t.Run("testnode can start a network with a custom chain ID", func(t *testing.T) {
		chainID := "custom-chain-id"

		// Set the chain ID on genesis. If this isn't done, the default chain ID
		// is used for the validator which results in a "signature verification
		// failed" error.
		genesis := genesis.NewDefaultGenesis().
			WithChainID(chainID).
			WithValidators(genesis.NewDefaultValidator(testnode.DefaultValidatorAccountName)).
			WithConsensusParams(testnode.DefaultConsensusParams())

		config := testnode.DefaultConfig()
		config.WithChainID(chainID)
		config.WithGenesis(genesis)
		testnode.NewNetwork(t, config)
	})
	t.Run("testnode can start with a custom MinGasPrice", func(t *testing.T) {
		// want := "0.000006stake"
		config := testnode.DefaultConfig()
		appConfig := testnode.DefaultAppConfig()
		// appConfig.MinGasPrices = want
		config.WithAppConfig(appConfig)
		_, _, grpcAddr := testnode.NewNetwork(t, config)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		grpcConn := setup(t, ctx, grpcAddr)
		got, err := user.QueryMinimumGasPrice(ctx, grpcConn)
		require.NoError(t, err)
		assert.Equal(t, ".002utia", got)
	})
}

func setup(t *testing.T, ctx context.Context, grpcAddr string) *grpc.ClientConn {
	client, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	// this ensures we can't start the node without core connection
	client.Connect()
	if !client.WaitForStateChange(ctx, connectivity.Ready) {
		// hits the case when context is canceled
		t.Fatalf("couldn't connect to core endpoint(%s): %v", grpcAddr, ctx.Err())
	}
	return client
}
