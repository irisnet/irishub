package types

// Service module event types
const (
	EventTypeRequestSvc     = "request_service"
	EventTypeRespondSvc     = "respond_service"
	EventTypeSvcCallTimeout = "service_call_expiration"

	AttributeKeyProvider   = "provider"
	AttributeKeyConsumer   = "consumer"
	AttributeKeyRequestID  = "request_id"
	AttributeKeyServiceFee = "service_fee"
	AttributeKeySlashCoins = "service_slash_coins"

	AttributeValueCategory = "service"
)
