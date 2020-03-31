package types

import (
	service "github.com/irisnet/irishub/app/v3/service/exported"
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Entries []FeedEntry `json:"entries"`
}

type FeedEntry struct {
	Feed   Feed                        `json:"feed"`
	State  service.RequestContextState `json:"state"`
	Values []FeedValue                 `json:"values"`
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
