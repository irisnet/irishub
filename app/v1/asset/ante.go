package asset

import (
	"github.com/irisnet/irishub/app/v1/auth"
    sdk "github.com/irisnet/irishub/types"
)

// NewAnteHandler returns an AnteHandler that checks if the balance of
// the fee payer is sufficient for asset related fee
func NewAnteHandler(ak Keeper) sdk.AnteHandler {
    return func(
        ctx sdk.Context, tx sdk.Tx, simulate bool,
    ) (newCtx sdk.Context, res sdk.Result, abort bool) {
		
		// get the payer
		payer:= auth.GetSigners(ctx)[0]
		
        var totalFee sdk.Coins

        for _, msg := range tx.GetMsgs() {
            var msgFee sdk.Coin

            switch msg := msg.(type) {
            case MsgCreateGateway:
                msgFee = getGatewayCreateFee(ctx, ak, msg.Moniker)
                break

            case MsgIssueToken:
                if msg.Source == NATIVE {
                    msgFee = getTokenIssueFee(ctx, ak, msg.Symbol)
                } else if msg.Source == GATEWAY {
                    msgFee = getGatewayTokenIssueFee(ctx, ak, msg.Symbol)
                }

                break

            case MsgMintToken:
                prefix, symbol := GetTokenIDParts(msg.TokenId)

                if prefix == "i" {
                    msgFee = getTokenMintFee(ctx, ak, symbol)
                } else if prefix != "x" {
                    msgFee = getGatewayTokenMintFee(ctx, ak, symbol)
                }

                break

            default:
                msgFee = sdk.NewCoin(sdk.NativeTokenMinDenom, sdk.ZeroInt())
            }

            totalFee = totalFee.Plus(sdk.Coins{msgFee})
        }

        if !totalFee.IsAllLT(.GetCoins()) {
            // return error result and abort
            return ctx, ErrInsufficientFee(DefaultCodespace, "insufficient asset fee").Result(), true
        }

        // continue
        return ctx, sdk.Result{}, false
    }
}

