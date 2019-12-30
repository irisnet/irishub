package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/irisnet/irishub/modules/guardian/exported"
)

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

// GuardianKeeper defines the expected guardian keeper (noalias)
type GuardianKeeper interface {
	GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (guardian exported.GuardianI, found bool)
	GetTrustee(ctx sdk.Context, addr sdk.AccAddress) (guardian exported.GuardianI, found bool)
}

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) supplyexported.ModuleAccountI
	GetModuleAddress(moduleName string) sdk.AccAddress

	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error

	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}
