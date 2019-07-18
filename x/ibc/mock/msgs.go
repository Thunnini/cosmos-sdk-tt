package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	client "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
)

type MsgRelay struct {
	Relayer sdk.AccAddress `json:"relayer"`
	Data    []byte         `json:"data"`
}

var _ sdk.Msg = &MsgRelay{}

func NewMsgRelay(relayer sdk.AccAddress, data []byte) MsgRelay {
	return MsgRelay{
		Relayer: relayer,
		Data:    data,
	}
}

//nolint
func (msg MsgRelay) Route() string { return RouterKey }
func (msg MsgRelay) Type() string  { return "ibc-relay" }
func (msg MsgRelay) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Relayer}
}

// get the bytes for the message signer to sign on
func (msg MsgRelay) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgRelay) ValidateBasic() sdk.Error {

	return nil
}

type MsgCreateClient = client.MsgCreateClient
type MsgUpdateClient = client.MsgUpdateClient
