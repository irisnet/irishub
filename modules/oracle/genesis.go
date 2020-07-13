package oracle

import (
	"fmt"
	"github.com/irisnet/irishub/modules/oracle/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/oracle/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	//init feed
	for _, entry := range data.Entries {
		k.SetFeed(ctx, entry.Feed)

		reqCtx, found := k.GetRequestContext(ctx, entry.Feed.RequestContextID)
		if !found {
			panic(fmt.Errorf("unknown servcie request context: %s", entry.Feed.RequestContextID))
		}

		for _, value := range entry.Values {
			k.SetFeedValue(
				ctx,
				entry.Feed.FeedName,
				reqCtx.BatchCounter,
				entry.Feed.LatestHistory,
				value,
			)
		}

		k.Enqueue(ctx, entry.Feed.FeedName, entry.State)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	// export created feed and value
	var entries []types.FeedEntry
	k.IteratorFeeds(ctx, func(feed types.Feed) {
		reqCtx, found := k.GetRequestContext(ctx, feed.RequestContextID)
		if found {
			entries = append(
				entries,
				types.FeedEntry{
					Feed:   feed,
					Values: k.GetFeedValues(ctx, feed.FeedName),
					State:  reqCtx.State,
				},
			)
		}
	})
	return types.GenesisState{
		Entries: entries,
	}
}

// PrepForZeroHeightGenesis refunds the deposits, service fees and earned fees
func PrepForZeroHeightGenesis(ctx sdk.Context, k keeper.Keeper) {
	// reset request contexts state and batch
	if err := k.ResetFeedEntryState(ctx); err != nil {
		panic(fmt.Sprintf("failed to reset the feed entry state: %s", err))
	}
}
