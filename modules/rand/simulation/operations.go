package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irisnet/irishub/modules/rand/internal/keeper"
	"github.com/irisnet/irishub/modules/rand/internal/types"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(appParams simulation.AppParams, cdc *codec.Codec, ak types.AccountKeeper, k keeper.Keeper) simulation.WeightedOperations {
	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			128,
			SimulateMsgRequestRand(ak, k),
		),
	}
}

// SimulateMsgRequestRand generates a MsgRequestRand with random values.
func SimulateMsgRequestRand(ak types.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		simAccount, _ := simulation.RandomAcc(r, accs)
		blockInterval := simulation.RandIntBetween(r, 10, 100)
		account := ak.GetAccount(ctx, simAccount.Address)

		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgRequestRand(simAccount.Address, uint64(blockInterval))

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
