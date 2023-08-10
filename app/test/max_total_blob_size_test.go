package app_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	appns "github.com/celestiaorg/celestia-app/pkg/namespace"
	"github.com/celestiaorg/celestia-app/pkg/square"
	"github.com/celestiaorg/celestia-app/test/util/testfactory"
	"github.com/celestiaorg/celestia-app/test/util/testnode"
	"github.com/celestiaorg/celestia-app/x/blob/types"
	blobtypes "github.com/celestiaorg/celestia-app/x/blob/types"
	sdk_tx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	coretypes "github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
)

const (
	mebibyte   = 1_048_576 // one mebibyte in bytes
	squareSize = 64
)

func FuzzMaxTotalBlobSizeAnteHandler(f *testing.F) {
	f.Add([]byte("seed data"))
	accounts := testfactory.GenerateAccounts(1)

	tmConfig := testnode.DefaultTendermintConfig()
	tmConfig.Mempool.MaxTxBytes = 10 * mebibyte

	cParams := testnode.DefaultParams()
	cParams.Block.MaxBytes = 10 * mebibyte

	cfg := testnode.DefaultConfig().
		WithAccounts(accounts).
		WithTendermintConfig(tmConfig).
		WithConsensusParams(cParams)

	f.Fuzz(func(t *testing.T, data []byte) {
		cctx, _, _ := testnode.NewNetwork(t, cfg)
		require.NoError(t, cctx.WaitForNextBlock())

		signer := blobtypes.NewKeyringSigner(cctx.Keyring, accounts[0], cctx.ChainID)

		ns1 := appns.MustNewV0(bytes.Repeat([]byte{1}, appns.NamespaceVersionZeroIDSize))
		blob := &tmproto.Blob{
			NamespaceId:      ns1.ID,
			NamespaceVersion: uint32(ns1.Version),
			Data:             data,
			ShareVersion:     0,
		}

		blobTx := newBlobTx(t, signer, cctx.GRPCClient, blob)
		res, err := types.BroadcastTx(context.TODO(), cctx.GRPCClient, sdk_tx.BroadcastMode_BROADCAST_MODE_BLOCK, blobTx)
		require.NoError(t, err)
		require.NotNil(t, res)

		sq, err := square.Construct([][]byte{blobTx}, appconsts.LatestVersion, squareSize)
		if res.TxResponse.Code == abci.CodeTypeOK {
			// verify that if the tx was accepted, the blob can fit in a square
			assert.NoError(t, err)
			assert.False(t, sq.IsEmpty())
			fmt.Printf("verified that tx was accepted and blob can fit in a square\n")
		} else {
			// verify that if the tx was rejected, the blob can not fit in a square
			assert.Error(t, err)
			fmt.Printf("verified that tx was rejected and blob can not fit in a square\n")
		}
	})
}

func newBlobTx(t *testing.T, signer *blobtypes.KeyringSigner, conn *grpc.ClientConn, blob *tmproto.Blob) coretypes.Tx {
	addr, err := signer.GetSignerInfo().GetAddress()
	require.NoError(t, err)

	msg, err := types.NewMsgPayForBlobs(addr.String(), blob)
	require.NoError(t, err)

	err = signer.QueryAccountNumber(context.TODO(), conn)
	require.NoError(t, err)

	options := []blobtypes.TxBuilderOption{blobtypes.SetGasLimit(1e9)} // set gas limit to 1 billion to avoid gas exhaustion
	builder := signer.NewTxBuilder(options...)
	stx, err := signer.BuildSignedTx(builder, msg)
	require.NoError(t, err)

	rawTx, err := signer.EncodeTx(stx)
	require.NoError(t, err)

	blobTx, err := coretypes.MarshalBlobTx(rawTx, blob)
	require.NoError(t, err)

	return blobTx
}
