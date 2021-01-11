package app

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"

	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
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
			if strings.HasPrefix(msg.Token.Denom, coinswaptypes.FormatUniABSPrefix) {
				return ctx, sdkerrors.Wrap(
					sdkerrors.ErrInvalidRequest, "can't transfer coinswap coin from the ibc module")
			}
		case *tokentypes.MsgBurnToken:
			if _, err := ctd.tk.GetToken(ctx, msg.Symbol); err != nil {
				return ctx, sdkerrors.Wrapf(
					sdkerrors.ErrInvalidRequest, "can't burn token %s，only the token managed by the token module can be burned", msg.Symbol)
			}
		}
	}
	return next(ctx, tx, simulate)
}
