package stakingibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the staking module
	ModuleName = "staking-ibc"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName
)

type PacketIBCDelegate struct {
	From       sdk.AccAddress `json:"from"`
	DstAddress []byte         `json:"dest"` // Recipient address on counterparty chain
	Validator  sdk.ValAddress `json:"validator"`
	Amount     sdk.Coin       `json:"amount"`
}

type delegateInfo struct {
	Delegator   sdk.AccAddress
	Validator   sdk.ValAddress
	DestAddress []byte
}
