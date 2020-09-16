package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/guardian/types"
)

// NewQuerier creates a querier for guardian REST endpoints
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, _ abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryProfilers:
			return queryProfilers(ctx, k, legacyQuerierCdc)
		case types.QueryTrustees:
			return queryTrustees(ctx, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryProfilers(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var profilers []types.Guardian
	k.IterateProfilers(
		ctx,
		func(profiler types.Guardian) bool {
			profilers = append(profilers, profiler)
			return false
		},
	)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, profilers)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryTrustees(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var trustees []types.Guardian
	k.IterateTrustees(
		ctx,
		func(trustee types.Guardian) bool {
			trustees = append(trustees, trustee)
			return false
		},
	)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, trustees)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
