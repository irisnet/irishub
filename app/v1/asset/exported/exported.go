package exported

import (
	"github.com/irisnet/irishub/app/v1/asset/internal/types"
)

type FungibleToken = types.FungibleToken

var (
	MaximumAssetMaxSupply  = types.MaximumAssetMaxSupply
	NewFungibleToken       = types.NewFungibleToken
	NewMsgIssueToken       = types.NewMsgIssueToken
	EXTERNAL               = types.EXTERNAL
	FUNGIBLE               = types.FUNGIBLE
	CodeInvalidAssetSource = types.CodeInvalidAssetSource
	ErrAssetAlreadyExists  = types.ErrAssetAlreadyExists
)
