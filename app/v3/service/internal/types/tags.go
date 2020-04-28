package types

// nolint
var (
	ActionCreateContext   = "create-context"
	ActionPauseContext    = "pause-context"
	ActionCompleteContext = "complete-context"
	ActionNewBatch        = "new-batch"
	ActionNewBatchRequest = "new-batch-request"
	ActionCompleteBatch   = "complete-batch"

	TagAuthor           = "author"
	TagServiceName      = "service-name"
	TagOwner            = "owner"
	TagProvider         = "provider"
	TagConsumer         = "consumer"
	TagRequestContextID = "request-context-id"
	TagRequestID        = "request-id"
	TagServiceFee       = "service-fee"
	TagRequestHeight    = "request-height"
	TagExpirationHeight = "expiration-height"
	TagSlashedCoins     = "slashed-coins"
)

type BatchState struct {
	BatchCounter           uint64                   `json:"batch_counter"`
	State                  RequestContextBatchState `json:"state"`
	BatchResponseThreshold uint16                   `json:"batch_response_threshold"`
	BatchRequestCount      uint16                   `json:"batch_request_count"`
	BatchResponseCount     uint16                   `json:"batch_response_count"`
}
