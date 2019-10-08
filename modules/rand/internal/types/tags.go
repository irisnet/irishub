// nolint
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ActionRequestRand = []byte("request_rand")

	TagAction     = sdk.TagAction
	TagReqID      = "request-id"
	TagRandHeight = "rand-height"
	TagRand       = "rand"
)
