package exported

import (
	"github.com/irisnet/irishub/modules/asset/internal/types"
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
