package bind

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"log"
)

type SvcBindService struct {
	keeper SvcBindKeeper
}

func NewSvcBindService(keeper SvcBindKeeper) SvcBindService {
	return SvcBindService{
		keeper,
	}
}

func (service SvcBindService) CheckTx(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	log.Printf("enter into %s", "SvcBindService.CheckTx")
	return sdk.Result{}
}

func (service SvcBindService) DeliverTx(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	log.Printf("enter into %s", "SvcBindService.DeliverTx")
	return sdk.Result{}
}
