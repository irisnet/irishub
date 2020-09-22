package types

const (
	QueryRandom             = "random" // random query endpoint supported by the random querier
	QueryRandomRequestQueue = "queue"  // random request queue query endpoint supported by the random querier
)

// QueryRandomParams is the query parameters for 'custom/random/random'
type QueryRandomParams struct {
	ReqID string `json:"req_id" yaml:"req_id"` // request id
}

// QueryRandomRequestQueueParams is the query parameters for 'custom/random/queue'
type QueryRandomRequestQueueParams struct {
	Height int64 `json:"height" yaml:"height"` // the height of the block where the random number is generated
}
