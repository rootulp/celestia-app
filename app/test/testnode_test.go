package app_test

import (
	"context"
	"testing"

	"github.com/celestiaorg/celestia-app/v2/app"
	"github.com/celestiaorg/celestia-app/v2/test/util/testnode"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func Test_testnode(t *testing.T) {
	t.Run("testnode can start a network with default chain ID", func(t *testing.T) {
		testnode.NewNetwork(t, testnode.DefaultConfig())
	})
	t.Run("testnode can start with a custom MinGasPrice", func(t *testing.T) {
		want := "0.003utia"
		appConfig := testnode.DefaultAppConfig()
		appConfig.MinGasPrices = want
		config := testnode.DefaultConfig().WithAppConfig(appConfig)
		cctx, _, _ := testnode.NewNetwork(t, config)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		got, err := queryMinimumGasPrice(ctx, cctx.GRPCClient)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func queryMinimumGasPrice(ctx context.Context, grpcConn *grpc.ClientConn) (float64, error) {
	cfgRsp, err := nodeservice.NewServiceClient(grpcConn).Config(ctx, &nodeservice.ConfigRequest{})
	if err != nil {
		return 0, err
	}

	localMinCoins, err := sdktypes.ParseDecCoins(cfgRsp.MinimumGasPrice)
	if err != nil {
		return 0, err
	}
	localMinPrice := localMinCoins.AmountOf(app.BondDenom).MustFloat64()
	return localMinPrice, nil
}
