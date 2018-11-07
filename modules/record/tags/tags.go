// nolint
package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ActionSubmitRecord = []byte("submit-record")

	Action       = sdk.TagAction
	OwnerAddress = "ownerAddress"
	RecordID     = "record-id"
)
