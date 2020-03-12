package types

import sdk "github.com/irisnet/irishub/types"

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)

	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)

	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SetTotalSupply(ctx sdk.Context, totalSupply sdk.Coin)

	GetTotalSupply(ctx sdk.Context, denom string) (coin sdk.Coin, found bool)

	IncreaseTotalSupply(ctx sdk.Context, amt sdk.Coin) sdk.Error

	BurnCoins(ctx sdk.Context, fromAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)

	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
}

// GuardianKeeper defines the expected guardian keeper
type GuardianKeeper interface {
	ProfilersIterator(ctx sdk.Context) sdk.Iterator
}
