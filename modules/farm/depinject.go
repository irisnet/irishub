package farm

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	modulev1 "mods.irisnet.org/api/irismod/farm/module/v1"
	"mods.irisnet.org/modules/farm/keeper"
	"mods.irisnet.org/modules/farm/types"
)

// App Wiring Setup

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule, ProvideKeyTable),
	)
}

// ProvideKeyTable returns the KeyTable for the farm module's parameters.
//
// No parameters.
// Returns a types.KeyTable.
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

	AccountKeeper  types.AccountKeeper
	BankKeeper     types.BankKeeper
	DistrKeeper    types.DistrKeeper
	GovKeeper      types.GovKeeper
	CoinswapKeeper types.CoinswapKeeper

	// LegacySubspace is used solely for migration of x/params managed parameters
	LegacySubspace types.Subspace `optional:"true"`
}

// Outputs define the module outputs for the depinject.
type Outputs struct {
	depinject.Out

	FarmKeeper keeper.Keeper
	Module     appmodule.AppModule
}

// ProvideModule creates and returns the farm module with the specified inputs.
//
// It takes Inputs as the parameter, which includes the configuration, codec, key, account keeper, bank keeper, governance keeper, coinswap keeper, and legacy subspace.
// It returns Outputs containing the farm keeper and the app module.
func ProvideModule(in Inputs) Outputs {
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
		in.DistrKeeper,
		in.GovKeeper,
		in.CoinswapKeeper,
		in.Config.FeeCollectorName,
		in.Config.CommunityPoolName,
		authority.String(),
	)
	m := NewAppModule(in.Cdc, keeper, in.AccountKeeper, in.BankKeeper, in.LegacySubspace)

	return Outputs{FarmKeeper: keeper, Module: m}
}
