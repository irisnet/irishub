package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/app/v3/asset/internal/types"
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
		// total fee
		feeMap := make(map[string]sdk.Coins)
		var senders []string
		for _, msg := range tx.GetMsgs() {
			// only check consecutive msgs which are routed to asset from the beginning
			if msg.Route() != types.MsgRoute {
				break
			}
			sender := msg.GetSigners()[0].String()
			if _, ok := feeMap[sender]; !ok {
				feeMap[sender] = sdk.Coins{}
				senders = append(senders, sender)
			}

			var fee sdk.Coin
			switch msg := msg.(type) {
			case types.MsgIssueToken:
				fee = k.getTokenIssueFee(ctx, msg.Symbol)
				break
			case types.MsgMintToken:
				fee = k.getTokenMintFee(ctx, msg.Symbol)
				break
			}
			feeMap[sender] = feeMap[sender].Add(sdk.NewCoins(fee))
		}

		for _, addr := range senders {
			owner, _ := sdk.AccAddressFromBech32(addr)
			balance := k.bk.GetCoins(ctx, owner)
			if balance.IsAllLT(feeMap[addr]) {
				return newCtx, types.ErrInsufficientCoins(types.DefaultCodespace, fmt.Sprintf("insufficient coins for asset fee: %s needed", feeMap[addr].MainUnitString())).Result(), true
			}
		}
		// continue
		return newCtx, sdk.Result{}, false
	}
}
