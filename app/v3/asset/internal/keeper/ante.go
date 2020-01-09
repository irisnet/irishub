package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/app/v3/asset/internal/types"

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
			return newCtx, types.ErrSignersMissingInContext(types.DefaultCodespace, "signers missing in context").Result(), true
		}

		// get the payer
		payer := signerAccs[0]

		// total fee
		totalFee := sdk.Coins{}

		for _, msg := range tx.GetMsgs() {
			// only check consecutive msgs which are routed to asset from the beginning
			if msg.Route() != types.MsgRoute {
				break
			}

			var msgFee sdk.Coin

			switch msg := msg.(type) {
			case types.MsgIssueToken:
				msgFee = k.getTokenIssueFee(ctx, msg.Symbol)
				break

			case types.MsgMintToken:
				_, symbol := types.GetTokenIDParts(msg.TokenId)
				msgFee = k.getTokenMintFee(ctx, symbol)
				break

			default:
				msgFee = sdk.NewCoin(sdk.IrisAtto, sdk.ZeroInt())
			}

			totalFee = totalFee.Add(sdk.Coins{msgFee})
		}

		if !totalFee.IsAllLTE(payer.GetCoins()) {
			// return error result and abort
			return newCtx, types.ErrInsufficientCoins(types.DefaultCodespace, fmt.Sprintf("insufficient coins for asset fee: %s needed", totalFee.MainUnitString())).Result(), true
		}

		// continue
		return newCtx, sdk.Result{}, false
	}
}
