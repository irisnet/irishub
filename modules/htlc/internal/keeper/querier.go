package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/htlc/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryHTLC:
			return queryHTLC(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown HTLC query endpoint")
		}
	}
}

func queryHTLC(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryHTLCParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	if len(params.HashLock) != types.HashLockLength {
		return nil, types.ErrInvalidHashLock(types.DefaultCodespace, fmt.Sprintf("the hash lock must be %d bytes long", types.HashLockLength))
	}

	htlc, err2 := keeper.GetHTLC(ctx, params.HashLock)
	if err2 != nil {
		return nil, err2
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, htlc)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}
