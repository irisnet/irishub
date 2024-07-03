package oracle

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"

	modulev1 "mods.irisnet.org/api/irismod/oracle/module/v1"
	"mods.irisnet.org/modules/oracle/keeper"
	"mods.irisnet.org/modules/oracle/types"
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
	ServiceKeeper types.ServiceKeeper
}

// Outputs define the module outputs for the depinject.
type Outputs struct {
	depinject.Out

	OracleKeeper keeper.Keeper
	Module       appmodule.AppModule
}

// ProvideModule creates a new OracleKeeper and AppModule using the provided inputs and returns the Outputs.
//
// Parameters:
// - in: the Inputs struct containing the necessary dependencies for creating the OracleKeeper and AppModule.
//
// Returns:
// - Outputs: the struct containing the OracleKeeper and AppModule.
func ProvideModule(in Inputs) Outputs {
	keeper := keeper.NewKeeper(
		in.Cdc,
		in.Key,
		in.ServiceKeeper,
	)
	m := NewAppModule(in.Cdc, keeper, in.AccountKeeper, in.BankKeeper)

	return Outputs{OracleKeeper: keeper, Module: m}
}
