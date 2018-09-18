package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
	"github.com/irisnet/irishub/modules/parameter"
)

type Op string

const (
	Add    Op = "add"
	Update Op = "update"
)

type Param struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Op    Op     `json:"op"`
}

type Params []Param

// Implements Proposal Interface
var _ Proposal = (*ParameterProposal)(nil)

type ParameterProposal struct {
	TextProposal
	Param Param `json:"params"`
}

func (pp *ParameterProposal) Execute(ctx sdk.Context, k Keeper) (err error) {

	logger := ctx.Logger().With("module", "x/gov")
	logger.Info("Execute ParameterProposal begin", "info", fmt.Sprintf("current height:%d", ctx.BlockHeight()))
	parameter.ParamMapping[pp.Param.Key].Update(ctx,pp.Param.Value)

	return
}