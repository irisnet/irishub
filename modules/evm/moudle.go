package evm

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	ethermint "github.com/evmos/ethermint/x/evm"
	"github.com/evmos/ethermint/x/evm/keeper"
	"github.com/evmos/ethermint/x/evm/types"

	iristypes "github.com/irisnet/irishub/types"
)

var (
	_ module.AppModule = AppModule{}
)

// ____________________________________________________________________________

// AppModule implements an application module for the evm module.
type AppModule struct {
	ethermint.AppModule
	k *Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k *keeper.Keeper, ak types.AccountKeeper, bankKeeper types.BankKeeper) AppModule {
	return AppModule{
		AppModule: ethermint.NewAppModule(k, ak),
		k:         &Keeper{k, bankKeeper, false},
	}
}

// BeginBlock returns the begin block for the evm module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	ethChainID := iristypes.BuildEthChainID(ctx.ChainID())
	am.AppModule.BeginBlock(ctx.WithChainID(ethChainID), req)
}

// InitGenesis performs genesis initialization for the evm module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	ethChainID := iristypes.BuildEthChainID(ctx.ChainID())
	return am.AppModule.InitGenesis(ctx.WithChainID(ethChainID), cdc, data)
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), am.k)
	types.RegisterQueryServer(cfg.QueryServer(), am.k.evmkeeper)
}
