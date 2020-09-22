package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irisnet/irismod/modules/record/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreateRecord = "op_weight_msg_create_record"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONMarshaler,
	ak types.AccountKeeper,
	bk types.BankKeeper) simulation.WeightedOperations {
	var weightCreate int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateRecord, &weightCreate, nil,
		func(_ *rand.Rand) {
			weightCreate = 50
		},
	)
	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightCreate,
			SimulateCreateRecord(ak, bk),
		),
	}
}

// SimulateCreateRecord tests and runs a single msg create a new record
func SimulateCreateRecord(ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		record, err := genRecord(r, accs)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeCreateRecord, err.Error()), nil, err
		}

		msg := types.NewMsgCreateRecord(record.Contents, record.Creator)

		simAccount, found := simtypes.FindAccount(accs, record.Creator)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeCreateRecord, err.Error()), nil, fmt.Errorf("account %s not found", record.Creator)
		}

		account := ak.GetAccount(ctx, msg.Creator)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeCreateRecord, err.Error()), nil, err
		}
		txGen := simappparams.MakeEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if _, _, err = app.Deliver(tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeCreateRecord, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate issue token"), nil, nil
	}
}

func genRecord(r *rand.Rand, accs []simtypes.Account) (types.Record, error) {
	var record types.Record
	txHash := make([]byte, 32)
	_, err := r.Read(txHash)
	if err != nil {
		return record, err
	}

	record.TxHash = txHash

	for i := 0; i <= r.Intn(10); i++ {
		record.Contents = append(record.Contents, types.Content{
			Digest:     simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 1, 50)),
			DigestAlgo: simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 1, 50)),
			URI:        simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 0, 50)),
			Meta:       simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 0, 50)),
		})
	}

	acc, _ := simtypes.RandomAcc(r, accs)
	record.Creator = acc.Address

	return record, nil
}
