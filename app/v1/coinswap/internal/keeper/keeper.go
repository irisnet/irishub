package keeper

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"

	"github.com/irisnet/irishub/app/v1/coinswap/internal/types"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"

	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the coinswap store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	bk         bank.Keeper
	ak         auth.AccountKeeper
	paramSpace params.Subspace
}

// NewKeeper returns a coinswap keeper. It handles:
// - creating new ModuleAccounts for each trading pair
// - burning minting liquidity coins
// - sending to and from ModuleAccounts
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk bank.Keeper, ak auth.AccountKeeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		bk:         bk,
		ak:         ak,
		cdc:        cdc,
		paramSpace: paramSpace.WithTypeTable(types.ParamTypeTable()),
	}
}

// CreateReservePool initializes a new reserve pool by creating a
// ModuleAccount with minting and burning permissions.
func (keeper Keeper) CreateReservePool(ctx sdk.Context, moduleName string) {
	moduleAcc := keeper.ak.GetAccount(ctx, auth.SwapPoolAccAddr)
	if moduleAcc != nil {
		panic(fmt.Sprintf("reserve pool for %s already exists", moduleName))
	}
	keeper.bk.AddCoins(ctx, auth.SwapPoolAccAddr, sdk.Coins{})
}

// HasCoins returns whether or not an account has at least coins.
func (keeper Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, coins ...sdk.Coin) bool {
	return keeper.bk.HasCoins(ctx, addr, coins)
}

// BurnCoins burns liquidity coins from the ModuleAccount at moduleName. The
// moduleName and denomination of the liquidity coins are the same.
func (keeper Keeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Int) {
	_, err := keeper.bk.BurnCoins(ctx, auth.SwapPoolAccAddr, sdk.Coins{sdk.NewCoin(moduleName, amt)})
	if err != nil {
		panic(err)
	}
}

// MintCoins mints liquidity coins to the ModuleAccount at moduleName. The
// moduleName and denomination of the liquidity coins are the same.
func (keeper Keeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Int) {
	//err := keeper.bk.MintCoins(ctx, moduleName, sdk.NewCoins(sdk.NewCoin(moduleName, amt)))
	//if err != nil {
	//	panic(err)
	//}
}

// SendCoin sends coins from the address to the ModuleAccount at moduleName.
func (keeper Keeper) SendCoins(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins ...sdk.Coin) {
	_, err := keeper.bk.SendCoins(ctx, addr, auth.SwapPoolAccAddr, coins)
	if err != nil {
		panic(err)
	}
}

// RecieveCoin sends coins from the ModuleAccount at moduleName to the
// address provided.
func (keeper Keeper) RecieveCoins(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins ...sdk.Coin) {
	_, err := keeper.bk.SendCoins(ctx, auth.SwapPoolAccAddr, addr, coins)
	if err != nil {
		panic(err)
	}
}

// GetReservePool returns the total balance of an reserve pool at the
// provided denomination.
func (keeper Keeper) GetReservePool(ctx sdk.Context, moduleName string) (coins sdk.Coins, found bool) {
	acc := keeper.ak.GetAccount(ctx, auth.SwapPoolAccAddr)
	if acc == nil {
		return nil, false
	}
	return acc.GetCoins(), true
}

// GetNativeDenom returns the native denomination for this module from the
// global param store.
func (keeper Keeper) GetNativeDenom(ctx sdk.Context) (nativeDenom string) {
	return keeper.GetParams(ctx).NativeDenom
}

// GetFeeParam returns the current FeeParam from the global param store
func (keeper Keeper) GetFeeParam(ctx sdk.Context) (feeParam types.FeeParam) {
	return keeper.GetParams(ctx).Fee
}

// GetParams gets the parameters for the coinswap module.
func (keeper Keeper) GetParams(ctx sdk.Context) types.Params {
	var swapParams types.Params
	keeper.paramSpace.GetParamSet(ctx, &swapParams)
	return swapParams
}

// SetParams sets the parameters for the coinswap module.
func (keeper Keeper) SetParams(ctx sdk.Context, params types.Params) {
	keeper.paramSpace.SetParamSet(ctx, &params)
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
