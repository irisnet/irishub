package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/guardian/internal/types"
)

// NewQuerier creates a querier for guardian REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, _ abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryProfilers:
			return queryProfilers(ctx, k)
		case types.QueryTrustees:
			return queryTrustees(ctx, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryProfilers(ctx sdk.Context, k Keeper) ([]byte, error) {
	profilersIterator := k.ProfilersIterator(ctx)
	defer profilersIterator.Close()
	var profilers []types.Guardian
	for ; profilersIterator.Valid(); profilersIterator.Next() {
		var profiler types.Guardian
		k.cdc.MustUnmarshalBinaryLengthPrefixed(profilersIterator.Value(), &profiler)
		profilers = append(profilers, profiler)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, profilers)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryTrustees(ctx sdk.Context, k Keeper) ([]byte, error) {
	trusteesIterator := k.TrusteesIterator(ctx)
	defer trusteesIterator.Close()
	var trustees []types.Guardian
	for ; trusteesIterator.Valid(); trusteesIterator.Next() {
		var trustee types.Guardian
		k.cdc.MustUnmarshalBinaryLengthPrefixed(trusteesIterator.Value(), &trustee)
		trustees = append(trustees, trustee)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, trustees)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
