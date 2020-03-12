package oracle

import (
	"time"

	cmn "github.com/tendermint/tendermint/libs/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/migrate/v0_17/service"
)

type GenesisState struct {
	Entries []FeedEntry `json:"entries"`
}

type FeedEntry struct {
	Feed   Feed                        `json:"feed"`
	State  service.RequestContextState `json:"state"`
	Values []FeedValue                 `json:"values"`
}

type Feed struct {
	FeedName         string         `json:"feed_name"`
	Description      string         `json:"description"`
	AggregateFunc    string         `json:"aggregate_func"`
	ValueJsonPath    string         `json:"value_json_path"`
	LatestHistory    uint64         `json:"latest_history"`
	RequestContextID cmn.HexBytes   `json:"request_context_id"`
	Creator          sdk.AccAddress `json:"creator"`
}

type FeedValue struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
