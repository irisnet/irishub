package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"mods.irisnet.org/modules/random/keeper"
	"mods.irisnet.org/modules/random/types"
	irishelpers "mods.irisnet.org/simapp/helpers"
)

// WeightedOperations generates a MsgRequestRandom with random values.
func WeightedOperations(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simulation.WeightedOperations {
	weightMsgRequestRandom := 100
	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgRequestRandom,
			func(
				r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
			) (
				simtypes.OperationMsg, []simtypes.FutureOperation, error,
			) {
				simAccount, _ := simtypes.RandomAcc(r, accs)
				blockInterval := simtypes.RandIntBetween(r, 10, 100)

				account := ak.GetAccount(ctx, simAccount.Address)

				spendable := bk.SpendableCoins(ctx, account.GetAddress())

				msg := types.NewMsgRequestRandom(
					simAccount.Address.String(),
					uint64(blockInterval),
					false,
					nil,
				)

				fees, err := simtypes.RandomFees(r, ctx, spendable)
				if err != nil {
					return simtypes.NoOpMsg(
						types.ModuleName,
						msg.Type(),
						"unable to generate fees",
					), nil, err
				}

				txConfig := moduletestutil.MakeTestEncodingConfig().TxConfig
				tx, err := irishelpers.GenTx(
					r,
					txConfig,
					[]sdk.Msg{msg},
					fees,
					simtestutil.DefaultGenTxGas,
					chainID,
					[]uint64{account.GetAccountNumber()},
					[]uint64{account.GetSequence()},
					simAccount.PrivKey,
				)
				if err != nil {
					return simtypes.NoOpMsg(
						types.ModuleName,
						msg.Type(),
						"unable to generate mock tx",
					), nil, err
				}

				if _, _, err := app.SimDeliver(txConfig.TxEncoder(), tx); err != nil {
					return simtypes.NoOpMsg(
						types.ModuleName,
						msg.Type(),
						"unable to deliver tx",
					), nil, err
				}

				return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
			},
		),
	}
}
