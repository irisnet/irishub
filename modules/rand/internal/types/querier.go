package types

const (
	QueryRand             = "rand"  // rand query endpoint supported by the rand querier
	QueryRandRequestQueue = "queue" // rand request queue query endpoint supported by the rand querier
)

// QueryRandParams is the query parameters for 'custom/rand/rand'
type QueryRandParams struct {
	ReqID string
}

// QueryRandRequestQueueParams is the query parameters for 'custom/rand/queue'
type QueryRandRequestQueueParams struct {
	Height int64
}
