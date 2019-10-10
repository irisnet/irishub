package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/irisnet/irishub/modules/asset/internal/types"
)

// NewAnteHandler returns an AnteHandler that checks if the balance of
// the fee payer is sufficient for asset related fee
func NewAnteHandler(ak keeper.AccountKeeper, k Keeper) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, simulate bool,
	) (newCtx sdk.Context, res sdk.Result, abort bool) {
		// new ctx
		newCtx = sdk.Context{}

		// all transactions must be of type auth.StdTx
		stdTx, ok := tx.(auth.StdTx)
		if !ok {
			// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
			// during runTx.
			return newCtx, sdk.ErrInternal("tx must be StdTx").Result(), true
		}

		// get the signing accouts
		signerAccs := stdTx.GetSigners()
		if len(signerAccs) == 0 {
			return newCtx, types.ErrSignersMissingInContext(types.DefaultCodespace, "signers missing in context").Result(), true
		}

		// get the payer
		payer, res := auth.GetSignerAcc(newCtx, ak, signerAccs[0])
		if !res.IsOK() {
			return newCtx, res, true
		}

		// total fee
		totalFee := sdk.Coins{}

		for _, msg := range tx.GetMsgs() {
			// only check consecutive msgs which are routed to asset from the beginning
			if msg.Route() != types.MsgRoute {
				break
			}

			var msgFee sdk.Coin

			switch msg := msg.(type) {
			case types.MsgCreateGateway:
				msgFee = GetGatewayCreateFee(ctx, k, msg.Moniker)
				break

			case types.MsgIssueToken:
				if msg.Source == types.NATIVE {
					msgFee = GetTokenIssueFee(ctx, k, msg.Symbol)
				} else if msg.Source == types.GATEWAY {
					msgFee = GetGatewayTokenIssueFee(ctx, k, msg.Symbol)
				}

				break

			case types.MsgMintToken:
				prefix, symbol := types.GetTokenIDParts(msg.TokenId)

				if prefix == "" || prefix == "i" {
					msgFee = GetTokenMintFee(ctx, k, symbol)
				} else if prefix != "x" {
					msgFee = GetGatewayTokenMintFee(ctx, k, symbol)
				}

				break

			default:
				msgFee = sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt())
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
