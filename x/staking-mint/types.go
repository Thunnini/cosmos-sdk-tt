package stakingmint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the staking module
	ModuleName = "liquid-bond"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName
)

type liquidDelegateInfo struct {
	Delegator   sdk.AccAddress
	Validator   sdk.ValAddress
	DestAddress sdk.AccAddress
	Amount      sdk.Coin
}
