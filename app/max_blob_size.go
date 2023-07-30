package app

import (
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MaxBlobSize returns an upper bound for the maximum blob size that can be
// contained in a single Celestia block. Blobs that are larger than this value
// should be rejected.
func (app *App) MaxBlobSize(ctx sdk.Context) int {
	maxSquareSize := app.GovSquareSizeUpperBound(ctx)
	maxShares := maxSquareSize * maxSquareSize
	// Subtract one from maxShares because at least one share must be occupied
	// by the PFB tx associated with this blob.
	maxBlobShares := maxShares - 1
	maxShareBytes := maxBlobShares * appconsts.ContinuationSparseShareContentSize

	// TODO(rootulp): get MaxBytes consensus params from core
	maxBlockBytes := appconsts.DefaultMaxBytes

	return min(maxShareBytes, maxBlockBytes)
}

// min returns the smaller of a and b.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
