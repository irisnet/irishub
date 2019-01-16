package mint

import (
	sdk "github.com/irisnet/irishub/types"
)

// expected fee collection keeper interface
type FeeKeeper interface {
	AddCollectedFees(sdk.Context, sdk.Coins) sdk.Coins
}
