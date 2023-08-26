package evm

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	ethermint "github.com/evmos/ethermint/x/evm"
	"github.com/evmos/ethermint/x/evm/keeper"
	"github.com/evmos/ethermint/x/evm/types"
)

var (
	_ module.AppModule = AppModule{}
)

// ____________________________________________________________________________

// AppModule implements an application module for the evm module.
type AppModule struct {
	ethermint.AppModule
	k  *Keeper
	ss types.Subspace
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	k *keeper.Keeper,
	ak types.AccountKeeper,
	bankKeeper types.BankKeeper,
	ss types.Subspace,
) AppModule {
	return AppModule{
		AppModule: ethermint.NewAppModule(k, ak, ss),
		k:         &Keeper{k, bankKeeper, false},
		ss:        ss,
	}
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), am.k)
	types.RegisterQueryServer(cfg.QueryServer(), am.k.evmkeeper)

	m := keeper.NewMigrator(*am.k.evmkeeper, am.ss)

	if err := cfg.RegisterMigration(types.ModuleName, 3, m.Migrate3to4); err != nil {
		panic(err)
	}

	if err := cfg.RegisterMigration(types.ModuleName, 4, m.Migrate4to5); err != nil {
		panic(err)
	}
}
