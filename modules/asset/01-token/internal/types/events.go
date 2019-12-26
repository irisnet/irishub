// nolint
package types

// asset module event types
const (
	EventTypeIssueToken    = "issue_token"
	EventTypeEditToken     = "edit_token"
	EventTypeTransferToken = "transfer_token"
	EventTypeMintToken     = "mint_token"
	EventTypeBurnToken     = "burn_token"

	AttributeKeyTokenID     = "token_id"
	AttributeKeyTokenDenom  = "token_denom"
	AttributeKeyTokenSymbol = "token_symbol"
	AttributeKeyTokenOwner  = "token_owner"

	AttributeValueCategory = SubModuleName
)
