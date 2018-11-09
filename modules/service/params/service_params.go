package serviceparams

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/irisnet/irishub/iparam"
)

var MaxRequestTimeoutParameter MaxRequestTimeoutParam

var _ iparam.SignalParameter = (*MaxRequestTimeoutParam)(nil)

type MaxRequestTimeoutParam struct {
	Value      int64
	paramSpace params.Subspace
}

func (param *MaxRequestTimeoutParam) InitGenesis(genesisState interface{}) {
	param.Value = genesisState.(int64)
}

func (param *MaxRequestTimeoutParam) SetReadWriter(paramSpace params.Subspace) {
	param.paramSpace = paramSpace
}

func (param *MaxRequestTimeoutParam) GetStoreKey() []byte {
	return []byte("serviceMaxRequestTimeout")
}

func (param *MaxRequestTimeoutParam) SaveValue(ctx sdk.Context) {
	param.paramSpace.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *MaxRequestTimeoutParam) LoadValue(ctx sdk.Context) bool {
	if param.paramSpace.Has(ctx, param.GetStoreKey()) == false {
		return false
	}
	param.paramSpace.Get(ctx, param.GetStoreKey(), &param.Value)
	return true
}

var MinProviderDepositParameter MinProviderDepositParam
var _ iparam.SignalParameter = (*MinProviderDepositParam)(nil)

type MinProviderDepositParam struct {
	Value      sdk.Coins
	paramSpace params.Subspace
}

func (param *MinProviderDepositParam) InitGenesis(genesisState interface{}) {
	param.Value = genesisState.(sdk.Coins)
}

func (param *MinProviderDepositParam) SetReadWriter(paramSpace params.Subspace) {
	param.paramSpace = paramSpace
}

func (param *MinProviderDepositParam) GetStoreKey() []byte {
	return []byte("serviceMinProviderDeposit")
}

func (param *MinProviderDepositParam) SaveValue(ctx sdk.Context) {
	param.paramSpace.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *MinProviderDepositParam) LoadValue(ctx sdk.Context) bool {
	if param.paramSpace.Has(ctx, param.GetStoreKey()) == false {
		return false
	}
	param.paramSpace.Get(ctx, param.GetStoreKey(), &param.Value)
	return true
}
