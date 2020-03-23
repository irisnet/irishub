// nolint
package types

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
)

var (
	ActionRequestRand = []byte("request_rand")

	TagAction           = sdk.TagAction
	TagReqID            = "request-id"
	TagRequestContextID = "request-context-id"
	TagRandHeight       = "rand-height"
	TagRand             = func(reqID string) string {
		return fmt.Sprintf("rand.%s", reqID)
	}
)
