package types

// NFT module event types
var (
	EventTypeIssueDenom    = "issue_denom"
	EventTypeTransfer      = "transfer_nft"
	EventTypeEditNFT       = "edit_nft"
	EventTypeMintNFT       = "mint_nft"
	EventTypeBurnNFT       = "burn_nft"
	EventTypeTransferDenom = "transfer_denom"

	AttributeValueCategory = ModuleName

	AttributeKeySender    = "sender"
	AttributeKeyCreator   = "creator"
	AttributeKeyRecipient = "recipient"
	AttributeKeyOwner     = "owner"
	AttributeKeyTokenID   = "token_id"
	AttributeKeyTokenURI  = "token_uri"
	AttributeKeyDenomID   = "denom_id"
	AttributeKeyDenomName = "denom_name"
)
