package nft

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	store "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"

	modulev1 "mods.irisnet.org/api/irismod/nft/module/v1"
	"mods.irisnet.org/modules/nft/keeper"
	"mods.irisnet.org/modules/nft/types"
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

// Inputs define the arguments used to instantiate an app module.
type Inputs struct {
	depinject.In

	Config *modulev1.Module
	Cdc    codec.Codec
	Key    *store.KVStoreKey

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
}

// Outputs define the read-only arguments return by depinject.
type Outputs struct {
	depinject.Out

	NFTKeeper keeper.Keeper
	Module    appmodule.AppModule
}

// ProvideModule provides a module for the NFT with the given inputs and returns the NFT keeper and module.
//
// Takes Inputs as input parameters and returns Outputs.
func ProvideModule(in Inputs) Outputs {
	keeper := keeper.NewKeeper(
		in.Cdc,
		in.Key,
		in.AccountKeeper,
		in.BankKeeper,
	)
	m := NewAppModule(in.Cdc, keeper, in.AccountKeeper, in.BankKeeper)

	return Outputs{NFTKeeper: keeper, Module: m}
}
