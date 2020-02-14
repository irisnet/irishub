package keeper

import (
	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// Keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	sk types.ServiceKeeper
	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType, sk types.ServiceKeeper) Keeper {
	keeper := Keeper{
		storeKey:  key,
		cdc:       cdc,
		sk:        sk,
		codespace: codespace,
	}

	return keeper
}

func (k Keeper) CreateFeed(ctx sdk.Context, msg types.MsgCreateFeed) sdk.Error {
	if k.hasFeed(ctx, msg.FeedName) {
		return types.ErrExistedFeedName(types.DefaultCodespace, msg.FeedName)
	}
	requestContextID, err := k.sk.CreateRequestContext(ctx)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	feed := Feed{
		FeedName:              msg.FeedName,
		AggregateMethod:       msg.AggregateMethod,
		AggregateArgsJsonPath: msg.AggregateArgsJsonPath,
		LatestHistory:         msg.LatestHistory,
		RequestContextID:      requestContextID,
		Owner:                 msg.Owner,
	}
	k.setFeed(ctx, feed)
	return nil
}

func (k Keeper) StartFeed(ctx sdk.Context, msg types.MsgStartFeed) sdk.Error {
	//TODO
	return nil
}

func (k Keeper) PauseFeed(ctx sdk.Context, msg types.MsgPauseFeed) sdk.Error {
	//TODO
	return nil
}

func (k Keeper) KillFeed(ctx sdk.Context, msg types.MsgKillFeed) sdk.Error {
	//TODO
	return nil
}

func (k Keeper) EditFeed(ctx sdk.Context, msg types.MsgEditFeed) sdk.Error {
	//TODO
	return nil
}
