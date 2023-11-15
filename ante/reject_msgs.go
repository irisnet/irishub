package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// RejectMessagesDecorator prevents invalid msg types from being executed
type RejectMessagesDecorator struct{}

// AnteHandle rejects the following messages:
// 1. Messages that requires ethereum-specific authentication.
// For example `MsgEthereumTx` requires fee to be deducted in the antehandler in
// order to perform the refund.
// 2. Messages that creates vesting accounts.
func (rmd RejectMessagesDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		switch msg.(type) {
		case *evmtypes.MsgEthereumTx:
			return ctx, errorsmod.Wrapf(
				errortypes.ErrInvalidType,
				"MsgEthereumTx needs to be contained within a tx with 'ExtensionOptionsEthereumTx' option",
			)

		case *vestingtypes.MsgCreateVestingAccount,
			*vestingtypes.MsgCreatePermanentLockedAccount,
			*vestingtypes.MsgCreatePeriodicVestingAccount:
			return ctx, errortypes.Wrap(
				errortypes.ErrInvalidType,
				"currently doesn't support creating vesting account")
		}
	}
	return next(ctx, tx, simulate)
}
