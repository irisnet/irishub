package app

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"

	guardiankeeper "github.com/irisnet/irishub/modules/guardian/keeper"
	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	tokenkeeper "github.com/irisnet/irismod/modules/token/keeper"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
)

// CheckTokenDecorator is responsible for restricting the token participation of the swap prefix
type CheckTokenDecorator struct {
	tk tokenkeeper.Keeper
}

// NewCheckTokenDecorator return a instance of CheckTokenDecorator
func NewCheckTokenDecorator(tk tokenkeeper.Keeper) CheckTokenDecorator {
	return CheckTokenDecorator{
		tk: tk,
	}
}

// AnteHandle check the transaction
func (ctd CheckTokenDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *ibctransfertypes.MsgTransfer:
			if containSwapCoin(msg.Token) {
				return ctx, sdkerrors.Wrap(
					sdkerrors.ErrInvalidRequest, "can't transfer coinswap liquidity tokens through the IBC module")
			}
		case *tokentypes.MsgBurnToken:
			if _, err := ctd.tk.GetToken(ctx, msg.Symbol); err != nil {
				return ctx, sdkerrors.Wrap(
					sdkerrors.ErrInvalidRequest, "burnt failed, only native tokens can be burnt")
			}
		case *govtypes.MsgSubmitProposal:
			if containSwapCoin(msg.InitialDeposit...) {
				return ctx, sdkerrors.Wrap(
					sdkerrors.ErrInvalidRequest, "can't deposit coinswap liquidity token for proposal")
			}
		case *govtypes.MsgDeposit:
			if containSwapCoin(msg.Amount...) {
				return ctx, sdkerrors.Wrap(
					sdkerrors.ErrInvalidRequest, "can't deposit coinswap liquidity token for proposal")
			}
		}
	}
	return next(ctx, tx, simulate)
}

// ServiceAuthDecorator is responsible for checking the permission to execute MsgCallService
type ServiceAuthDecorator struct {
	gk guardiankeeper.Keeper
}

// NewServiceAuthDecorator return a instance of ServiceAuthDecorator
func NewServiceAuthDecorator(gk guardiankeeper.Keeper) ServiceAuthDecorator {
	return ServiceAuthDecorator{
		gk: gk,
	}
}

// AnteHandle check the transaction
func (sad ServiceAuthDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *servicetypes.MsgCallService:
			if !msg.Repeated {
				continue
			}

			consumer, err := sdk.AccAddressFromBech32(msg.Consumer)
			if err != nil {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid consumer")
			}

			if !sad.gk.Authorized(ctx, consumer) {
				return ctx, sdkerrors.Wrap(
					sdkerrors.ErrInvalidRequest, "authentication failed, only super accounts can create repeated service invocation")
			}
		}
	}
	return next(ctx, tx, simulate)
}

func containSwapCoin(coins ...sdk.Coin) bool {
	for _, coin := range coins {
		if strings.HasPrefix(coin.Denom, coinswaptypes.FormatUniABSPrefix) {
			return true
		}
	}
	return false
}
