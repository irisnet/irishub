package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/exported"
	supply "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// SupplyKeeper defines the expected supply keeper for module accounts (noalias)
type SupplyKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error

	GetSupply(ctx sdk.Context) (supply supply.SupplyI)

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) supply.ModuleAccountI

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

// AccountKeeper defines the expected account keeper for query account
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.Account
}
