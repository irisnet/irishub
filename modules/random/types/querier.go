package types

const (
	QueryRandom             = "rand"  // rand query endpoint supported by the rand querier
	QueryRandomRequestQueue = "queue" // rand request queue query endpoint supported by the rand querier
)

// QueryRandomParams is the query parameters for 'custom/rand/rand'
type QueryRandomParams struct {
	ReqID string `json:"req_id" yaml:"req_id"` // request id
}

// QueryRandomRequestQueueParams is the query parameters for 'custom/rand/queue'
type QueryRandomRequestQueueParams struct {
	Height int64 `json:"height" yaml:"height"` // the height of the block where the random number is generated
}
