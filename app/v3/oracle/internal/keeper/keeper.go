package keeper

import (
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// Keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	// codespace
	codespace sdk.CodespaceType
	// params subspace
	paramSpace params.Subspace
}

// NewKeeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType, paramSpace params.Subspace) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		codespace:  codespace,
		paramSpace: paramSpace.WithTypeTable(types.ParamTypeTable()),
	}

	return keeper
}

func (k Keeper) CreateFeed(ctx sdk.Context, msg types.MsgCreateFeed) (sdk.Tags, sdk.Error) {
	//TODO
	return nil, nil
}

func (k Keeper) StartFeed(ctx sdk.Context, msg types.MsgStartFeed) (sdk.Tags, sdk.Error) {
	//TODO
	return nil, nil
}

func (k Keeper) StopFeed(ctx sdk.Context, msg types.MsgStopFeed) (sdk.Tags, sdk.Error) {
	//TODO
	return nil, nil
}

func (k Keeper) EditFeed(ctx sdk.Context, msg types.MsgEditFeed) (sdk.Tags, sdk.Error) {
	//TODO
	return nil, nil
}
