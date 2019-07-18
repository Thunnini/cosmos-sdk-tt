package stakingibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
)

func (k StakingIBCKeeper) AfterValidatorBonded(ctx sdk.Context, address sdk.ConsAddress, _ sdk.ValAddress) {

}

//_________________________________________________________________________________________

// Hooks wrapper struct for slashing StakingIBCKeeper
type Hooks struct {
	k StakingIBCKeeper
}

var _ types.StakingHooks = Hooks{}

// TODO: hook after unbonding delegation, and reward to delegator from module account.

// Return the wrapper struct
func (k StakingIBCKeeper) Hooks() Hooks {
	return Hooks{k}
}

func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}

func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, _ sdk.ValAddress) {
}

func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
}

// nolint - unused hooks
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)  {}
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                          {}
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) AfterDelegationModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ sdk.ValAddress, _ sdk.Dec)                {}
