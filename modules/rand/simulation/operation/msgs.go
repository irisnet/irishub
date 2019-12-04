package operations

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	randmodule "github.com/irisnet/irishub/modules/rand"
)

// SimulateMsgRequestRand generates a MsgRequestRand with random values.
func SimulateMsgRequestRand(k randmodule.Keeper) simulation.Operation {
	handler := randmodule.NewHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		consumer := simulation.RandomAcc(r, accs)
		interval := simulation.RandIntBetween(r, 10, 100)

		msg := randmodule.NewMsgRequestRand(consumer.Address, uint64(interval))

		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(randmodule.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		ok := handler(ctx, msg).IsOK()
		if ok {
			write()
		}

		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}
