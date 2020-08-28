package types

// get raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Entries: []FeedEntry{},
	}
}

// ValidateGenesis validates the provided asset genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	for _, entry := range data.Entries {
		feed := entry.Feed
		if err := ValidateFeedName(feed.FeedName); err != nil {
			return err
		}
		if err := ValidateDescription(feed.Description); err != nil {
			return err
		}
		if err := ValidateAggregateFunc(feed.AggregateFunc); err != nil {
			return err
		}
		if err := ValidateValueJsonPath(feed.ValueJsonPath); err != nil {
			return err
		}
		if err := ValidateLatestHistory(feed.LatestHistory); err != nil {
			return err
		}
		if err := ValidateCreator(feed.Creator); err != nil {
			return err
		}
	}
	return nil
}
