package mt

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"

	modulev1 "mods.irisnet.org/api/irismod/mt/module/v1"
	"mods.irisnet.org/modules/mt/keeper"
	"mods.irisnet.org/modules/mt/types"
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

	MTKeeper keeper.Keeper
	Module   appmodule.AppModule
}

// ProvideModule creates a new MTKeeper and AppModule using the provided inputs and returns the Outputs.
//
// Parameters:
// - in: the Inputs struct containing the necessary dependencies for creating the MTKeeper and AppModule.
//
// Returns:
// - Outputs: the struct containing the MTKeeper and AppModule.
func ProvideModule(in Inputs) Outputs {
	keeper := keeper.NewKeeper(
		in.Cdc,
		in.Key,
	)
	m := NewAppModule(in.Cdc, keeper, in.AccountKeeper, in.BankKeeper)

	return Outputs{MTKeeper: keeper, Module: m}
}
