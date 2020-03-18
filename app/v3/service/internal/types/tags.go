package types

// nolint
var (
	ActionCreateContext   = "create-context"
	ActionPauseContext    = "pause-context"
	ActionCompleteContext = "complete-context"
	ActionNewBatch        = "new-batch"
	ActionNewBatchRequest = "new-batch-request"

	TagAuthor           = "author"
	TagServiceName      = "service-name"
	TagProvider         = "provider"
	TagConsumer         = "consumer"
	TagRequestContextID = "request-context-id"
	TagRequestID        = "request-id"
	TagServiceFee       = "service-fee"
	TagRequestHeight    = "request-height"
	TagExpirationHeight = "expiration-height"
	TagSlashedCoins     = "slashed-coins"
)
