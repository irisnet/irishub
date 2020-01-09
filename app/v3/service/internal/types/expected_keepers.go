package types

import (
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
)

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
	BurnCoins(ctx sdk.Context, fromAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
}

// GuardianKeeper defines the expected guardian keeper (noalias)
type GuardianKeeper interface {
	GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (guardian guardian.Guardian, found bool)
	GetTrustee(ctx sdk.Context, addr sdk.AccAddress) (guardian guardian.Guardian, found bool)
}
