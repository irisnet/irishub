package types

// service module event types
const (
	EventTypeCreateContext           = "create_context"
	EventTypePauseContext            = "pause_context"
	EventTypeCompleteContext         = "complete_context"
	EventTypeNewBatch                = "new_batch"
	EventTypeNewBatchRequest         = "new_batch_request"
	EventTypeNewBatchRequestProvider = "new_batch_request_provider"
	EventTypeCompleteBatch           = "complete_batch"
	EventTypeServiceSlash            = "service_slash"
	EventTypeNoExchangeRate          = "no_exchange_rate"

	AttributeValueCategory          = ModuleName
	AttributeKeyAuthor              = "author"
	AttributeKeyServiceName         = "service_name"
	AttributeKeyProvider            = "provider"
	AttributeKeyConsumer            = "consumer"
	AttributeKeyRequestContextID    = "request_context_id"
	AttributeKeyRequestContextState = "request_context_state"
	AttributeKeyRequests            = "requests"
	AttributeKeyRequestID           = "request_id"
	AttributeKeyServiceFee          = "service_fee"
	AttributeKeyRequestHeight       = "request_height"
	AttributeKeyExpirationHeight    = "expiration_height"
	AttributeKeySlashedCoins        = "slashed_coins"
	AttributeKeyPriceDenom          = "price_denom"
)

type BatchState struct {
	BatchCounter           uint64                   `json:"batch_counter"`
	State                  RequestContextBatchState `json:"state"`
	BatchResponseThreshold uint32                   `json:"batch_response_threshold"`
	BatchRequestCount      uint32                   `json:"batch_request_count"`
	BatchResponseCount     uint32                   `json:"batch_response_count"`
}

// ActionTag appends action and all tagKeys
func ActionTag(action string, tagKeys ...string) string {
	tag := action
	for _, key := range tagKeys {
		tag = tag + "." + key
	}
	return tag
}
