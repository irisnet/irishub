package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	servicetypes "github.com/irisnet/irismod/modules/service/types"
	tokenkeeper "github.com/irisnet/irismod/modules/token/keeper"
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
