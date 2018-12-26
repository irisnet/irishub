package service

import (
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/service/params"
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all service state that must be provided at genesis
type GenesisState struct {
	ServiceParams serviceparams.Params `json:service_govparams`
}

func NewGenesisState(maxRequestTimeout int64, minDepositMultiple int64, serviceFeeTax, slashFraction sdk.Dec) GenesisState {
	return GenesisState{
		ServiceParams: serviceparams.Params{
			MaxRequestTimeout:  maxRequestTimeout,
			MinDepositMultiple: minDepositMultiple,
			ServiceFeeTax:      serviceFeeTax,
			SlashFraction:      slashFraction,
		},
	}
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	params.InitGenesisParameter(&serviceparams.ServiceParameter, ctx, data.ServiceParams)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		ServiceParams: serviceparams.GetSericeParams(ctx),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		ServiceParams: serviceparams.NewSericeParams(),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	serviceParams := serviceparams.NewSericeParams()
	serviceParams.MaxRequestTimeout = 10
	serviceParams.MinDepositMultiple = 10
	return GenesisState{
		ServiceParams: serviceParams,
	}
}
