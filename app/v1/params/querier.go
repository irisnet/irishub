package params

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the params Querier
const (
	QueryModule = "module"
)

// creates a querier for params REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryModule:
			var params QueryModuleParams
			if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
				return nil, sdk.ParseParamsErr(err)
			}
			if len(params.Module) == 0 {
				return queryAll(ctx, keeper)
			}
			ps, err := queryModule(ctx, params.Module, keeper)
			if err != nil {
				return nil, sdk.ParseParamsErr(err)
			}
			bz, _ := keeper.cdc.MarshalJSON(ps)
			return bz, nil

		default:
			return nil, sdk.ErrUnknownRequest("unknown params query endpoint")
		}
	}
}

// defines the params for query: "custom/params/module"
type QueryModuleParams struct {
	Module string
}

func queryAll(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	var paramSets ParamSets
	for key := range keeper.spaces {
		ps, _ := queryModule(ctx, key, keeper)
		paramSets = append(paramSets, ps)
	}
	bz, err := keeper.cdc.MarshalJSON(paramSets)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}
	return bz, nil
}

func queryModule(ctx sdk.Context, module string, keeper Keeper) (ParamSet, sdk.Error) {
	subspace, ok := keeper.GetSubspace(module)
	if !ok {
		return nil, sdk.NewError(DefaultCodespace, CodeInvalidModule, fmt.Sprintf("The module %s is not existed or does not support params change", module))
	}

	ps, ok := keeper.GetParamSet(module)
	if !ok {
		return nil, sdk.NewError(DefaultCodespace, CodeInvalidModule, fmt.Sprintf("The module %s is does not support params change", module))
	}

	subspace.GetParamSet(ctx, ps)
	return ps, nil
}
