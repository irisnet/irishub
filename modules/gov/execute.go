package gov

import (
	"fmt"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/types/common"
	sdk "github.com/irisnet/irishub/types"
	govtypes "github.com/irisnet/irishub/types/gov"
	protocolKeeper "github.com/irisnet/irishub/app/protocol/keeper"
)


func Execute(ctx sdk.Context, k Keeper, p govtypes.Proposal) (err error){
	switch p.GetProposalType(){
	case govtypes.ProposalTypeParameterChange:
		return ParameterProposalExecute(ctx, k, p.(*govtypes.ParameterProposal))
	case govtypes.ProposalTypeSoftwareHalt:
		return HaltProposalExecute(ctx, k)
	case govtypes.ProposalTypeTxTaxUsage:
		return TaxUsageProposalExecute(ctx,k, p.(*govtypes.TaxUsageProposal))
	case govtypes.ProposalTypeSoftwareUpgrade:
		return SoftwareUpgradeProposalExecute(ctx, k, p.(*govtypes.SoftwareUpgradeProposal))
	}
	return nil
}

func TaxUsageProposalExecute(ctx sdk.Context, k Keeper, p *govtypes.TaxUsageProposal) (err error) {
	burn := false
	if p.Usage == govtypes.UsageTypeBurn {
		burn = true
	}
	k.dk.AllocateFeeTax(ctx, p.DestAddress, p.Percent, burn)
	return
}

func ParameterProposalExecute(ctx sdk.Context, k Keeper, pp *govtypes.ParameterProposal) (err error) {

	logger := ctx.Logger().With("module", "x/gov")
	logger.Info("Execute ParameterProposal begin", "info", fmt.Sprintf("current height:%d", ctx.BlockHeight()))
	if pp.Param.Op == govtypes.Update {
		params.ParamMapping[pp.Param.Key].Update(ctx, pp.Param.Value)
	} else if pp.Param.Op == govtypes.Insert {
		//Todo: insert
	}
	return
}

func SoftwareUpgradeProposalExecute(ctx sdk.Context, k Keeper, sp *govtypes.SoftwareUpgradeProposal) error {

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
			common.ProtocolDefinition{sp.Version, sp.Software, sp.SwitchHeight}})

	logger.Info("Execute SoftwareProposal Success", "info",
		fmt.Sprintf("current height:%d", ctx.BlockHeight()))

	return nil
}

func HaltProposalExecute(ctx sdk.Context, k Keeper) error {
	logger := ctx.Logger().With("module", "x/gov")

	if k.GetTerminatorHeight(ctx) == -1 {
		k.SetTerminatorHeight(ctx, ctx.BlockHeight()+k.GetTerminatorPeriod(ctx))
		logger.Info("Execute TerminatorProposal begin", "info", fmt.Sprintf("Terminator height:%d", k.GetTerminatorHeight(ctx)))
	} else {
		logger.Info("Terminator Period is in process.", "info", fmt.Sprintf("Terminator height:%d", k.GetTerminatorHeight(ctx)))

	}
	return nil
}