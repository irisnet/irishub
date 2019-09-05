package keeper

import (
	"encoding/hex"

	"github.com/irisnet/irishub/app/v2/htlc/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryHTLC:
			return queryHTLC(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown htlc query endpoint")
		}
	}
}

func queryHTLC(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryHTLCParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	secretHash, err := hex.DecodeString(params.SecretHashLock)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	htlc, err2 := keeper.GetHTLC(ctx, secretHash)
	if err2 != nil {
		return nil, err2
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, htlc)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}
