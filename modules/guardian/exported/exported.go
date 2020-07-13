package exported

import (
	"github.com/irisnet/irishub/modules/guardian/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GuardianI expected guardian functions
type GuardianI interface {
	GetDescription() string
	GetAccountType() types.AccountType
	GetAddress() sdk.AccAddress
	GetAddedBy() sdk.AccAddress
}
