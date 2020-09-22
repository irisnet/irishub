// nolint
package types

// HTLC module event types and attributes
const (
	EventTypeCreateHTLC  = "create_htlc"
	EventTypeClaimHTLC   = "claim_htlc"
	EventTypeRefundHTLC  = "refund_htlc"
	EventTypeHTLCExpired = "htlc_expired"

	AttributeValueCategory           = ModuleName
	AttributeKeySender               = "sender"
	AttributeKeyReceiver             = "receiver"
	AttributeKeyReceiverOnOtherChain = "receiver_on_other_chain"
	AttributeKeyAmount               = "amount"
	AttributeKeyHashLock             = "hash_lock"
	AttributeKeyTimeLock             = "time_lock"
	AttributeKeySecret               = "secret"
)
