package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ActionSvcDef           = []byte("service-define")
	ActionSvcBind          = []byte("service-bind")
	ActionSvcBindUpdate    = []byte("service-update-binding")
	ActionSvcRefundDeposit = []byte("service-refund-deposit")

	Action = sdk.TagAction
)
