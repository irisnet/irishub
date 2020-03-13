package types

// nolint
var (
	ActionNewBatch = "new-batch"

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

// BatchState defines the state of a batch, used in enblock
type BatchState struct {
	BatchCounter       uint64                   `json:"batch_counter"`
	State              RequestContextBatchState `json:"state"`
	ResponseThreshold  uint16                   `json:"response_threshold"`
	BatchRequestCount  uint16                   `json:"batch_request_count"`
	BatchResponseCount uint16                   `json:"batch_response_count"`
}

// ActionTag returns action.tagsKeys
func ActionTag(action string, tagKeys ...string) string {
	tag := action
	for _, key := range tagKeys {
		tag = tag + "." + key
	}
	return tag
}
