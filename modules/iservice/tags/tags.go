package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ActionSvcDef        = []byte("service-define")
	ActionSvcBind       = []byte("service-bind")
	ActionSvcBindUpdate = []byte("service-binding-update")

	Action = sdk.TagAction
)
