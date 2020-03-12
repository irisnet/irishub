// nolint
package types

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	ActionRequestRand = []byte("request_rand")

	TagAction           = sdk.TagAction
	TagReqID            = "request-id"
	TagRequestContextID = "request-context-id"
	TagRandHeight       = "rand-height"
	TagRand             = "rand"
)
