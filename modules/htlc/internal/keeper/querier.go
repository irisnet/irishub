package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/htlc/internal/types"
)

// NewQuerier creates a new HTLC Querier instance
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryHTLC:
			return queryHTLC(ctx, req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryHTLC(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryHTLCParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if len(params.HashLock) != types.HashLockLength {
		return nil, sdkerrors.Wrapf(types.ErrInvalidHashLock, "the hash lock must be %d bytes long", types.HashLockLength)
	}

	htlc, err := keeper.GetHTLC(ctx, params.HashLock)
	if err != nil {
		return nil, err
	}

	bz, errRes := codec.MarshalJSONIndent(keeper.cdc, htlc)
	if errRes != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
