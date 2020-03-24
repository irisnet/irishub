package types

import (
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v3/asset/exported"
	sdk "github.com/irisnet/irishub/types"
)

type BankKeeper interface {
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	IncreaseLoosenToken(ctx sdk.Context, amt sdk.Coins)
	DecreaseLoosenToken(ctx sdk.Context, amt sdk.Coins)
}

type AssetKeeper interface {
	GetAllTokens(ctx sdk.Context) (tokens []exported.TokenI)
}

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.Account
}
