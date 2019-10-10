package types

// asset module event types
var (
	EventTypeIssueToken           = "issue_token"
	EventTypeEditToken            = "edit_token"
	EventTypeMintToken            = "mint_token"
	EventTypeTransferTokenOwner   = "transfer_token_owner"
	EventTypeCreateGateway        = "create_gateway"
	EventTypeEditGateway          = "edit_gateway"
	EventTypeTransferGatewayOwner = "transfer_gateway_owner"

	AttributeKeyTokenId  = "toke-id"
	AttributeKeyDenom    = "denom"
	AttributeKeySymbol   = "symbol"
	AttributeKeyGateway  = "gateway"
	AttributeKeySource   = "source"
	AttributeKeyMoniker  = "moniker"
	AttributeKeyIdentity = "identity"
	AttributeKeyDetails  = "details"
	AttributeKeyWebsite  = "website"
	AttributeKeyOwner    = "owner"
	AttributeKeyPreOwner = "pre-owner"
	AttributeKeyNewOwner = "new-owner"

	AttributeValueCategory = ModuleName
)
