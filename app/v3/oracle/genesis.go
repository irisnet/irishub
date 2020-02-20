package oracle

import (
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
		if reqCtx, found := k.GetRequestContext(ctx, entry.Feed.RequestContextID); found {
			for _, value := range entry.Values {
				k.SetFeedValue(ctx,
					entry.Feed.FeedName,
					reqCtx.BatchCounter,
					entry.Feed.LatestHistory,
					value)
			}
			k.Enqueue(ctx, entry.Feed.FeedName, entry.State)
		}
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// export created feed and value
	var entries []FeedEntry
	k.IteratorFeeds(ctx, func(feed types.Feed) {
		reqCtx, found := k.GetRequestContext(ctx, feed.RequestContextID)
		if found {
			entries = append(entries, FeedEntry{
				Feed:   feed,
				Values: k.GetFeedValues(ctx, feed.FeedName),
				State:  reqCtx.State,
			})
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

// ValidateGenesis validates the provided asset genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	for _, entry := range data.Entries {
		feed := entry.Feed
		if err := types.ValidateFeedName(feed.FeedName); err != nil {
			panic(err)
		}
		if err := types.ValidateDescription(feed.Description); err != nil {
			panic(err)
		}

		if err := types.ValidateAggregateFunc(feed.AggregateFunc); err != nil {
			return err
		}

		if err := types.ValidateValueJsonPath(feed.ValueJsonPath); err != nil {
			return err
		}

		if err := types.ValidateLatestHistory(feed.LatestHistory); err != nil {
			return err
		}

		if err := types.ValidateCreator(feed.Creator); err != nil {
			return err
		}
	}
	return nil
}
