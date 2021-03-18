// nolint
package types

// HTLC module event types and attributes
const (
	EventTypeCreateHTLC = "create_htlc"
	EventTypeClaimHTLC  = "claim_htlc"
	EventTypeRefundHTLC = "refund_htlc"

	AttributeValueCategory = ModuleName

	AttributeKeySender               = "sender"
	AttributeKeyReceiver             = "receiver"
	AttributeKeyReceiverOnOtherChain = "receiver_on_other_chain"
	AttributeKeySenderOnOtherChain   = "sender_on_other_chain"
	AttributeKeyAmount               = "amount"
	AttributeKeyHashLock             = "hash_lock"
	AttributeKeyID                   = "id"
	AttributeKeyTimeLock             = "time_lock"
	AttributeKeySecret               = "secret"
	AttributeKeyTransfer             = "transfer"
	AttributeKeyDirection            = "direction"
)
