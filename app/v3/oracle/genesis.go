package oracle

import (
	"fmt"

	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	//init feed
	for _, entry := range data.Entries {
		k.SetFeed(ctx, entry.Feed)

		reqCtx, found := k.GetRequestContext(ctx, entry.Feed.RequestContextID)
		if !found {
			panic(fmt.Errorf("no servcie request context"))
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
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// export created feed and value
	var entries []FeedEntry
	k.IteratorFeeds(ctx, func(feed types.Feed) {
		reqCtx, found := k.GetRequestContext(ctx, feed.RequestContextID)
		if found {
			entries = append(
				entries,
				FeedEntry{
					Feed:   feed,
					Values: k.GetFeedValues(ctx, feed.FeedName),
					State:  reqCtx.State,
				},
			)
		}
	})
	return GenesisState{
		Entries: entries,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Entries: []FeedEntry{},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		Entries: []FeedEntry{},
	}
}

// PrepForZeroHeightGenesis refunds the deposits, service fees and earned fees
func PrepForZeroHeightGenesis(ctx sdk.Context, k Keeper) {
	// reset request contexts state and batch
	if err := k.ResetFeedEntryState(ctx); err != nil {
		panic(fmt.Sprintf("failed to reset the feed entry state: %s", err))
	}
}
