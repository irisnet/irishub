package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"mods.irisnet.org/modules/token/contracts"
	"mods.irisnet.org/modules/token/types"
)

type erc20Hook struct {
	k Keeper
}

// PostTxProcessing processes the transaction receipt for ERC20 token swap to native.
//
// Parameters:
//
//	ctx: the context in which the function is executed.
//	msg: the core message associated with the transaction.
//	receipt: the Ethereum receipt containing the transaction logs.
//
// Return type: error
func (hook erc20Hook) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	disable := !hook.k.ERC20Enabled(ctx)
	erc20 := contracts.ERC20TokenContract.ABI
	for _, log := range receipt.Logs {
		// Note: the `SwapToNative` event contains 1 topics
		//  SwapToNative(address from, string to, uint256 amount)
		if len(log.Topics) != 1 {
			continue
		}

		// Check if event is included in ERC20
		eventID := log.Topics[0]
		event, err := erc20.EventByID(eventID)
		if err != nil {
			continue
		}

		// Check if event is a `SwapToNative` event.
		if event.Name != contracts.EventSwapToNative {
			continue
		}

		token, err := hook.k.getTokenByContract(ctx, log.Address)
		if err != nil {
			hook.k.Logger(ctx).Error("invalid SwapToNative event", "contract", log.Address.Hex())
			continue
		}

		if disable {
			return types.ErrERC20Disabled
		}

		eventArgs, err := erc20.Unpack(event.Name, log.Data)
		if err != nil {
			return errorsmod.Wrap(types.ErrInvalidContract, "failed to unpack SwapToNative event")
		}

		if len(eventArgs) != 3 {
			return errorsmod.Wrapf(
				types.ErrInvalidContract,
				"swapToNative event has wrong number of parameters, expected 3, actual: %d",
				len(eventArgs),
			)
		}

		to, ok := eventArgs[1].(string)
		if !ok || len(to) == 0 {
			return errorsmod.Wrap(
				types.ErrInvalidContract,
				"swapToNative event `to` parameters is invalid, expected string",
			)
		}

		receiver, err := sdk.AccAddressFromBech32(to)
		if err != nil {
			return errorsmod.Wrapf(
				types.ErrInvalidContract,
				"swapToNative event `to` parameters is invalid, expected iaa address, actual: %s",
				to,
			)
		}

		amount, ok := eventArgs[2].(*big.Int)
		if !ok || amount.Cmp(big.NewInt(0)) == 0 {
			return errorsmod.Wrap(
				types.ErrInvalidContract,
				"swapToNative event `amount` parameters is invalid, expected `*big.Int`",
			)
		}

		mintedCoins := sdk.NewCoins(sdk.NewCoin(token.MinUnit, sdkmath.NewIntFromBigInt(amount)))
		if err := hook.k.bankKeeper.MintCoins(ctx, types.ModuleName, mintedCoins); err != nil {
			return err
		}

		if err := hook.k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, mintedCoins); err != nil {
			return err
		}
	}
	return nil
}
