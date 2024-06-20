package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"

	modulev1 "github.com/irisnet/irismod/api/irismod/nft/module/v1"
	"github.com/irisnet/irismod/nft/keeper"
	"github.com/irisnet/irismod/nft/types"
)

// App Wiring Setup
func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

var _ appmodule.AppModule = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

type NFTInputs struct {
	depinject.In

	Config *modulev1.Module
	Cdc    codec.Codec
	Key    *store.KVStoreKey

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
}

type NFTOutputs struct {
	depinject.Out

	NFTKeeper keeper.Keeper
	Module    appmodule.AppModule
}

func ProvideModule(in NFTInputs) NFTOutputs {
	keeper := keeper.NewKeeper(
		in.Cdc,
		in.Key,
		in.AccountKeeper,
		in.BankKeeper,
	)
	m := NewAppModule(in.Cdc, keeper, in.AccountKeeper, in.BankKeeper)

	return NFTOutputs{NFTKeeper: keeper, Module: m}
}
