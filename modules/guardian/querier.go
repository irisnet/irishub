package guardian

import (
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryProfilers = "profilers"
	QueryTrustees  = "trustees"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryProfilers:
			return queryProfilers(ctx, k)
		case QueryTrustees:
			return queryTrustees(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown guardian query endpoint")
		}
	}
}

func queryProfilers(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	profilersIterator := k.ProfilersIterator(ctx)
	defer profilersIterator.Close()
	var profilers []Guardian
	for ; profilersIterator.Valid(); profilersIterator.Next() {
		var profiler Guardian
		k.cdc.MustUnmarshalBinaryLengthPrefixed(profilersIterator.Value(), &profiler)
		profilers = append(profilers, profiler)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, profilers)
	if err != nil {
		return nil, sdk.MarshalErr(err)
	}
	return bz, nil
}

func queryTrustees(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	trusteesIterator := k.TrusteesIterator(ctx)
	defer trusteesIterator.Close()
	var trustees []Guardian
	for ; trusteesIterator.Valid(); trusteesIterator.Next() {
		var trustee Guardian
		k.cdc.MustUnmarshalBinaryLengthPrefixed(trusteesIterator.Value(), &trustee)
		trustees = append(trustees, trustee)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, trustees)
	if err != nil {
		return nil, sdk.MarshalErr(err)
	}
	return bz, nil
}
