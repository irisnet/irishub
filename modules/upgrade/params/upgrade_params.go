package upgradeparams

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/irisnet/irishub/modules/iparam"
)

var CurrentUpgradeProposalIdParameter CurrentUpgradeProposalIdParam

var _ iparam.SignalParameter = (*CurrentUpgradeProposalIdParam)(nil)

type CurrentUpgradeProposalIdParam struct {
	Value   int64
	psetter params.Setter
	pgetter params.Getter
}

func (param *CurrentUpgradeProposalIdParam) InitGenesis(genesisState interface{}) {
	param.Value = -1
}

func (param *CurrentUpgradeProposalIdParam) SetReadWriter(setter params.Setter) {
	param.psetter = setter
	param.pgetter = setter.Getter
}

func (param *CurrentUpgradeProposalIdParam) GetStoreKey() string {
	return "Sig/upgrade/proposalId"
}

func (param *CurrentUpgradeProposalIdParam) SaveValue(ctx sdk.Context) {
	param.psetter.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *CurrentUpgradeProposalIdParam) LoadValue(ctx sdk.Context) bool {
	err := param.pgetter.Get(ctx, param.GetStoreKey(), &param.Value)
	if err != nil {
		return false
	}
	return true
}


var ProposalAcceptHeightParameter ProposalAcceptHeightParam

var _ iparam.SignalParameter = (*ProposalAcceptHeightParam)(nil)

type ProposalAcceptHeightParam struct {
	Value   int64
	psetter params.Setter
	pgetter params.Getter
}

func (param *ProposalAcceptHeightParam) InitGenesis(genesisState interface{}) {
	param.Value = -1
}

func (param *ProposalAcceptHeightParam) SetReadWriter(setter params.Setter) {
	param.psetter = setter
	param.pgetter = setter.Getter
}

func (param *ProposalAcceptHeightParam) GetStoreKey() string {
	return "Sig/upgrade/proposalAcceptHeight"
}

func (param *ProposalAcceptHeightParam) SaveValue(ctx sdk.Context) {
	param.psetter.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *ProposalAcceptHeightParam) LoadValue(ctx sdk.Context) bool {
	err := param.pgetter.Get(ctx, param.GetStoreKey(), &param.Value)
	if err != nil {
		return false
	}
	return true
}


var SwitchPeriodParameter SwitchPeriodParam

var _ iparam.SignalParameter = (*SwitchPeriodParam)(nil)

type SwitchPeriodParam struct {
	Value   int64
	psetter params.Setter
	pgetter params.Getter
}

func (param *SwitchPeriodParam) InitGenesis(genesisState interface{}) {
	param.Value = genesisState.(int64)
}

func (param *SwitchPeriodParam) SetReadWriter(setter params.Setter) {
	param.psetter = setter
	param.pgetter = setter.Getter
}

func (param *SwitchPeriodParam) GetStoreKey() string {
	return "Sig/upgrade/switchperiod"
}

func (param *SwitchPeriodParam) SaveValue(ctx sdk.Context) {
	param.psetter.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *SwitchPeriodParam) LoadValue(ctx sdk.Context) bool {
	err := param.pgetter.Get(ctx, param.GetStoreKey(), &param.Value)
	if err != nil {
		return false
	}
	return true
}
