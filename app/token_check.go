package app

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
			if containCoinPrefix(sdk.NewCoins(msg.Token), coinswaptypes.FormatUniABSPrefix) {
				return ctx, sdkerrors.Wrap(
					sdkerrors.ErrInvalidRequest, "can't transfer coinswap liquidity tokens through the IBC module")
			}
		case *tokentypes.MsgBurnToken:
			if _, err := ctd.tk.GetToken(ctx, msg.Symbol); err != nil {
				return ctx, sdkerrors.Wrap(
					sdkerrors.ErrInvalidRequest, "burnt failed, only native tokens can be burnt")
			}
		case *govtypes.MsgSubmitProposal:
			if containCoinPrefix(msg.InitialDeposit, coinswaptypes.FormatUniABSPrefix) {
				return ctx, sdkerrors.Wrap(
					sdkerrors.ErrInvalidRequest, "can't deposit coinswap liquidity token for proposal")
			}
		case *govtypes.MsgDeposit:
			if containCoinPrefix(msg.Amount, coinswaptypes.FormatUniABSPrefix) {
				return ctx, sdkerrors.Wrap(
					sdkerrors.ErrInvalidRequest, "can't deposit coinswap liquidity  token %s for proposal")
			}
		}
	}
	return next(ctx, tx, simulate)
}

func containCoinPrefix(coins sdk.Coins, prefix string) bool {
	for _, coin := range coins {
		if strings.HasPrefix(coin.Denom, prefix) {
			return true
		}
	}
	return false
}
