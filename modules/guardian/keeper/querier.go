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
		case types.QuerySupers:
			return querySupers(ctx, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func querySupers(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var supers []types.Super
	k.IterateSupers(
		ctx,
		func(super types.Super) bool {
			supers = append(supers, super)
			return false
		},
	)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, supers)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
