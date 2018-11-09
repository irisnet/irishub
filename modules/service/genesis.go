package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
	"github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/iparam"
	"github.com/irisnet/irishub/modules/service/params"
)

// GenesisState - all service state that must be provided at genesis
type GenesisState struct {
	MaxRequestTimeout  int64
	MinProviderDeposit sdk.Coins
}

func NewGenesisState(maxRequestTimeout int64, minProviderDeposit sdk.Coins) GenesisState {
	return GenesisState{
		MaxRequestTimeout:  maxRequestTimeout,
		MinProviderDeposit: minProviderDeposit,
	}
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, data GenesisState) {
	iparam.InitGenesisParameter(&serviceparams.MaxRequestTimeoutParameter, ctx, data.MaxRequestTimeout)
	iparam.InitGenesisParameter(&serviceparams.MinProviderDepositParameter, ctx, data.MinProviderDeposit)
}

// WriteGenesis - output genesis parameters
func WriteGenesis(ctx sdk.Context) GenesisState {
	maxRequestTimeout := serviceparams.GetMaxRequestTimeout(ctx)
	minProviderDeposit := serviceparams.GetMinProviderDeposit(ctx)

	return GenesisState{
		MaxRequestTimeout:  maxRequestTimeout,
		MinProviderDeposit: minProviderDeposit,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	Denom := "iris"
	IrisCt := types.NewDefaultCoinType(Denom)
	minDeposit, err := IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", 1000, Denom))
	if err != nil {
		panic(err)
	}
	return GenesisState{
		MaxRequestTimeout:  100,
		MinProviderDeposit: sdk.Coins{minDeposit},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	Denom := "iris"
	IrisCt := types.NewDefaultCoinType(Denom)
	minDeposit, err := IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", 10, Denom))
	if err != nil {
		panic(err)
	}
	return GenesisState{
		MaxRequestTimeout:  10,
		MinProviderDeposit: sdk.Coins{minDeposit},
	}
}
