package types

// HTLC module event types
const (
	EventTypeCreateHTLC  = "create_htlc"
	EventTypeClaimHTLC   = "claim_htlc"
	EventTypeRefundHTLC  = "refund_htlc"
	EventTypeExpiredHTLC = "expired_htlc"

	AttributeValueSender               = "sender"
	AttributeValueReceiver             = "receiver"
	AttributeValueReceiverOnOtherChain = "receiver-on-other-chain"
	AttributeValueAmount               = "amount"
	AttributeValueHashLock             = "hash-lock"
	AttributeValueSecret               = "secret"
)
