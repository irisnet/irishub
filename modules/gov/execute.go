package gov

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
)

func Execute(ctx sdk.Context, gk Keeper, p Proposal) (err error) {
	switch p.GetProposalType() {
	case ProposalTypeParameterChange:
		return ParameterProposalExecute(ctx, gk, p.(*ParameterProposal))
	case ProposalTypeSystemHalt:
		return SystemHaltProposalExecute(ctx, gk)
	case ProposalTypeTxTaxUsage:
		return TaxUsageProposalExecute(ctx, gk, p.(*TaxUsageProposal))
	case ProposalTypeSoftwareUpgrade:
		return SoftwareUpgradeProposalExecute(ctx, gk, p.(*SoftwareUpgradeProposal))
	}
	return nil
}

func TaxUsageProposalExecute(ctx sdk.Context, gk Keeper, p *TaxUsageProposal) (err error) {
	burn := false
	if p.TaxUsage.Usage == UsageTypeBurn {
		burn = true
	} else {
		_, found := gk.guardianKeeper.GetTrustee(ctx, p.TaxUsage.DestAddress)
		if !found {
			ctx.Logger().Error("Execute TaxUsageProposal Failure", "info",
				"the destination address is not a trustee now", "destinationAddress", p.TaxUsage.DestAddress)
			return
		}
	}
	gk.dk.AllocateFeeTax(ctx, p.TaxUsage.DestAddress, p.TaxUsage.Percent, burn)
	return
}

func ParameterProposalExecute(ctx sdk.Context, gk Keeper, pp *ParameterProposal) (err error) {
	ctx.Logger().Info("Execute ParameterProposal begin")
	for _, param := range pp.Params {
		paramSet, _ := gk.paramsKeeper.GetParamSet(param.Subspace)
		value, _ := paramSet.Validate(param.Key, param.Value)
		subspace, found := gk.paramsKeeper.GetSubspace(param.Subspace)
		if found {
			SetParameterMetrics(gk.metrics, param.Key, value)
			subspace.Set(ctx, []byte(param.Key), value)
			ctx.Logger().Info("Execute ParameterProposal Successed", "key", param.Key, "value", param.Value)
		} else {
			ctx.Logger().Info("Execute ParameterProposal Failed", "key", param.Key, "value", param.Value)
		}

	}

	return
}

func SoftwareUpgradeProposalExecute(ctx sdk.Context, gk Keeper, sp *SoftwareUpgradeProposal) error {

	if _, ok := gk.protocolKeeper.GetUpgradeConfig(ctx); ok {
		ctx.Logger().Info("Execute SoftwareProposal Failure", "info",
			fmt.Sprintf("Software Upgrade Switch Period is in process."))
		return nil
	}
	if !gk.protocolKeeper.IsValidVersion(ctx, sp.ProtocolDefinition.Version) {
		ctx.Logger().Info("Execute SoftwareProposal Failure", "info",
			fmt.Sprintf("version [%v] in SoftwareUpgradeProposal isn't valid ", sp.ProposalID))
		return nil
	}
	if uint64(ctx.BlockHeight())+1 >= sp.ProtocolDefinition.Height {
		ctx.Logger().Info("Execute SoftwareProposal Failure", "info",
			fmt.Sprintf("switch height must be more than blockHeight + 1"))
		return nil
	}

	gk.protocolKeeper.SetUpgradeConfig(ctx, sdk.NewUpgradeConfig(sp.ProposalID, sp.ProtocolDefinition))

	ctx.Logger().Info("Execute SoftwareProposal Success")

	return nil
}

func SystemHaltProposalExecute(ctx sdk.Context, gk Keeper) error {
	logger := ctx.Logger()

	if gk.GetSystemHaltHeight(ctx) == -1 {
		gk.SetSystemHaltHeight(ctx, ctx.BlockHeight()+gk.GetSystemHaltPeriod(ctx))
		logger.Info("Execute SystemHaltProposal begin", "SystemHaltHeight", gk.GetSystemHaltHeight(ctx))
	} else {
		logger.Info("SystemHalt Period is in process", "SystemHaltHeight", gk.GetSystemHaltHeight(ctx))

	}
	return nil
}
