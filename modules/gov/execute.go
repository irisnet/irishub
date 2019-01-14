package gov

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/params"
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
	logger := ctx.Logger().With("module", "gov")
	burn := false
	if p.TaxUsage.Usage == UsageTypeBurn {
		burn = true
	} else {
		_, found := gk.guardianKeeper.GetTrustee(ctx, p.TaxUsage.DestAddress)
		if !found {
			logger.Error("Execute TaxUsageProposal Failure", "info",
				fmt.Sprintf("the destination address [%s] is not a trustee now", p.TaxUsage.DestAddress))
			return
		}
	}
	gk.dk.AllocateFeeTax(ctx, p.TaxUsage.DestAddress, p.TaxUsage.Percent, burn)
	return
}

func ParameterProposalExecute(ctx sdk.Context, gk Keeper, pp *ParameterProposal) (err error) {
	logger := ctx.Logger().With("module", "gov")
	logger.Info("Execute ParameterProposal begin", "info", fmt.Sprintf("current height:%d", ctx.BlockHeight()))
	for _, param := range pp.Params {
		paramSet := params.ParamSetMapping[param.Subspace]
		value, _ := paramSet.Validate(param.Key, param.Value)
		subspace, bool := gk.paramsKeeper.GetSubspace(param.Subspace)
		if bool {
			subspace.Set(ctx, []byte(param.Key), value)
		}

		logger.Info("Execute ParameterProposal begin", "info", fmt.Sprintf("%s = %s", param.Key, param.Value))
	}

	return
}

func SoftwareUpgradeProposalExecute(ctx sdk.Context, gk Keeper, sp *SoftwareUpgradeProposal) error {
	logger := ctx.Logger().With("module", "x/gov")

	if _, ok := gk.protocolKeeper.GetUpgradeConfig(ctx); ok {
		logger.Info("Execute SoftwareProposal Failure", "info",
			fmt.Sprintf("Software Upgrade Switch Period is in process. current height:%d", ctx.BlockHeight()))
		return nil
	}
	if !gk.protocolKeeper.IsValidVersion(ctx, sp.ProtocolDefinition.Version) {
		logger.Info("Execute SoftwareProposal Failure", "info",
			fmt.Sprintf("version [%v] in SoftwareUpgradeProposal isn't valid ", sp.ProposalID))
		return nil
	}
	if uint64(ctx.BlockHeight())+1 >= sp.ProtocolDefinition.Height {
		logger.Info("Execute SoftwareProposal Failure", "info",
			fmt.Sprintf("switch height must be more than blockHeight + 1"))
		return nil
	}

	gk.protocolKeeper.SetUpgradeConfig(ctx, sdk.NewUpgradeConfig(sp.ProposalID, sp.ProtocolDefinition))

	logger.Info("Execute SoftwareProposal Success", "info",
		fmt.Sprintf("current height:%d", ctx.BlockHeight()))

	return nil
}

func SystemHaltProposalExecute(ctx sdk.Context, gk Keeper) error {
	logger := ctx.Logger().With("module", "x/gov")

	if gk.GetSystemHaltHeight(ctx) == -1 {
		gk.SetSystemHaltHeight(ctx, ctx.BlockHeight()+gk.GetSystemHaltPeriod(ctx))
		logger.Info("Execute SystemHaltProposal begin", "info", fmt.Sprintf("SystemHalt height:%d", gk.GetSystemHaltHeight(ctx)))
	} else {
		logger.Info("SystemHalt Period is in process.", "info", fmt.Sprintf("SystemHalt height:%d", gk.GetSystemHaltHeight(ctx)))

	}
	return nil
}
