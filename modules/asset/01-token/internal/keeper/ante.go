package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/asset/01-token/internal/types"
)

// ValidateTokenFeeDecorator is responsible for withholding fees on transactions issued in msg and additional tokens.
type ValidateTokenFeeDecorator struct {
	tk Keeper
	ak types.AccountKeeper
}

func NewValidateTokenFeeDecorator(tk Keeper, ak types.AccountKeeper) ValidateTokenFeeDecorator {
	return ValidateTokenFeeDecorator{
		tk: tk,
		ak: ak,
	}
}

func (dtf ValidateTokenFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeMap := make(map[string]sdk.Coin)
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case types.MsgIssueToken:
			fee := GetTokenIssueFee(ctx, dtf.tk, msg.Symbol)
			if fe, ok := feeMap[msg.Owner.String()]; ok {
				feeMap[msg.Owner.String()] = fe.Add(fee)
			} else {
				feeMap[msg.Owner.String()] = fee
			}
		case types.MsgMintToken:
			fee := GetTokenMintFee(ctx, dtf.tk, msg.Symbol)
			if fe, ok := feeMap[msg.Owner.String()]; ok {
				feeMap[msg.Owner.String()] = fe.Add(fee)
			} else {
				feeMap[msg.Owner.String()] = fee
			}
		}
	}

	for addr, fee := range feeMap {
		owner, _ := sdk.AccAddressFromBech32(addr)
		account := dtf.ak.GetAccount(ctx, owner)
		if account.GetCoins().IsAllLT(sdk.NewCoins(fee)) {
			return ctx, sdk.ErrInsufficientCoins(fmt.Sprintf("account[%s] insufficient coins for asset fee: %s needed", addr, fee.String()))
		}
	}
	return next(ctx, tx, simulate)
}
