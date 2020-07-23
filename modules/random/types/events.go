// nolint
package types

// rand module event types
const (
	EventTypeRequestRandom  = "request_rand"
	EventTypeGenerateRandom = "generate_rand"

	AttributeKeyRequestID        = "request_id"
	AttributeKeyGenHeight        = "generate_height"
	AttributeKeyRandom           = "rand"
	AttributeKeyRequestContextID = "request_context_id"

	AttributeValueCategory = ModuleName
)
