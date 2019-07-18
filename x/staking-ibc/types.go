package stakingibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/mock"
)

const (
	// ModuleName is the name of the staking module
	ModuleName = "stakingibc"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName
)

type PacketIBCDelegated struct {
	From       sdk.AccAddress `json:"from"`
	DstAddress []byte         `json:"dest"` // Recipient address on counterparty chain
	Validator  sdk.ValAddress `json:"validator"`
	Amount     sdk.Coin       `json:"amount"`
}

var _ mock.Packet = PacketIBCDelegated{}

func (PacketIBCDelegated) Type() string {
	return "delegated"
}

type PacketIBCUndelegate struct {
	From      sdk.AccAddress `json:"from"`
	Validator sdk.ValAddress `json:"validator"`
	Amount    sdk.Coin       `json:"amount"`
}

var _ mock.Packet = PacketIBCUndelegate{}

func (PacketIBCUndelegate) Type() string {
	return "undelegate"
}

type delegateInfo struct {
	Delegator   sdk.AccAddress
	Validator   sdk.ValAddress
	DestAddress []byte
}
