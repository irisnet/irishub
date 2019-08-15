package bank

import (
	"fmt"
	"regexp"

	"github.com/irisnet/irishub/app/v1/auth"
	sdk "github.com/irisnet/irishub/types"
)

// NewAnteHandler returns an AnteHandler that checks if the memo of
// the tx satisfies the rule set by the recipient
func NewAnteHandler(ak auth.AccountKeeper) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, simulate bool,
	) (newCtx sdk.Context, res sdk.Result, abort bool) {
		// new ctx
		newCtx = sdk.Context{}

		// get memo
		stdTx := tx.(auth.StdTx)
		memo := stdTx.Memo

		for _, msg := range tx.GetMsgs() {
			// only check bank.MsgSend msg
			if msg.Route() == "bank" && msg.Type() == "send" {
				sendMsg := msg.(MsgSend)

				for _, output := range sendMsg.Outputs {
					acc := ak.GetAccount(ctx, output.Address)

					if acc != nil {
						memoRegexp := acc.GetMemoRegexp()

						if memoRegexp != "" {
							_, err := regexp.MatchString(memoRegexp, memo)
							if err != nil {
								return newCtx, ErrInvalidMemo(DefaultCodespace, fmt.Sprintf("invalid memo %s, expected regular expression: %s", memo, memoRegexp)).Result(), true
							}
						}
					}
				}
			}
		}

		// continue
		return newCtx, sdk.Result{}, false
	}
}
