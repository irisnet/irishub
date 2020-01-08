package exported

import (
	"github.com/irisnet/irishub/app/v3/asset/internal/types"
)

type FungibleToken = types.FungibleToken

var (
	MaximumAssetMaxSupply = types.MaximumAssetMaxSupply
	NewFungibleToken      = types.NewFungibleToken
	EXTERNAL              = types.EXTERNAL
	FUNGIBLE              = types.FUNGIBLE
	NewMsgIssueToken      = types.NewMsgIssueToken
	ValidateMsgIssueToken = types.ValidateMsgIssueToken
	ErrAssetAlreadyExists = types.ErrAssetAlreadyExists
)
