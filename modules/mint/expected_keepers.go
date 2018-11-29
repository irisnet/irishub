package mint

import (
	sdk "github.com/irisnet/irishub/types"
)

// expected stake keeper
type StakeKeeper interface {
	TotalPower(ctx sdk.Context) sdk.Dec
	BondedRatio(ctx sdk.Context) sdk.Dec
	InflateSupply(ctx sdk.Context, newTokens sdk.Dec)
}

// expected fee collection keeper interface
type FeeCollectionKeeper interface {
	AddCollectedFees(sdk.Context, sdk.Coins) sdk.Coins
}
