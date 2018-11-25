package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ActionSvcDef           = []byte("service-define")
	ActionSvcBind          = []byte("service-bind")
	ActionSvcBindUpdate    = []byte("service-update-binding")
	ActionSvcRefundDeposit = []byte("service-refund-deposit")
	ActionSvcDisable       = []byte("service-disable")
	ActionSvcEnable        = []byte("service-enable")

	ActionSvcCall         = []byte("service-call")
	ActionSvcRespond      = []byte("service-respond")
	ActionSvcRefundFees   = []byte("service-refund-fees")
	ActionSvcWithdrawFees = []byte("service-withdraw-fees")

	ActionSvcCallTimeOut = []byte("service-call-expiration")

	Action = sdk.TagAction

	Provider = "provider"
	Consumer = "consumer"
	RequestID = "request-id"
)
