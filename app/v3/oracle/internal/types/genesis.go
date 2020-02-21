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
