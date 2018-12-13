package gov

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/params"
)

const (
	Insert string = "insert"
	Update string = "update"
)

type Param struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Op    string `json:"op"`
}

// Implements Proposal Interface
var _ Proposal = (*ParameterProposal)(nil)

type ParameterProposal struct {
	TextProposal
	Param Param `json:"params"`
}

func (pp *ParameterProposal) Execute(ctx sdk.Context, k Keeper) (err error) {

	logger := ctx.Logger().With("module", "x/gov")
	logger.Info("Execute ParameterProposal begin", "info", fmt.Sprintf("current height:%d", ctx.BlockHeight()))
	if pp.Param.Op == Update {
		params.ParamMapping[pp.Param.Key].Update(ctx, pp.Param.Value)
	} else if pp.Param.Op == Insert {
		//Todo: insert
	}
	return
}
