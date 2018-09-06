package gov

import (
	"fmt"
	"strings"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	Params Params `json:"params"`
}

func (pp *ParameterProposal) Execute(ctx sdk.Context, k Keeper) (err error) {

	logger := ctx.Logger().With("module", "x/gov")
	logger.Info("Execute ParameterProposal begin", "info", fmt.Sprintf("current height:%d", ctx.BlockHeight()))

	for _, data := range pp.Params {
		//param only begin with "gov/" can be update
		if !strings.HasPrefix(data.Key, Prefix) {
			errMsg := fmt.Sprintf("Parameter %s is not begin with %s", data.Key, Prefix)
			logger.Error("Execute ParameterProposal ", "err", errMsg)
			continue
		}
		if data.Op == Add {
			k.ps.Set(ctx, data.Key, data.Value)
		} else if data.Op == Update {
			bz := k.ps.GetRaw(ctx, data.Key)
			if bz == nil || len(bz) == 0 {
				logger.Error("Execute ParameterProposal ", "err", "Parameter "+data.Key+" is not exist")
			} else {
				k.ps.SetString(ctx, data.Key, data.Value)
			}
		}
	}
	return
}