// nolint
package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/modules/coinswap/types"
)

// DeductPoolCreationFee performs fee handling for creating liquidity pool
func (k Keeper) DeductPoolCreationFee(ctx sdk.Context, creator sdk.AccAddress) error {
	params := k.GetParams(ctx)
	poolCreationFee := params.PoolCreationFee

	// compute community tax and burned coin
	communityTaxCoin := sdk.NewCoin(poolCreationFee.Denom,
		sdk.NewDecFromInt(poolCreationFee.Amount).Mul(params.TaxRate).TruncateInt())
	burnedCoins := sdk.NewCoins(poolCreationFee.Sub(communityTaxCoin))

	// send all fees to module account
	if err := k.bk.SendCoinsFromAccountToModule(
		ctx, creator, types.ModuleName, sdk.NewCoins(poolCreationFee),
	); err != nil {
		return err
	}

	// send community tax to feeCollector
	if err := k.bk.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, sdk.NewCoins(communityTaxCoin)); err != nil {
		return err
	}

	// burn burnedCoin
	return k.bk.BurnCoins(ctx, types.ModuleName, burnedCoins)
}
