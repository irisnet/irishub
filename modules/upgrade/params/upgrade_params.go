package upgradeparams

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/irisnet/irishub/modules/parameter"
)

var CurrentUpgradeProposalIdParameter CurrentUpgradeProposalIdParam

var _ parameter.SignalParameter = (*CurrentUpgradeProposalIdParam)(nil)

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
