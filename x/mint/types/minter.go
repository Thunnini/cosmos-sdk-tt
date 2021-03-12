package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinter returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinter(annualProvisions sdk.Dec) Minter {
	return Minter{
		AnnualProvisions: annualProvisions,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value.
func InitialMinter(inflation sdk.Dec) Minter {
	return NewMinter(sdk.NewDec(0))
}

// DefaultInitialMinter returns a default initial Minter object for a new chain
// which uses an inflation rate of 13%.
func DefaultInitialMinter() Minter {
	return InitialMinter(
		sdk.NewDecWithPrec(13, 2),
	)
}

// validate minter
func ValidateMinter(minter Minter) error {
	return nil
}

// NextAnnualProvisions returns the annual provisions based on current total
// supply and inflation rate.
func (m Minter) NextAnnualProvisions(_ Params, totalSupply sdk.Int) sdk.Dec {
	return m.AnnualProvisions.QuoInt(sdk.NewInt(2))
}

// BlockProvision returns the provisions for a block based on the annual
// provisions rate.
func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisionAmt := m.AnnualProvisions.QuoInt(sdk.NewInt(int64(params.EpochsPerYear)))
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}
