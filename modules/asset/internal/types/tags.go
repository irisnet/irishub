// nolint
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	TagActionIssueToken = []byte("issue_token")

	TagAction  = sdk.TagAction
	TagId      = "token-id"
	TagDenom   = "token-denom"
	TagSymbol  = "token-symbol"
	TagOwner   = "token-owner"
	TagGateway = "token-gateway"
	TagSource  = "token-source"
)
