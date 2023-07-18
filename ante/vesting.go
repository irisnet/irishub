package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// RejectVestingDecorator is responsible for rejecting the vesting msg
type RejectVestingDecorator struct{}

// NewRejectVestingDecorator returns an instance of ValidateVestingDecorator
func NewRejectVestingDecorator() RejectVestingDecorator {
	return RejectVestingDecorator{}
}

// AnteHandle checks the transaction
func (vvd RejectVestingDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		switch msg.(type) {
		case *vestingtypes.MsgCreateVestingAccount, *vestingtypes.MsgCreatePermanentLockedAccount, *vestingtypes.MsgCreatePeriodicVestingAccount:
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "currently doesn't support creating vesting account")
		}
	}
	return next(ctx, tx, simulate)
}
