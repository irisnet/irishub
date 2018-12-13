package gov

import (
	"fmt"
	protocolKeeper "github.com/irisnet/irishub/app/protocol/keeper"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/types/common"
)

var _ Proposal = (*SoftwareUpgradeProposal)(nil)

type SoftwareUpgradeProposal struct {
	TextProposal
	Version      uint64
	Software     string
	SwitchHeight uint64
}

func (sp *SoftwareUpgradeProposal) Execute(ctx sdk.Context, k Keeper) error {

	logger := ctx.Logger().With("module", "x/gov")
	emptyUpgradeConfig := protocolKeeper.UpgradeConfig{}
	if k.pk.GetUpgradeConfig(ctx) == emptyUpgradeConfig {
		k.pk.SetUpgradeConfig(ctx,
			protocolKeeper.UpgradeConfig{sp.ProposalID,
				common.ProtocolDefinition{sp.Version, sp.Software, sp.SwitchHeight}})
		logger.Info("Execute SoftwareProposal begin", "info", fmt.Sprintf("current height:%d", ctx.BlockHeight()))

	} else {
		logger.Info("Software Upgrade Switch Period is in process.", "info", fmt.Sprintf("current height:%d", ctx.BlockHeight()))

	}

	return nil
}
