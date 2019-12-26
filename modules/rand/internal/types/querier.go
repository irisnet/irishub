package types

const (
	QueryRand             = "rand"  // rand query endpoint supported by the rand querier
	QueryRandRequestQueue = "queue" // rand request queue query endpoint supported by the rand querier
)

// QueryRandParams is the query parameters for 'custom/rand/rand'
type QueryRandParams struct {
	ReqID string `json:"req_id" yaml:"req_id"` // request id
}

// QueryRandRequestQueueParams is the query parameters for 'custom/rand/queue'
type QueryRandRequestQueueParams struct {
	Height int64 `json:"height" yaml:"height"` // the height of the block where the random number is generated
}
