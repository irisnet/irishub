// nolint
package types

const (
	EventTypeIssueToken         = "issue_token"
	EventTypeEditToken          = "edit_token"
	EventTypeMintToken          = "mint_token"
	EventTypeBurnToken          = "burn_token"
	EventTypeTransferTokenOwner = "transfer_token_owner"
	EventTypeSwapFeeToken       = "swap_fee_token"

	AttributeValueCategory = ModuleName

	AttributeKeyCreator   = "creator"
	AttributeKeySymbol    = "symbol"
	AttributeKeyAmount    = "amount"
	AttributeKeyOwner     = "owner"
	AttributeKeyDstOwner  = "dst_owner"
	AttributeKeyRecipient = "recipient"
	AttributeKeySender    = "sender"
	AttributeKeyFeePaid   = "fee_paid"
	AttributeKeyFeeGot    = "fee_git"
)
