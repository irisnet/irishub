package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mods.irisnet.org/modules/token/types"
	v1 "mods.irisnet.org/modules/token/types/v1"
)

type ValidateTokenFeeDecorator struct {
	k  Keeper
	bk types.BankKeeper
}

func NewValidateTokenFeeDecorator(k Keeper, bk types.BankKeeper) ValidateTokenFeeDecorator {
	return ValidateTokenFeeDecorator{
		k:  k,
		bk: bk,
	}
}

// AnteHandle returns an AnteHandler that checks if the balance of
// the fee payer is sufficient for token related fee
func (dtf ValidateTokenFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// total fee
	feeMap := make(map[string]sdk.Coin)
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *v1.MsgIssueToken:
			fee, err := dtf.k.GetTokenIssueFee(ctx, msg.Symbol)
			if err != nil {
				return ctx, errorsmod.Wrap(types.ErrInvalidBaseFee, err.Error())
			}

			if fe, ok := feeMap[msg.Owner]; ok {
				feeMap[msg.Owner] = fe.Add(fee)
			} else {
				feeMap[msg.Owner] = fee
			}
		case *v1.MsgMintToken:
			symbol, err := dtf.k.getSymbolByMinUnit(ctx, msg.Coin.Denom)
			if err != nil {
				return ctx, err
			}

			fee, err := dtf.k.GetTokenMintFee(ctx, symbol)
			if err != nil {
				return ctx, errorsmod.Wrap(types.ErrInvalidBaseFee, err.Error())
			}

			if fe, ok := feeMap[msg.Owner]; ok {
				feeMap[msg.Owner] = fe.Add(fee)
			} else {
				feeMap[msg.Owner] = fee
			}
		}
	}

	for addr, fee := range feeMap {
		owner, _ := sdk.AccAddressFromBech32(addr)
		balance := dtf.bk.GetBalance(ctx, owner, fee.Denom)
		if balance.IsLT(fee) {
			return ctx, errorsmod.Wrapf(
				sdkerrors.ErrInsufficientFunds, "insufficient coins for token fee; %s < %s", balance, fee,
			)
		}
	}
	// continue
	return next(ctx, tx, simulate)
}
