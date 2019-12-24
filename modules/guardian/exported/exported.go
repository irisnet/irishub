package exported

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountTypeI expected account type functions
type AccountTypeI interface {
	Format(s fmt.State, verb rune)
	String() string
	MarshalJSON() ([]byte, error)
}

// GuardianI expected guardian functions
type GuardianI interface {
	GetDescription() string
	GetAccountType() AccountTypeI
	GetAddress() sdk.AccAddress
	GetAddedBy() sdk.AccAddress
}
