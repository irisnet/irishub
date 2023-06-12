package v2

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/coinswap/exported"
	"github.com/irisnet/irismod/modules/coinswap/types"
)

// Parameter store keys
var (
	KeyFee                 = []byte("Fee") // fee key
	DefaultPoolCreationFee = sdk.NewCoin("uiris", sdkmath.NewIntWithDecimal(5000, 6))
	DefaultTaxRate         = sdk.NewDecWithPrec(4, 1)
)

type (
	CoinswapKeeper interface {
		GetParams(ctx sdk.Context) types.Params
		SetParams(ctx sdk.Context, params types.Params) error
	}

	Params struct {
		Fee sdk.Dec `protobuf:"bytes,1,opt,name=fee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"fee"`
	}
)

func Migrate(ctx sdk.Context, k CoinswapKeeper, paramSpace exported.Subspace) error {
	params := GetLegacyParams(ctx, paramSpace)
	newParams := types.Params{
		Fee:             params.Fee,
		PoolCreationFee: DefaultPoolCreationFee,
		TaxRate:         DefaultTaxRate,
	}
	return k.SetParams(ctx, newParams)
}

// GetLegacyParams gets the parameters for the coinswap module.
func GetLegacyParams(ctx sdk.Context, paramSpace exported.Subspace) Params {
	var swapParams Params
	paramSpace.GetParamSet(ctx, &swapParams)
	return swapParams
}

// ParamSetPairs implements paramtypes.KeyValuePairs
func (p *Params) ParamSetPairs() exported.ParamSetPairs {
	return exported.ParamSetPairs{
		exported.NewParamSetPair(KeyFee, &p.Fee, nil),
	}
}
