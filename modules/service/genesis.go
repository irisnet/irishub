package service

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/service/params"
	"github.com/irisnet/irishub/modules/params"
)

// GenesisState - all service state that must be provided at genesis
type GenesisState struct {
	MaxRequestTimeout  int64
	MinDepositMultiple int64
	ServiceFeeTax      sdk.Dec
	SlashFraction      sdk.Dec
}

func NewGenesisState(maxRequestTimeout int64, minDepositMultiple int64, serviceFeeTax, slashFraction sdk.Dec) GenesisState {
	return GenesisState{
		MaxRequestTimeout:  maxRequestTimeout,
		MinDepositMultiple: minDepositMultiple,
		ServiceFeeTax:      serviceFeeTax,
		SlashFraction:      slashFraction,
	}
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetServiceFeeTax(ctx, data.ServiceFeeTax)
	k.SetServiceSlashFraction(ctx, data.SlashFraction)
	params.InitGenesisParameter(&serviceparams.MaxRequestTimeoutParameter, ctx, data.MaxRequestTimeout)
	params.InitGenesisParameter(&serviceparams.MinDepositMultipleParameter, ctx, data.MinDepositMultiple)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	maxRequestTimeout := serviceparams.GetMaxRequestTimeout(ctx)
	minDepositMultiple := serviceparams.GetMinDepositMultiple(ctx)

	return GenesisState{
		MaxRequestTimeout:  maxRequestTimeout,
		MinDepositMultiple: minDepositMultiple,
		ServiceFeeTax:      k.GetServiceFeeTax(ctx),
		SlashFraction:      k.GetServiceSlashFraction(ctx),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		MaxRequestTimeout:  100,
		MinDepositMultiple: 1000,
		ServiceFeeTax:      sdk.NewDecWithPrec(2, 2), //2%
		SlashFraction:      sdk.NewDecWithPrec(1, 2), //1%
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		MaxRequestTimeout:  10,
		MinDepositMultiple: 10,
		ServiceFeeTax:      sdk.NewDecWithPrec(2, 2), //2%
		SlashFraction:      sdk.NewDecWithPrec(1, 2), //1%
	}
}
