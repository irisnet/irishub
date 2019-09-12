package gov

import sdk "github.com/irisnet/irishub/types"

const (
	Insert string = "insert"
	Update string = "update"
)

type Param struct {
	Subspace string `json:"subspace"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type Params []Param

// Implements Proposal Interface
var _ Proposal = (*ParameterProposal)(nil)

type ParameterProposal struct {
	BasicProposal
	Params Params `json:"params"`
}

func (pp *ParameterProposal) Validate(ctx sdk.Context, k Keeper, verify bool) sdk.Error {
	if err := pp.BasicProposal.Validate(ctx, k, verify); err != nil {
		return err
	}

	param := pp.Params[0]
	if p, ok := k.paramsKeeper.GetParamSet(param.Subspace); ok {
		if _, err := p.Validate(param.Key, param.Value); err != nil {
			return err
		}
	} else {
		return ErrInvalidParam(DefaultCodespace, param.Subspace)
	}
	return nil
}

func (pp *ParameterProposal) Execute(ctx sdk.Context, k Keeper) sdk.Error {
	ctx.Logger().Info("Execute ParameterProposal begin")
	//check again
	if len(pp.Params) != 1 {
		ctx.Logger().Error("Execute ParameterProposal Failure", "info",
			"the length of ParameterProposal's param should be one", "ProposalId", pp.ProposalID)
		return nil
	}
	param := pp.Params[0]
	paramSet, _ := k.paramsKeeper.GetParamSet(param.Subspace)
	value, _ := paramSet.Validate(param.Key, param.Value)
	subspace, found := k.paramsKeeper.GetSubspace(param.Subspace)
	if found {
		k.metrics.AddParameter(param.Key, value)
		subspace.Set(ctx, []byte(param.Key), value)
		ctx.Logger().Info("Execute ParameterProposal Success", "key", param.Key, "value", param.Value)
	} else {
		ctx.Logger().Info("Execute ParameterProposal Failed", "key", param.Key, "value", param.Value)
	}

	return nil
}
