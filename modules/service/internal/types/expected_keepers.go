package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/guardian"
)

// BankKeeper defines the expected bank keeper interface
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
}

// GuardianKeeper defines the expected guardian keeper interface
type GuardianKeeper interface {
	GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (guardian guardian.Guardian, found bool)
	GetTrustee(ctx sdk.Context, addr sdk.AccAddress) (guardian guardian.Guardian, found bool)
}
