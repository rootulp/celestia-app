package ante

import (
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	blobtypes "github.com/celestiaorg/celestia-app/x/blob/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MaxBlobSizeDecorator helps to prevent a PFB from being included in a block
// but not fitting in a data square.
type MaxBlobSizeDecorator struct {
	k BlobKeeper
}

func NewMaxBlobSizeDecorator(k BlobKeeper) MaxBlobSizeDecorator {
	return MaxBlobSizeDecorator{k}
}

// AnteHandle implements the AnteHandler interface. It checks to see
// if the transaction contains a MsgPayForBlobs and if so, checks that
// the total blob size in the MsgPayForBlobs are less than the max blob size.
func (d MaxBlobSizeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	upperBound := d.blobSizeUpperBound(ctx)
	for _, m := range tx.GetMsgs() {
		if pfb, ok := m.(*blobtypes.MsgPayForBlobs); ok {
			totalBlobSize := sum(pfb.BlobSizes)
			if totalBlobSize > upperBound {
				return ctx, errors.Wrapf(blobtypes.ErrBlobSizeTooLarge, "total blob size %d exceeds upper bound %d", totalBlobSize, upperBound)
			}
		}
	}

	return next(ctx, tx, simulate)
}

// blobSizeUpperBound returns an upper bound for the number of bytes available
// for blobs in a data square based on state parameters (namely the max square
// size). Note it is possible that txs with a total blobSize less than this
// upper bound still fail to be included in a block due to overhead from the PFB
// tx and/or padding shares. As a result, this upper bound should only be used
// to reject transactions that are guaranteed to be too large.
func (d MaxBlobSizeDecorator) blobSizeUpperBound(ctx sdk.Context) int {
	squareSize := d.getMaxSquareSize(ctx)
	return squareBytes(squareSize)
}

func (d MaxBlobSizeDecorator) getMaxSquareSize(ctx sdk.Context) int {
	// NOTE: it is possible to remove upperBound if we enforce that GovMaxSquareSize <= MaxSquareSize
	// See https://github.com/celestiaorg/celestia-app/pull/2203
	upperBound := appconsts.SquareSizeUpperBound(ctx.ConsensusParams().Version.AppVersion)
	govSquareSize := d.k.GovMaxSquareSize(ctx)
	return min(upperBound, int(govSquareSize))
}

// sum returns the total size of the given sizes.
func sum(sizes []uint32) (total int) {
	for _, size := range sizes {
		total += int(size)
	}
	return total
}

// squareBytes returns the number of bytes in a square for the given squareSize.
func squareBytes(squareSize int) int {
	numShares := squareSize * squareSize
	return numShares * appconsts.ShareSize
}

// min returns the minimum of two ints. This function can be removed once we
// upgrade to Go 1.21.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}