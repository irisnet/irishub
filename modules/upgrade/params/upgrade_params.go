package upgradeparams

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/iparam"
)

var CurrentUpgradeProposalIdParameter CurrentUpgradeProposalIdParam

var _ iparam.SignalParameter = (*CurrentUpgradeProposalIdParam)(nil)

type CurrentUpgradeProposalIdParam struct {
	Value      uint64
	paramSpace params.Subspace
}

func (param *CurrentUpgradeProposalIdParam) InitGenesis(genesisState interface{}) {
	param.Value = 0
}

func (param *CurrentUpgradeProposalIdParam) SetReadWriter(paramSpace params.Subspace) {
	param.paramSpace = paramSpace
}

func (param *CurrentUpgradeProposalIdParam) GetStoreKey() []byte {
	return []byte("upgradeProposalId")
}

func (param *CurrentUpgradeProposalIdParam) SaveValue(ctx sdk.Context) {
	param.paramSpace.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *CurrentUpgradeProposalIdParam) LoadValue(ctx sdk.Context) bool {
	if param.paramSpace.Has(ctx, param.GetStoreKey()) == false {
		return false
	}
	param.paramSpace.Get(ctx, param.GetStoreKey(), &param.Value)
	return true
}

var ProposalAcceptHeightParameter ProposalAcceptHeightParam

var _ iparam.SignalParameter = (*ProposalAcceptHeightParam)(nil)

type ProposalAcceptHeightParam struct {
	Value      int64
	paramSpace params.Subspace
}

func (param *ProposalAcceptHeightParam) InitGenesis(genesisState interface{}) {
	param.Value = -1
}

func (param *ProposalAcceptHeightParam) SetReadWriter(paramSpace params.Subspace) {
	param.paramSpace = paramSpace
}

func (param *ProposalAcceptHeightParam) GetStoreKey() []byte {
	return []byte("upgradeProposalAcceptHeight")
}

func (param *ProposalAcceptHeightParam) SaveValue(ctx sdk.Context) {
	param.paramSpace.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *ProposalAcceptHeightParam) LoadValue(ctx sdk.Context) bool {
	if param.paramSpace.Has(ctx, param.GetStoreKey()) == false {
		return false
	}
	param.paramSpace.Get(ctx, param.GetStoreKey(), &param.Value)
	return true
}

var SwitchPeriodParameter SwitchPeriodParam

var _ iparam.SignalParameter = (*SwitchPeriodParam)(nil)

type SwitchPeriodParam struct {
	Value      int64
	paramSpace params.Subspace
}

func (param *SwitchPeriodParam) InitGenesis(genesisState interface{}) {
	param.Value = genesisState.(int64)
}

func (param *SwitchPeriodParam) SetReadWriter(paramSpace params.Subspace) {
	param.paramSpace = paramSpace
}

func (param *SwitchPeriodParam) GetStoreKey() []byte {
	return []byte("upgradeSwitchPeriod")
}

func (param *SwitchPeriodParam) SaveValue(ctx sdk.Context) {
	param.paramSpace.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *SwitchPeriodParam) LoadValue(ctx sdk.Context) bool {
	if param.paramSpace.Has(ctx, param.GetStoreKey()) == false {
		return false
	}
	param.paramSpace.Get(ctx, param.GetStoreKey(), &param.Value)
	return true
}
