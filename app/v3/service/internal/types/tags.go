package types

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	TagActionSvcCallTimeOut = []byte("service_call_expiration")

	TagAction = sdk.TagAction

	TagProvider   = "provider"
	TagConsumer   = "consumer"
	TagRequestID  = "request-id"
	TagServiceFee = "service-fee"
	TagSlashCoins = "service-slash-coins"
)
