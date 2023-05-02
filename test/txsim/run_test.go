//go:build !race

// known race in testnode
// ref: https://github.com/celestiaorg/celestia-app/issues/1369
package txsim_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/celestiaorg/celestia-app/test/txsim"
	"github.com/celestiaorg/celestia-app/test/util/testnode"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"

	blob "github.com/celestiaorg/celestia-app/x/blob/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestTxSimulator(t *testing.T) {
	testCases := []struct {
		name        string
		sequences   []txsim.Sequence
		expMessages map[string]int64
	}{
		{
			name:      "send sequence",
			sequences: []txsim.Sequence{txsim.NewSendSequence(2, 1000, 100)},
			// we expect at least 5 bank send messages within 30 seconds
			expMessages: map[string]int64{sdk.MsgTypeURL(&bank.MsgSend{}): 5},
		},
		{
			name:      "stake sequence",
			sequences: []txsim.Sequence{txsim.NewStakeSequence(1000)},
			expMessages: map[string]int64{
				sdk.MsgTypeURL(&staking.MsgDelegate{}):                     1,
				sdk.MsgTypeURL(&distribution.MsgWithdrawDelegatorReward{}): 5,
				// NOTE: this sequence also makes redelegations but because the
				// testnet has only one validator, this never happens
			},
		},
		{
			name: "blob sequence",
			sequences: []txsim.Sequence{
				txsim.NewBlobSequence(
					txsim.NewRange(100, 1000),
					txsim.NewRange(1, 3)),
			},
			expMessages: map[string]int64{sdk.MsgTypeURL(&blob.MsgPayForBlobs{}): 10},
		},
		{
			name: "multi blob sequence",
			sequences: txsim.NewBlobSequence(
				txsim.NewRange(1000, 1000),
				txsim.NewRange(3, 3),
			).Clone(4),
			expMessages: map[string]int64{sdk.MsgTypeURL(&blob.MsgPayForBlobs{}): 20},
		},
		{
			name: "multi mixed sequence",
			sequences: append(append(
				txsim.NewSendSequence(2, 1000, 100).Clone(3),
				txsim.NewStakeSequence(1000).Clone(3)...),
				txsim.NewBlobSequence(txsim.NewRange(1000, 1000), txsim.NewRange(1, 3)).Clone(3)...),
			expMessages: map[string]int64{
				sdk.MsgTypeURL(&bank.MsgSend{}):                            15,
				sdk.MsgTypeURL(&staking.MsgDelegate{}):                     2,
				sdk.MsgTypeURL(&distribution.MsgWithdrawDelegatorReward{}): 10,
				sdk.MsgTypeURL(&blob.MsgPayForBlobs{}):                     10,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			keyring, rpcAddr, grpcAddr := Setup(t)

			err := txsim.Run(
				ctx,
				[]string{rpcAddr},
				[]string{grpcAddr},
				keyring,
				9001,
				time.Second,
				tc.sequences...,
			)
			// Expect all sequences to run for at least 30 seconds without error
			require.True(t, errors.Is(err, context.DeadlineExceeded), err.Error())

			blocks, err := testnode.ReadBlockchain(context.Background(), rpcAddr)
			require.NoError(t, err)
			for _, block := range blocks {
				msgs, err := testnode.DecodeBlockData(block.Data)
				require.NoError(t, err, block.Height)
				for _, msg := range msgs {
					if _, ok := tc.expMessages[sdk.MsgTypeURL(msg)]; ok {
						tc.expMessages[sdk.MsgTypeURL(msg)]--
					}
				}
			}
			for msg, count := range tc.expMessages {
				if count > 0 {
					t.Errorf("missing %d messages of type %s (blocks: %d)", count, msg, len(blocks))
				}
			}
		})
	}
}

func Setup(t testing.TB) (keyring.Keyring, string, string) {
	t.Helper()
	genesis, keyring, err := testnode.DefaultGenesisState()
	require.NoError(t, err)

	tmCfg := testnode.DefaultTendermintConfig()
	tmCfg.RPC.ListenAddress = fmt.Sprintf("tcp://127.0.0.1:%d", testnode.GetFreePort())
	tmCfg.P2P.ListenAddress = fmt.Sprintf("tcp://127.0.0.1:%d", testnode.GetFreePort())
	tmCfg.RPC.GRPCListenAddress = fmt.Sprintf("tcp://127.0.0.1:%d", testnode.GetFreePort())

	node, app, cctx, err := testnode.New(
		t,
		testnode.DefaultParams(),
		tmCfg,
		true,
		genesis,
		keyring,
		"testnet",
	)
	require.NoError(t, err)

	cctx, stopNode, err := testnode.StartNode(node, cctx)
	require.NoError(t, err)

	appConf := testnode.DefaultAppConfig()
	appConf.GRPC.Address = fmt.Sprintf("127.0.0.1:%d", testnode.GetFreePort())
	appConf.API.Address = fmt.Sprintf("tcp://127.0.0.1:%d", testnode.GetFreePort())

	_, cleanupGRPC, err := testnode.StartGRPCServer(app, appConf, cctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		t.Log("tearing down testnode")
		require.NoError(t, stopNode())
		require.NoError(t, cleanupGRPC())
	})

	return keyring, tmCfg.RPC.ListenAddress, appConf.GRPC.Address
}
