package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//expected bank keeper
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
}
