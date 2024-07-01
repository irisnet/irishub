package coinswap

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	modulev1 "mods.irisnet.org/api/irismod/coinswap/module/v1"
	"mods.irisnet.org/modules/coinswap/keeper"
	"mods.irisnet.org/modules/coinswap/types"
)

// App Wiring Setup

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule, ProvideKeyTable),
	)
}

func ProvideKeyTable() types.KeyTable {
	return types.ParamKeyTable() 
}

var _ appmodule.AppModule = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

type CoinswapInputs struct {
	depinject.In

	Config *modulev1.Module
	Cdc    codec.Codec
	Key    *store.KVStoreKey

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper

	// LegacySubspace is used solely for migration of x/params managed parameters
	LegacySubspace types.Subspace `optional:"true"`
}

type CoinswapOutputs struct {
	depinject.Out

	CoinswapKeeper keeper.Keeper
	Module         appmodule.AppModule
}

func ProvideModule(in CoinswapInputs) CoinswapOutputs {
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
		in.Config.FeeCollectorName,
		authority.String(),
	)
	m := NewAppModule(in.Cdc, keeper, in.AccountKeeper, in.BankKeeper, in.LegacySubspace)

	return CoinswapOutputs{CoinswapKeeper: keeper, Module: m}
}
