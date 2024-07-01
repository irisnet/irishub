package token

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	modulev1 "mods.irisnet.org/api/irismod/token/module/v1"
	"mods.irisnet.org/modules/token/keeper"
	"mods.irisnet.org/modules/token/types"
	v1 "mods.irisnet.org/modules/token/types/v1"
)

// App Wiring Setup

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule, ProvideKeyTable),
	)
}

// ProvideKeyTable returns the key table for the Token module
func ProvideKeyTable() types.KeyTable {
	return v1.ParamKeyTable()
}

var _ appmodule.AppModule = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// TokenInputs is the input of the Token module
type TokenInputs struct {
	depinject.In

	Config *modulev1.Module
	Cdc    codec.Codec
	Key    *store.KVStoreKey

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
	EVMKeeper     types.EVMKeeper
	ICS20Keeper   types.ICS20Keeper

	// LegacySubspace is used solely for migration of x/params managed parameters
	LegacySubspace types.Subspace `optional:"true"`
}

// TokenOutputs is the output of the Token module
type TokenOutputs struct {
	depinject.Out

	TokenKeeper keeper.Keeper
	Module      appmodule.AppModule
}

// ProvideModule provides a module for the token with the given inputs and returns the token keeper and module.
//
// Takes TokenInputs as input parameters and returns TokenOutputs.
func ProvideModule(in TokenInputs) TokenOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	keeper := keeper.NewKeeper(
		in.Cdc,
		in.Key,
		in.BankKeeper,
		in.AccountKeeper,
		in.EVMKeeper,
		in.ICS20Keeper,
		in.Config.FeeCollectorName,
		authority.String(),
	)
	m := NewAppModule(in.Cdc, keeper, in.AccountKeeper, in.BankKeeper, in.LegacySubspace)

	return TokenOutputs{TokenKeeper: keeper, Module: m}
}
