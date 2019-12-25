// nolint
package types

// HTLC module event types
const (
	EventTypeCreateHTLC  = "create_htlc"
	EventTypeClaimHTLC   = "claim_htlc"
	EventTypeRefundHTLC  = "refund_htlc"
	EventTypeExpiredHTLC = "expired_htlc"

	AttributeValueSender               = "sender"
	AttributeValueReceiver             = "receiver"
	AttributeValueReceiverOnOtherChain = "receiver_on_other_chain"
	AttributeValueAmount               = "amount"
	AttributeValueHashLock             = "hash_lock"
	AttributeValueSecret               = "secret"

	AttributeValueCategory = ModuleName
)
