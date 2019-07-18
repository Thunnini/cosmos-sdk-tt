package stakingmint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
)

type MsgLiquidUndelegate struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	// ValidatorAddress sdk.ValAddress `json:"validator_address"`
	// Amount           sdk.Coin       `json:"amount"`
	// DestAddress      []byte         `json:"dest"`
}

var _ sdk.Msg = &MsgLiquidUndelegate{}

func NewMsgLiquidUndelegate(delAddr sdk.AccAddress) MsgLiquidUndelegate {
	return MsgLiquidUndelegate{
		DelegatorAddress: delAddr,
	}
}

//nolint
func (msg MsgLiquidUndelegate) Route() string { return RouterKey }
func (msg MsgLiquidUndelegate) Type() string  { return "liquid-undelegate" }
func (msg MsgLiquidUndelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddress}
}

// get the bytes for the message signer to sign on
func (msg MsgLiquidUndelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgLiquidUndelegate) ValidateBasic() sdk.Error {
	if msg.DelegatorAddress.Empty() {
		return staking.ErrNilDelegatorAddr(staking.DefaultCodespace)
	}
	return nil
}
