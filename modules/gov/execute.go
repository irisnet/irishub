package gov

import (
	"fmt"
	protocolKeeper "github.com/irisnet/irishub/app/protocol/keeper"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/params"
)

func Execute(ctx sdk.Context, k Keeper, p Proposal) (err error) {
	switch p.GetProposalType() {
	case ProposalTypeParameterChange:
		return ParameterProposalExecute(ctx, k, p.(*ParameterProposal))
	case ProposalTypeSystemHalt:
		return SystemHaltProposalExecute(ctx, k)
	case ProposalTypeTxTaxUsage:
		return TaxUsageProposalExecute(ctx, k, p.(*TaxUsageProposal))
	case ProposalTypeSoftwareUpgrade:
		return SoftwareUpgradeProposalExecute(ctx, k, p.(*SoftwareUpgradeProposal))
	}
	return nil
}

func TaxUsageProposalExecute(ctx sdk.Context, k Keeper, p *TaxUsageProposal) (err error) {
	burn := false
	if p.Usage == UsageTypeBurn {
		burn = true
	}
	k.dk.AllocateFeeTax(ctx, p.DestAddress, p.Percent, burn)
	return
}

func ParameterProposalExecute(ctx sdk.Context, k Keeper, pp *ParameterProposal) (err error) {

	logger := ctx.Logger().With("module", "x/gov")
	logger.Info("Execute ParameterProposal begin", "info", fmt.Sprintf("current height:%d", ctx.BlockHeight()))
    for _, param := range pp.Params{
    	paramSet := params.ParamSetMapping[param.Subspace]
    	value, _ := paramSet.Validate(param.Key,param.Value)
		subspace, bool := k.paramsKeeper.GetSubspace(param.Subspace)
		if bool {
			subspace.Set(ctx,[]byte(param.Key),value)
		}

		logger.Info("Execute ParameterProposal begin", "info", fmt.Sprintf("%s = %s",param.Key,param.Value ))
	}

	return
}

func SoftwareUpgradeProposalExecute(ctx sdk.Context, k Keeper, sp *SoftwareUpgradeProposal) error {

	logger := ctx.Logger().With("module", "x/gov")

	if _, ok := k.pk.GetUpgradeConfig(ctx); ok {
		logger.Info("Execute SoftwareProposal Failure", "info",
			fmt.Sprintf("Software Upgrade Switch Period is in process. current height:%d", ctx.BlockHeight()))
		return nil
	}
	if !k.pk.IsValidProtocolVersion(ctx, sp.Version) {
		logger.Info("Execute SoftwareProposal Failure", "info",
			fmt.Sprintf("version [%v] in SoftwareUpgradeProposal isn't valid ", sp.ProposalID))
		return nil
	}
	if uint64(ctx.BlockHeight())+1 >= sp.SwitchHeight {
		logger.Info("Execute SoftwareProposal Failure", "info",
			fmt.Sprintf("switch height must be more than blockHeight + 1"))
		return nil
	}

	k.pk.SetUpgradeConfig(ctx,
		protocolKeeper.UpgradeConfig{sp.ProposalID,
			sdk.ProtocolDefinition{sp.Version, sp.Software, sp.SwitchHeight}})

	logger.Info("Execute SoftwareProposal Success", "info",
		fmt.Sprintf("current height:%d", ctx.BlockHeight()))

	return nil
}

func SystemHaltProposalExecute(ctx sdk.Context, k Keeper) error {
	logger := ctx.Logger().With("module", "x/gov")

	if k.GetSystemHaltHeight(ctx) == -1 {
		k.SetSystemHaltHeight(ctx, ctx.BlockHeight()+k.GetSystemHaltPeriod(ctx))
		logger.Info("Execute SystemHaltProposal begin", "info", fmt.Sprintf("SystemHalt height:%d", k.GetSystemHaltHeight(ctx)))
	} else {
		logger.Info("SystemHalt Period is in process.", "info", fmt.Sprintf("SystemHalt height:%d", k.GetSystemHaltHeight(ctx)))

	}
	return nil
}
