package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker handles block beginning logic for mint
func BeginBlocker(ctx sdk.Context, k Keeper) {
	logger := k.Logger(ctx)
	// Get block BFT time and block height
	blockTime := ctx.BlockHeader().Time
	minter := k.GetMinter(ctx)
	if ctx.BlockHeight() <= 1 { // don't inflate token in the first block
		minter.LastUpdate = blockTime
		k.SetMinter(ctx, minter)
		return
	}

	// Calculate block mint amount
	params := k.GetParamSet(ctx)
	logger.Info("Mint parameters", "inflation_rate", params.Inflation.String(), "mint_denom", params.MintDenom)

	mintedCoin := minter.BlockProvision(params)
	logger.Info("Mint result", "block_provisions", mintedCoin.String(), "time", blockTime.String())

	mintedCoins := sdk.NewCoins(mintedCoin)
	// mint coins to submodule account
	if err := k.MintCoins(ctx, mintedCoins); err != nil {
		panic(err)
	}

	// send the minted coins to the fee collector account
	if err := k.AddCollectedFees(ctx, mintedCoins); err != nil {
		panic(err)
	}

	// Update last block BFT time
	lastInflationTime := minter.LastUpdate
	minter.LastUpdate = blockTime
	k.SetMinter(ctx, minter)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeMint,
			sdk.NewAttribute(AttributeKeyLastInflationTime, lastInflationTime.String()),
			sdk.NewAttribute(AttributeKeyInflationTime, blockTime.String()),
			sdk.NewAttribute(AttributeKeyMintCoin, mintedCoin.Amount.String()),
		),
	)
}
