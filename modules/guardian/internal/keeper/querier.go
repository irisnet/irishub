package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/modules/guardian/internal/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryProfilers:
			return queryProfilers(ctx, k)
		case types.QueryTrustees:
			return queryTrustees(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown guardian query endpoint")
		}
	}
}

func queryProfilers(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
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
		return nil, sdk.ConvertError(err)
	}
	return bz, nil
}

func queryTrustees(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
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
		return nil, sdk.ConvertError(err)
	}
	return bz, nil
}
