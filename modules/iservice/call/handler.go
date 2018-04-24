package call

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SvcCallService struct {
	keeper SvcCallKeeper
}

func NewSvcCallService(keeper SvcCallKeeper) SvcCallService {
	return SvcCallService{
		keeper,
	}
}

func (service SvcCallService) CheckTx(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	return sdk.Result{}
}

func (service SvcCallService) DeliverTx(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	return sdk.Result{}
}
