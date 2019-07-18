package stakingmint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper StakingMintKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case MsgLiquidUndelegate:
			return handleMsgLiquidUndelegate(ctx, msg, keeper)
		default:
			return sdk.ErrTxDecode("invalid message parse in staking ibc module").Result()
		}
	}
}

func handleMsgLiquidUndelegate(ctx sdk.Context, msg MsgLiquidUndelegate, keeper StakingMintKeeper) sdk.Result {
	tags, err := keeper.UndelegateLiquid(ctx, msg.DelegatorAddress)

	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}
