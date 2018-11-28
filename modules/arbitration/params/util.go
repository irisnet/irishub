package arbitrationparams

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func GetComplaintRetrospect(ctx sdk.Context) time.Duration {
	ComplaintRetrospectParameter.LoadValue(ctx)
	return ComplaintRetrospectParameter.Value
}

func GetArbitrationTimelimit(ctx sdk.Context) time.Duration {
	ArbitrationTimelimitParameter.LoadValue(ctx)
	return ArbitrationTimelimitParameter.Value
}