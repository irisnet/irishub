package keeper

import "github.com/cosmos/cosmos-sdk/types/query"

var (
	paginationDefaultLimit uint64 = 100
	paginationMaxLimit     uint64 = 100
)

// shapePageRequest shapes the PageRequest params to avoid querying all items.
// PageRequest.offset is forbidden and PageRequest.count_total must be zero.
// PageRequest.limit mustn't exceed paginationMaxLimit and is set to
// paginationDefaultLimit when unset.
func shapePageRequest(req *query.PageRequest) *query.PageRequest {
	res := newDefaultPageRequest()

	if req == nil {
		return res
	}

	res.Key = req.Key
	res.Reverse = req.Reverse
	if req.Limit > 0 && req.Limit <= paginationMaxLimit {
		res.Limit = req.Limit
	}

	return res
}

// newDefaultPageRequest returns a default PageRequest.
func newDefaultPageRequest() *query.PageRequest {
	return &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      paginationDefaultLimit,
		CountTotal: false,
		Reverse:    false,
	}
}
