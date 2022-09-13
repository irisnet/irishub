package app

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"

	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	tokenkeeper "github.com/irisnet/irismod/modules/token/keeper"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
)

// ValidateTokenDecorator is responsible for restricting the token participation of the swap prefix
type ValidateTokenDecorator struct {
	tk tokenkeeper.Keeper
}

// NewValidateTokenDecorator returns an instance of ValidateTokenDecorator
func NewValidateTokenDecorator(tk tokenkeeper.Keeper) ValidateTokenDecorator {
	return ValidateTokenDecorator{
		tk: tk,
	}
}

// AnteHandle checks the transaction
func (vtd ValidateTokenDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *ibctransfertypes.MsgTransfer:
			if containSwapCoin(msg.Token) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "can't transfer coinswap liquidity tokens through the IBC module")
			}
		case *tokentypes.MsgBurnToken:
			if _, err := vtd.tk.GetToken(ctx, msg.Symbol); err != nil {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "burnt failed, only native tokens can be burnt")
			}
		case *govv1.MsgSubmitProposal:
			if containSwapCoin(msg.InitialDeposit...) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "can't deposit coinswap liquidity token for proposal")
			}
		case *govv1.MsgDeposit:
			if containSwapCoin(msg.Amount...) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "can't deposit coinswap liquidity token for proposal")
			}
		}
	}
	return next(ctx, tx, simulate)
}

// ValidateServiceDecorator is responsible for checking the permission to execute MsgCallService
type ValidateServiceDecorator struct{}

// NewValidateServiceDecorator returns an instance of ServiceAuthDecorator
func NewValidateServiceDecorator() ValidateServiceDecorator {
	return ValidateServiceDecorator{}
}

// AnteHandle checks the transaction
func (vsd ValidateServiceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *servicetypes.MsgCallService:
			if msg.Repeated {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "currently does not support to create repeatable service invocation")
			}
		}
	}
	return next(ctx, tx, simulate)
}

func containSwapCoin(coins ...sdk.Coin) bool {
	for _, coin := range coins {
		if strings.HasPrefix(coin.Denom, coinswaptypes.LptTokenPrefix) {
			return true
		}
	}
	return false
}
