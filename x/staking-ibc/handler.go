package stakingibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper StakingIBCKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case MsgIBCDelegate:
			return handleMsgIBCDelegate(ctx, msg, keeper)
		default:
			return sdk.ErrTxDecode("invalid message parse in staking ibc module").Result()
		}
	}
}

func handleMsgIBCDelegate(ctx sdk.Context, msg MsgIBCDelegate, keeper StakingIBCKeeper) sdk.Result {
	tags, err := keeper.Delegate(ctx, msg.DelegatorAddress, msg.ValidatorAddress, msg.Amount, msg.DestAddress)

	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}
