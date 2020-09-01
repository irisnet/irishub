package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	servicetypes "github.com/irismod/service/types"
	tokenkeeper "github.com/irismod/token/keeper"
)

var (
	_ servicetypes.TokenKeeper = tokenAdapter{}
)

func WrapToken(tk tokenkeeper.Keeper) tokenAdapter {
	return tokenAdapter{tk}
}

type tokenAdapter struct {
	tk tokenkeeper.Keeper
}

func (t tokenAdapter) GetToken(ctx sdk.Context, denom string) (servicetypes.TokenI, error) {
	return t.tk.GetToken(ctx, denom)
}
