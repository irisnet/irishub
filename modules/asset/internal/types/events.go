// nolint
package types

const (
	EventTypeIssueToken         = "issue_token"
	EventTypeEditToken          = "edit_token"
	EventTypeTransferTokenOwner = "transfer_token_owner"
	EventTypeMintToken          = "mint_token"

	AttributeKeyTokenID     = "token_id"
	AttributeKeyTokenDenom  = "token_denom"
	AttributeKeyTokenSymbol = "token_symbol"
	AttributeKeyTokenOwner  = "token_owner"
	AttributeKeyTokenSource = "token_source"

	AttributeValueCategory = ModuleName
)
