package types

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	QueryRand             = "rand"
	QueryRands            = "rands"
	QueryRandRequestQueue = "queue"
)

// QueryRandParams is the query parameters for 'custom/rand/rand'
type QueryRandParams struct {
	ReqID string
}

// QueryRandsParams is the query parameters for 'custom/rand/rands'
type QueryRandsParams struct {
	Consumer sdk.AccAddress
}

// QueryRandRequestQueueParams is the query parameters for 'custom/rand/queue'
type QueryRandRequestQueueParams struct {
	Height int64
}
