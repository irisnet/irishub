package random

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"

	modulev1 "mods.irisnet.org/api/irismod/random/module/v1"
	"mods.irisnet.org/modules/random/keeper"
	"mods.irisnet.org/modules/random/types"
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

type RandomInputs struct {
	depinject.In

	Config *modulev1.Module
	Cdc    codec.Codec
	Key    *store.KVStoreKey

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
	ServiceKeeper types.ServiceKeeper
}

type RandomOutputs struct {
	depinject.Out

	RandomKeeper keeper.Keeper
	Module       appmodule.AppModule
}

func ProvideModule(in RandomInputs) RandomOutputs {
	keeper := keeper.NewKeeper(
		in.Cdc,
		in.Key,
		in.BankKeeper,
		in.ServiceKeeper,
	)
	m := NewAppModule(in.Cdc, keeper, in.AccountKeeper, in.BankKeeper)

	return RandomOutputs{RandomKeeper: keeper, Module: m}
}
