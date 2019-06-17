// nolint
package tags

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	ActionIssueToken = []byte("issue-token")

	Action  = sdk.TagAction
	Id      = "token-id"
	Denom   = "token-denom"
	Symbol  = "token-symbol"
	Owner   = "token-owner"
	Gateway = "token-gateway"
	Source  = "token-source"
)
