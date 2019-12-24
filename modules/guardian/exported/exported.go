package exported

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountTypeI
type AccountTypeI interface {
	Format(s fmt.State, verb rune)
	String() string
	MarshalJSON() ([]byte, error)
}

// GuardianI
type GuardianI interface {
	GetDescription() string
	GetAccountType() AccountTypeI
	GetAddress() sdk.AccAddress
	GetAddedBy() sdk.AccAddress
}
