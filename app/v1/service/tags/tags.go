package tags

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	ActionSvcCallTimeOut = []byte("service-call-expiration")

	Action = sdk.TagAction

	Provider   = "provider"
	Consumer   = "consumer"
	RequestID  = "request-id"
	ServiceFee = "service-fee"
	SlashCoins = "service-slash-coins"
)
