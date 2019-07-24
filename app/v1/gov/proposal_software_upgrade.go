package gov

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
)

var _ Proposal = (*SoftwareUpgradeProposal)(nil)

type SoftwareUpgradeProposal struct {
	BasicProposal
	ProtocolDefinition sdk.ProtocolDefinition `json:"protocol_definition"`
}

func (sp SoftwareUpgradeProposal) GetProtocolDefinition() sdk.ProtocolDefinition {
	return sp.ProtocolDefinition
}
func (sp *SoftwareUpgradeProposal) SetProtocolDefinition(upgrade sdk.ProtocolDefinition) {
	sp.ProtocolDefinition = upgrade
}

func (sp *SoftwareUpgradeProposal) Validate(ctx sdk.Context, k Keeper, verify bool) sdk.Error {
	if err := sp.BasicProposal.Validate(ctx, k, verify); err != nil {
		return err
	}

	if !k.protocolKeeper.IsValidVersion(ctx, sp.ProtocolDefinition.Version) {
		return ErrCodeInvalidVersion(k.codespace, sp.ProtocolDefinition.Version)
	}

	if uint64(ctx.BlockHeight()) > sp.ProtocolDefinition.Height {
		return ErrCodeInvalidSwitchHeight(k.codespace, uint64(ctx.BlockHeight()), sp.ProtocolDefinition.Height)
	}

	_, found := k.guardianKeeper.GetProfiler(ctx, sp.GetProposer())
	if !found {
		return ErrNotProfiler(k.codespace, sp.GetProposer())
	}

	if _, ok := k.protocolKeeper.GetUpgradeConfig(ctx); ok {
		return ErrSwitchPeriodInProcess(k.codespace)
	}
	return nil
}

func (sp *SoftwareUpgradeProposal) Execute(ctx sdk.Context, gk Keeper) sdk.Error {
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
