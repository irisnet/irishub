package v152

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/irisnet/irismod/modules/farm/types"
)

// Parameter store keys
var (
	DefaultPoolCreationFee = sdk.NewCoin("uiris", sdk.NewIntWithDecimal(2000, 6))
	DefaultTaxRate         = sdk.NewDecWithPrec(4, 1)
)

type (
	FarmKeeper interface {
		SetParams(ctx sdk.Context, params types.Params)
	}

	Params struct {
		PoolCreationFee     sdk.Coin `protobuf:"bytes,1,opt,name=pool_creation_fee,json=poolCreationFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Coin" json:"pool_creation_fee"`
		MaxRewardCategories uint32   `protobuf:"varint,2,opt,name=max_reward_categories,json=maxRewardCategories,proto3" json:"max_reward_categories,omitempty"`
		TaxRate             sdk.Dec  `protobuf:"bytes,3,opt,name=tax_rate,json=taxRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"tax_rate"`
	}
)

func Migrate(ctx sdk.Context, k FarmKeeper, ak types.AccountKeeper, paramSpace paramstypes.Subspace) error {
	params := GetLegacyParams(ctx, paramSpace)
	newParams := types.Params{
		MaxRewardCategories: params.MaxRewardCategories,
		PoolCreationFee:     DefaultPoolCreationFee,
		TaxRate:             DefaultTaxRate,
	}
	k.SetParams(ctx, newParams)

	//Grant burner permissions to the farm module account
	acc := ak.GetModuleAccount(ctx, types.ModuleName)
	if !acc.HasPermission(authtypes.Burner) {
		moduleAcc, _ := acc.(*authtypes.ModuleAccount)
		moduleAcc.Permissions = append(moduleAcc.Permissions, authtypes.Burner)
	}
	ak.SetModuleAccount(ctx, acc)
	return nil
}

// GetLegacyParams gets the parameters for the coinswap module.
func GetLegacyParams(ctx sdk.Context, paramSpace paramstypes.Subspace) Params {
	var swapParams Params
	paramSpace.GetParamSet(ctx, &swapParams)
	return swapParams
}

// ParamSetPairs implements paramtypes.KeyValuePairs
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(types.KeyPoolCreationFee, &p.PoolCreationFee, nil),
		paramstypes.NewParamSetPair(types.KeyMaxRewardCategories, &p.MaxRewardCategories, nil),
	}
}
