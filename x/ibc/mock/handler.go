package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func HandleMsgCreateClient(ctx sdk.Context, keeper Keeper, msg MsgCreateClient) sdk.Result {
	err := keeper.CreateClient(ctx, msg.ClientID)
	if err != nil {
		return err.Result()
	}

	// TODO: events
	return sdk.Result{}
}
