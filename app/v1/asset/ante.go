package asset

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/auth"
	sdk "github.com/irisnet/irishub/types"
)

// NewAnteHandler returns an AnteHandler that checks if the balance of
// the fee payer is sufficient for asset related fee
func NewAnteHandler(k Keeper) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, simulate bool,
	) (newCtx sdk.Context, res sdk.Result, abort bool) {
		// new ctx
		newCtx = sdk.Context{}

		// get the signing accouts
		signerAccs := auth.GetSigners(ctx)
		if len(signerAccs) == 0 {
			return newCtx, ErrSignersMissingInContext(DefaultCodespace, "signers missing in context").Result(), true
		}

		// get the payer
		payer := signerAccs[0]

		// total fee
		totalFee := sdk.Coins{}

		for _, msg := range tx.GetMsgs() {
			// only check consecutive msgs which are routed to asset from the beginning
			if msg.Route() != MsgRoute {
				break
			}

			var msgFee sdk.Coin

			switch msg := msg.(type) {
			case MsgCreateGateway:
				msgFee = getGatewayCreateFee(ctx, k, msg.Moniker)
				break

			case MsgIssueToken:
				if msg.Source == NATIVE {
					msgFee = getTokenIssueFee(ctx, k, msg.Symbol)
				} else if msg.Source == GATEWAY {
					msgFee = getGatewayTokenIssueFee(ctx, k, msg.Symbol)
				}

				break

			case MsgMintToken:
				prefix, symbol := GetTokenIDParts(msg.TokenId)

				if prefix == "i" {
					msgFee = getTokenMintFee(ctx, k, symbol)
				} else if prefix != "x" {
					msgFee = getGatewayTokenMintFee(ctx, k, symbol)
				}

				break

			default:
				msgFee = sdk.NewCoin(sdk.NativeTokenMinDenom, sdk.ZeroInt())
			}

			totalFee = totalFee.Plus(sdk.Coins{msgFee})
		}

		if !totalFee.IsAllLTE(payer.GetCoins()) {
			// return error result and abort
			return newCtx, ErrInsufficientCoins(DefaultCodespace, fmt.Sprintf("insufficient coins for asset fee: %s needed", totalFee.MainUnitString())).Result(), true
		}

		// continue
		return newCtx, sdk.Result{}, false
	}
}
