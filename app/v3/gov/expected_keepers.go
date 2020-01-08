package gov

import (
	"github.com/irisnet/irishub/app/v3/asset/exported"
	sdk "github.com/irisnet/irishub/types"
)

// AssetKeeper expected asset keeper
type AssetKeeper interface {
	IssueToken(ctx sdk.Context, token exported.FungibleToken) (sdk.Tags, sdk.Error)
	HasToken(ctx sdk.Context, tokenId string) bool
}
