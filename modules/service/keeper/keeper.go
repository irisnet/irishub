package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"mods.irisnet.org/modules/service/types"
)

// Keeper defines the service keeper
type Keeper struct {
	storeKey         storetypes.StoreKey
	cdc              codec.Codec
	accountKeeper    types.AccountKeeper
	bankKeeper       types.BankKeeper
	blockedAddrs     map[string]bool
	feeCollectorName string
	authority        string                            // name of the fee collector
	respCallbacks    map[string]types.ResponseCallback // used to map the module name to response callback
	stateCallbacks   map[string]types.StateCallback    // used to map the module name to state callback
	moduleServices   map[string]*types.ModuleService   // used to map the module name to module service
}

// NewKeeper creates a new service Keeper instance
func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	feeCollectorName string,
	authority string,
) Keeper {
	// ensure service module accounts are set
	if addr := accountKeeper.GetModuleAddress(types.DepositAccName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.DepositAccName))
	}

	if addr := accountKeeper.GetModuleAddress(types.RequestAccName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.RequestAccName))
	}

	keeper := Keeper{
		storeKey:         key,
		cdc:              cdc,
		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
		blockedAddrs:     bankKeeper.GetBlockedAddresses(),
		feeCollectorName: feeCollectorName,
		authority:        authority,
	}

	keeper.respCallbacks = make(map[string]types.ResponseCallback)
	keeper.stateCallbacks = make(map[string]types.StateCallback)
	keeper.moduleServices = make(map[string]*types.ModuleService)

	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// GetServiceDepositAccount returns the service depost ModuleAccount
func (k Keeper) GetServiceDepositAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.DepositAccName)
}

// GetServiceRequestAccount returns the service request ModuleAccount
func (k Keeper) GetServiceRequestAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.RequestAccName)
}
