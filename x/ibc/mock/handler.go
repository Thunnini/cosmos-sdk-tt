package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case MsgRelay:
			return handleMsgPacket(ctx, keeper, msg)
		default:
			return sdk.ErrTxDecode("invalid message parse in staking ibc module").Result()
		}
	}
}

func handleMsgPacket(ctx sdk.Context, keeper Keeper, msg MsgRelay) sdk.Result {
	err := keeper.ReceivePacket(ctx, "", MockPacket{
		Data: msg.Data,
	})
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgCreateClient(ctx sdk.Context, keeper Keeper, msg MsgCreateClient) sdk.Result {
	err := keeper.CreateClient(ctx, msg.ClientID)
	if err != nil {
		return err.Result()
	}

	// TODO: events
	return sdk.Result{}
}
