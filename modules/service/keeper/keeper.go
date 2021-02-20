package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/irisnet/irismod/modules/service/types"
)

// Keeper defines the service keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.Marshaler

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	paramSpace    paramstypes.Subspace

	// name of the fee collector
	feeCollectorName string

	// used to map the module name to response callback
	respCallbacks map[string]types.ResponseCallback

	// used to map the module name to state callback
	stateCallbacks map[string]types.StateCallback

	// used to map the module name to module service
	moduleServices map[string]*types.ModuleService
}

// NewKeeper creates a new service Keeper instance
func NewKeeper(
	cdc codec.Marshaler,
	key sdk.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	paramSpace paramstypes.Subspace,
	feeCollectorName string,
) Keeper {
	// ensure service module accounts are set
	if addr := accountKeeper.GetModuleAddress(types.DepositAccName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.DepositAccName))
	}

	if addr := accountKeeper.GetModuleAddress(types.RequestAccName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.RequestAccName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(ParamKeyTable())
	}

	keeper := Keeper{
		storeKey:         key,
		cdc:              cdc,
		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
		feeCollectorName: feeCollectorName,
		paramSpace:       paramSpace,
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
