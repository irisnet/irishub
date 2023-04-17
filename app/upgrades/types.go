package upgrades

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	tibckeeper "github.com/bianjieai/tibc-go/modules/tibc/core/keeper"

	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"

	htlckeeper "github.com/irisnet/irismod/modules/htlc/keeper"
	servicekeeper "github.com/irisnet/irismod/modules/service/keeper"
	tokenkeeper "github.com/irisnet/irismod/modules/token/keeper"
)

// Upgrade defines a struct containing necessary fields that a SoftwareUpgradeProposal
// must have written, in order for the state migration to go smoothly.
// An upgrade must implement this struct, and then set it in the app.go.
// The app.go will then define the handler.
type Upgrade struct {
	// Upgrade version name, for the upgrade handler, e.g. `v7`
	UpgradeName string

	// UpgradeHandlerConstructor defines the function that creates an upgrade handler
	UpgradeHandlerConstructor func(*module.Manager, module.Configurator, AppKeepers) upgradetypes.UpgradeHandler

	// Store upgrades, should be used for any new modules introduced, new modules deleted, or store names renamed.
	StoreUpgrades *store.StoreUpgrades
}

type ConsensusParamsReaderWriter interface {
	StoreConsensusParams(ctx sdk.Context, cp *abci.ConsensusParams)
	GetConsensusParams(ctx sdk.Context) *abci.ConsensusParams
}

type AppKeepers struct {
	AppCodec        codec.Codec
	HTLCKeeper      htlckeeper.Keeper
	AccountKeeper   authkeeper.AccountKeeper
	BankKeeper      bankkeeper.Keeper
	ServiceKeeper   servicekeeper.Keeper
	GetKey          func(moduleName string) *storetypes.KVStoreKey
	ModuleManager   *module.Manager
	TIBCkeeper      *tibckeeper.Keeper
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper
	TokenKeeper     tokenkeeper.Keeper
	ReaderWriter    ConsensusParamsReaderWriter
}
