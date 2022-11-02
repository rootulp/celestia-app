package keeper

import (
	"context"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/pkg/shares"
	"github.com/celestiaorg/celestia-app/x/payment/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	payForDataGasDescriptor = "pay for data"

	// GasPerMsgByte is the amount of gas to charge per byte of message data.
	// TODO: extract GasPerMsgByte as a parameter to this module.
	GasPerMsgByte  = 8
	GasPerMsgShare = appconsts.ShareSize * GasPerMsgByte
)

// Keeper handles all the state changes for the celestia-app module.
type Keeper struct {
	cdc codec.BinaryCodec
}

func NewKeeper(cdc codec.BinaryCodec) *Keeper {
	return &Keeper{
		cdc: cdc,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// MsgPayForData moves a user's coins to the module address and burns them.
func (k Keeper) PayForData(goCtx context.Context, msg *types.MsgPayForData) (*types.MsgPayForDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	gasToConsume := uint64(shares.MsgSharesUsed(int(msg.MessageSize)) * GasPerMsgShare)
	ctx.GasMeter().ConsumeGas(gasToConsume, payForDataGasDescriptor)

	ctx.EventManager().EmitEvent(
		types.NewPayForDataEvent(sdk.AccAddress(msg.Signer).String(), msg.GetMessageSize()),
	)

	return &types.MsgPayForDataResponse{}, nil
}
