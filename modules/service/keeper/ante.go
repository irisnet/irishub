package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/service/types"
)

// ValidateServiceAuthDecorator is responsible for checking the permission to execute MsgCallService
type ValidateServiceAuthDecorator struct {
	ak types.AuthKeeper
}

// NewValidateServiceAuthDecorator returns an instance of ServiceAuthDecorator
func NewValidateServiceAuthDecorator(ak types.AuthKeeper) ValidateServiceAuthDecorator {
	return ValidateServiceAuthDecorator{
		ak: ak,
	}
}

// AnteHandle checks the transaction
func (vsad ValidateServiceAuthDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *types.MsgCallService:
			if msg.Repeated || msg.SuperMode {
				consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
				if err != nil {
					return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid consumer")
				}

				if !vsad.ak.Authorized(ctx, consumer) {
					return ctx, sdkerrors.Wrap(
						sdkerrors.ErrUnauthorized, "authentication failed, only super accounts can create repeated service invocation or create service invocation with superMode")
				}
			}
		}
	}
	return next(ctx, tx, simulate)
}
