package stakingibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
)

// MsgDelegate - struct for bonding transactions
type MsgIBCDelegate struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Amount           sdk.Coin       `json:"amount"`
	DestAddress      []byte         `json:"dest"`
}

var _ sdk.Msg = &MsgIBCDelegate{}

func NewMsgIBCDelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin, destAddress []byte) MsgIBCDelegate {
	return MsgIBCDelegate{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
		Amount:           amount,
		DestAddress:      destAddress,
	}
}

//nolint
func (msg MsgIBCDelegate) Route() string { return RouterKey }
func (msg MsgIBCDelegate) Type() string  { return "delegate" }
func (msg MsgIBCDelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddress}
}

// get the bytes for the message signer to sign on
func (msg MsgIBCDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgIBCDelegate) ValidateBasic() sdk.Error {
	if msg.DelegatorAddress.Empty() {
		return staking.ErrNilDelegatorAddr(staking.DefaultCodespace)
	}
	if msg.ValidatorAddress.Empty() {
		return staking.ErrNilValidatorAddr(staking.DefaultCodespace)
	}
	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
		return staking.ErrBadDelegationAmount(staking.DefaultCodespace)
	}
	return nil
}

/*
// MsgUndelegate - struct for unbonding transactions
type MsgIBCUndelegate struct {
	Recipient        sdk.AccAddress `json:"recipient_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Amount           sdk.Coin       `json:"amount"`
}

var _ sdk.Msg = &MsgIBCUndelegate{}

func NewMsgIBCUndelegate(recipient sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) MsgIBCUndelegate {
	return MsgIBCUndelegate{
		Recipient:        recipient,
		ValidatorAddress: valAddr,
		Amount:           amount,
	}
}

//nolint
func (msg MsgIBCUndelegate) Route() string { return RouterKey }
func (msg MsgIBCUndelegate) Type() string  { return "begin_unbonding" }
func (msg MsgIBCUndelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Recipient}
}

// get the bytes for the message signer to sign on
func (msg MsgIBCUndelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgIBCUndelegate) ValidateBasic() sdk.Error {
	if msg.Recipient.Empty() {
		return staking.ErrNilDelegatorAddr(staking.DefaultCodespace)
	}
	if msg.ValidatorAddress.Empty() {
		return staking.ErrNilValidatorAddr(staking.DefaultCodespace)
	}
	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
		return staking.ErrBadSharesAmount(staking.DefaultCodespace)
	}
	return nil
}
*/
