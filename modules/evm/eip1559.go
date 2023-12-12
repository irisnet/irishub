package evm

import (
	"math/big"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k *Keeper) burnBaseFee(ctx sdk.Context, gasUsed uint64, baseFee *big.Int) error {
	evmDenom := k.evmkeeper.GetParams(ctx).EvmDenom
	burntAmt := new(big.Int).Mul(baseFee, new(big.Int).SetUint64(gasUsed))
	burnCoin := sdk.NewCoin(evmDenom, sdkmath.NewIntFromBigInt(burntAmt))
	if err := k.bankKeeper.BurnCoins(ctx, authtypes.FeeCollectorName, sdk.NewCoins(burnCoin)); err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		EventEIP1559Burnt,
		sdk.NewAttribute(AttributeKeyBurntFee, burnCoin.String()),
		sdk.NewAttribute(AttributeKeyBaseFee, baseFee.String()),
	))
	return nil
}
