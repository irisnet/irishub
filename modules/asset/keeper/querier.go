package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	token "github.com/irisnet/irishub/modules/asset/01-token"
)

// NewQuerier creates a querier for the IBC module
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		var (
			res []byte
			err error
		)

		switch path[0] {
		case token.SubModuleName:
			switch path[1] {
			case token.QueryToken:
				res, err = token.QuerierToken(ctx, req, k.TokenKeeper)
			case token.QueryTokens:
				res, err = token.QuerierTokens(ctx, req, k.TokenKeeper)
			case token.QueryFees:
				res, err = token.QuerierFees(ctx, req, k.TokenKeeper)
			case token.QueryParameters:
				res, err = token.QuerierParameters(ctx, k.TokenKeeper)
			default:
				err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown Asset %s query endpoint", token.SubModuleName)
			}
		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown Asset query endpoint")
		}

		return res, sdk.ConvertError(err)
	}
}
