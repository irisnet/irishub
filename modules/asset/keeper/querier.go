package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	token "github.com/irisnet/irishub/modules/asset/01-token"
)

// NewQuerier creates a querier for the asset module
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case token.SubModuleName:
			switch path[1] {
			case token.QueryToken:
				return token.QuerierToken(ctx, req, k.TokenKeeper)
			case token.QueryTokens:
				return token.QuerierTokens(ctx, req, k.TokenKeeper)
			case token.QueryFees:
				return token.QuerierFees(ctx, req, k.TokenKeeper)
			case token.QueryParameters:
				return token.QuerierParameters(ctx, k.TokenKeeper)
			default:
				return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown sub query path: %s", path[1])
			}
		}
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
	}
}
