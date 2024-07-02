package record

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"

	modulev1 "mods.irisnet.org/api/irismod/record/module/v1"
	"mods.irisnet.org/modules/record/keeper"
	"mods.irisnet.org/modules/record/types"
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

// Inputs define the module inputs for the depinject.
type Inputs struct {
	depinject.In

	Config *modulev1.Module
	Cdc    codec.Codec
	Key    *store.KVStoreKey

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
}

// Outputs define the module outputs for the depinject.
type Outputs struct {
	depinject.Out

	RecordKeeper keeper.Keeper
	Module       appmodule.AppModule
}

// ProvideModule creates and returns the record module with the specified inputs.
//
// Takes Inputs as the parameter, which includes the codec, key, account keeper, and bank keeper.
// Returns Outputs containing the record keeper and the app module.
func ProvideModule(in Inputs) Outputs {
	keeper := keeper.NewKeeper(
		in.Cdc,
		in.Key,
	)
	m := NewAppModule(in.Cdc, keeper, in.AccountKeeper, in.BankKeeper)

	return Outputs{RecordKeeper: keeper, Module: m}
}
