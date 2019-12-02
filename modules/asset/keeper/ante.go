package keeper

//import (
//	"fmt"
//
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/cosmos/cosmos-sdk/x/auth"
//	"github.com/irisnet/irishub/modules/asset/types"
//)
//
//// NewAnteHandler returns an AnteHandler that checks if the balance of
//// the fee payer is sufficient for asset related fee
//func NewAnteHandler(k Keeper) sdk.AnteHandler {
//	return func(
//		ctx sdk.Context, tx sdk.Tx, simulate bool,
//	) (newCtx sdk.Context, res error) {
//		// new ctx
//		newCtx = sdk.Context{}
//
//		// get the signing accouts
//		signerAccs := auth.GetSigners(ctx)
//		if len(signerAccs) == 0 {
//			return newCtx, types.ErrSignersMissingInContext(types.DefaultCodespace, "signers missing in context")
//		}
//
//		// get the payer
//		payer := signerAccs[0]
//
//		// total fee
//		totalFee := sdk.Coins{}
//
//		for _, msg := range tx.GetMsgs() {
//			// only check consecutive msgs which are routed to asset from the beginning
//			if msg.Route() != types.MsgRoute {
//				break
//			}
//
//			var msgFee sdk.Coin
//
//			switch msg := msg.(type) {
//			case types.MsgCreateGateway:
//				msgFee = GetGatewayCreateFee(ctx, k, msg.Moniker)
//				break
//
//			case types.MsgIssueToken:
//				if msg.Source == types.NATIVE {
//					msgFee = GetTokenIssueFee(ctx, k, msg.Symbol)
//				} else if msg.Source == types.GATEWAY {
//					msgFee = GetGatewayTokenIssueFee(ctx, k, msg.Symbol)
//				}
//
//				break
//
//			case types.MsgMintToken:
//				prefix, symbol := types.GetTokenIDParts(msg.TokenId)
//
//				if prefix == "" || prefix == "i" {
//					msgFee = GetTokenMintFee(ctx, k, symbol)
//				} else if prefix != "x" {
//					msgFee = GetGatewayTokenMintFee(ctx, k, symbol)
//				}
//
//				break
//
//			default:
//				msgFee = sdk.NewCoin(sdk.IrisAtto, sdk.ZeroInt())
//			}
//
//			totalFee = totalFee.Add(sdk.Coins{msgFee})
//		}
//
//		if !totalFee.IsAllLTE(payer.GetCoins()) {
//			// return error result and abort
//			return newCtx, types.ErrInsufficientCoins(types.DefaultCodespace, fmt.Sprintf("insufficient coins for asset fee: %s needed", totalFee.MainUnitString()))
//		}
//
//		// continue
//		return newCtx, nil
//	}
//}
