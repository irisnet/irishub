package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

const (
	FeedName  = "feed-name"
	FeedState = "state"
)

// RegisterHandlers registers oracle REST handlers to a router
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type createFeedReq struct {
	BaseReq           rest.BaseReq `json:"base_req" yaml:"base_req"`
	FeedName          string       `json:"feed_name"`
	AggregateFunc     string       `json:"aggregate_func"`
	ValueJsonPath     string       `json:"value_json_path"`
	LatestHistory     uint64       `json:"latest_history"`
	Description       string       `json:"description"`
	Creator           string       `json:"creator"`
	ServiceName       string       `json:"service_name"`
	Providers         []string     `json:"providers"`
	Input             string       `json:"input"`
	Timeout           int64        `json:"timeout"`
	ServiceFeeCap     string       `json:"service_fee_cap"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	ResponseThreshold uint32       `json:"response_threshold"`
}

type editFeedReq struct {
	BaseReq           rest.BaseReq `json:"base_req" yaml:"base_req"`
	Description       string       `json:"description"`
	LatestHistory     uint64       `json:"latest_history"`
	Creator           string       `json:"creator"`
	Providers         []string     `json:"providers"`
	Timeout           int64        `json:"timeout"`
	ServiceFeeCap     string       `json:"service_fee_cap"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	ResponseThreshold uint32       `json:"response_threshold"`
}

type startFeedReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Creator string       `json:"creator"`
}

type pauseFeedReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Creator string       `json:"creator"`
}
