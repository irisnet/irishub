package types

// service module event types
const (
	EventTypeCreateDefinition        = "create_definition"
	EventTypeCreateBinding           = "create_binding"
	EventTypeUpdateBinding           = "update_binding"
	EventTypeDisableBinding          = "disable_binding"
	EventTypeEnableBinding           = "enable_binding"
	EventTypeRefundDeposit           = "refund_deposit"
	EventTypeSetWithdrawAddress      = "set_withdraw_address"
	EventTypeRespondService          = "respond_service"
	EventTypeCreateContext           = "create_context"
	EventTypePauseContext            = "pause_context"
	EventTypeStartContext            = "start_context"
	EventTypeKillContext             = "kill_context"
	EventTypeUpdateContext           = "update_context"
	EventTypeWithdrawEarnedFees      = "withdraw_earned_fees"
	EventTypeCompleteContext         = "complete_context"
	EventTypeNewBatch                = "new_batch"
	EventTypeNewBatchRequest         = "new_batch_request"
	EventTypeNewBatchRequestProvider = "new_batch_request_provider"
	EventTypeCompleteBatch           = "complete_batch"
	EventTypeServiceSlash            = "service_slash"
	EventTypeNoExchangeRate          = "no_exchange_rate"

	AttributeValueCategory = ModuleName

	AttributeKeyAuthor              = "author"
	AttributeKeyServiceName         = "service_name"
	AttributeKeyProvider            = "provider"
	AttributeKeyOwner               = "owner"
	AttributeKeyWithdrawAddress     = "withdraw_address"
	AttributeKeyConsumer            = "consumer"
	AttributeKeyModuleService       = "module_service"
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
