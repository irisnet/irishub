package service

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	modulev1 "mods.irisnet.org/api/irismod/service/module/v1"
	"mods.irisnet.org/modules/service/keeper"
	"mods.irisnet.org/modules/service/types"
)

// App Wiring Setup

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule, ProvideKeyTable),
	)
}

// ProvideKeyTable returns the KeyTable for the service module.
//
// It calls the ParamKeyTable function from the types package to retrieve the KeyTable.
// The KeyTable is used to register parameter sets for the service module.
//
// Returns:
// - types.KeyTable: The KeyTable for the service module.
func ProvideKeyTable() types.KeyTable {
	return types.ParamKeyTable()
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

	// LegacySubspace is used solely for migration of x/params managed parameters
	LegacySubspace types.Subspace `optional:"true"`
}

// Outputs define the module outputs for the depinject.
type Outputs struct {
	depinject.Out

	ServiceKeeper keeper.Keeper
	Module        appmodule.AppModule
}

// ProvideModule creates and returns the HTLC module with the specified inputs.
//
// It takes Inputs as the parameter, which includes the configuration, codec, key, account keeper, and bank keeper.
// It returns Outputs containing the HTLC keeper and the app module.
func ProvideModule(in Inputs) Outputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	keeper := keeper.NewKeeper(
		in.Cdc,
		in.Key,
		in.AccountKeeper,
		in.BankKeeper,
		in.Config.FeeCollectorName,
		authority.String(),
	)
	m := NewAppModule(in.Cdc, keeper, in.AccountKeeper, in.BankKeeper, in.LegacySubspace)

	return Outputs{ServiceKeeper: keeper, Module: m}
}
