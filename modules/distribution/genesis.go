package distribution

import (
	"fmt"

	"github.com/irisnet/irishub/modules/distribution/types"
	sdk "github.com/irisnet/irishub/types"
)

// InitGenesis sets distribution information for genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	if err := types.ValidateParams(data.Params); err != nil {
		panic(err.Error())
	}
	keeper.SetParams(ctx, data.Params)

	if !data.FeePool.ValPool.IsZero() {
		panic(fmt.Sprintf("Global validator pool(%s) is not zero", data.FeePool.ValPool.ToString()))
	}
	keeper.SetGenesisFeePool(ctx, data.FeePool)

	for _, vdi := range data.ValidatorDistInfos {
		if !vdi.ValCommission.IsZero() || !vdi.DelPool.IsZero() {
			panic(fmt.Sprintf("Delegation pool(%s) or commission pool(%s) are not zero", vdi.ValCommission.ToString(), vdi.DelPool.ToString()))
		}
		keeper.SetValidatorDistInfo(ctx, vdi)
	}
	for _, ddi := range data.DelegationDistInfos {
		keeper.SetDelegationDistInfo(ctx, ddi)
	}
	for _, dw := range data.DelegatorWithdrawInfos {
		keeper.SetDelegatorWithdrawAddr(ctx, dw.DelegatorAddr, dw.WithdrawAddr)
	}
	keeper.SetPreviousProposerConsAddr(ctx, data.PreviousProposer)
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, and validator/delegator distribution info's
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	feePool := keeper.GetFeePool(ctx)
	vdis := keeper.GetAllValidatorDistInfos(ctx)
	ddis := keeper.GetAllDelegationDistInfos(ctx)
	dwis := keeper.GetAllDelegatorWithdrawInfos(ctx)
	pp := keeper.GetPreviousProposerConsAddr(ctx)
	return NewGenesisState(params, feePool, vdis, ddis, dwis, pp)
}
