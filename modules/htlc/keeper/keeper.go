package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"irismod.io/htlc/types"
)

// Keeper defines the HTLC keeper
type Keeper struct {
	storeKey      storetypes.StoreKey
	cdc           codec.Codec
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	blockedAddrs  map[string]bool
	authority     string
}

// NewKeeper creates a new HTLC Keeper instance
func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	authority string,
) Keeper {
	// ensure the HTLC module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		blockedAddrs:  bankKeeper.GetBlockedAddresses(),
		authority:     authority,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// GetHTLCAccount returns the HTLC module account
func (k Keeper) GetHTLCAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// EnsureModuleAccountPermissions syncs the bep3 module account's permissions with those in the supply keeper.
func (k Keeper) EnsureModuleAccountPermissions(ctx sdk.Context) error {
	maccI := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	macc, ok := maccI.(*authtypes.ModuleAccount)
	if !ok {
		return fmt.Errorf("expected %s account to be a module account type", types.ModuleName)
	}
	_, perms := k.accountKeeper.GetModuleAddressAndPermissions(types.ModuleName)
	macc.Permissions = perms
	k.accountKeeper.SetModuleAccount(ctx, macc)
	return nil
}
