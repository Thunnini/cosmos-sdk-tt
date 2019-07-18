package exported

import "github.com/cosmos/cosmos-sdk/x/auth"

// ModuleAccountI defines an account interface for modules that hold tokens in an escrow
type ModuleAccountI interface {
	auth.Account
	GetName() string
	GetPermission() string
}
