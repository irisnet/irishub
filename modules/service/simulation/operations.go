package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	cosmossimappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irisnet/irismod/modules/service/keeper"
	"github.com/irisnet/irismod/modules/service/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgDefineService         = "op_weight_msg_define_service"
	OpWeightMsgBindService           = "op_weight_msg_bind_service"
	OpWeightMsgUpdateServiceBinding  = "op_weight_msg_update_service_binding"
	OpWeightMsgSetWithdrawAddress    = "op_weight_msg_set_withdraw_address"
	OpWeightMsgDisableServiceBinding = "op_weight_msg_disable_service_binding"
	OpWeightMsgEnableServiceBinding  = "op_weight_msg_enable_service_binding"
	OpWeightMsgRefundServiceDeposit  = "op_weight_msg_refund_service_deposit"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONMarshaler,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgDefineService         int
		weightMsgBindService           int
		weightMsgUpdateServiceBinding  int
		weightMsgSetWithdrawAddress    int
		weightMsgDisableServiceBinding int
		weightMsgEnableServiceBinding  int
		weightMsgRefundServiceDeposit  int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgDefineService, &weightMsgDefineService, nil,
		func(_ *rand.Rand) {
			weightMsgDefineService = DefaultWeightMsgDefineService
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgBindService, &weightMsgBindService, nil,
		func(_ *rand.Rand) {
			weightMsgBindService = DefaultWeightMsgBindService
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgUpdateServiceBinding, &weightMsgUpdateServiceBinding, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateServiceBinding = DefaultWeightMsgUpdateServiceBinding
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgSetWithdrawAddress, &weightMsgSetWithdrawAddress, nil,
		func(_ *rand.Rand) {
			weightMsgSetWithdrawAddress = DefaultWeightMsgSetWithdrawAddress
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgDisableServiceBinding, &weightMsgDisableServiceBinding, nil,
		func(_ *rand.Rand) {
			weightMsgDisableServiceBinding = DefaultWeightMsgDisableServiceBinding
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgEnableServiceBinding, &weightMsgEnableServiceBinding, nil,
		func(_ *rand.Rand) {
			weightMsgEnableServiceBinding = DefaultWeightMsgEnableServiceBinding
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgRefundServiceDeposit, &weightMsgRefundServiceDeposit, nil,
		func(_ *rand.Rand) {
			weightMsgRefundServiceDeposit = DefaultWeightMsgRefundServiceDeposit
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgDefineService,
			SimulateMsgDefineService(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgBindService,
			SimulateMsgBindService(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgUpdateServiceBinding,
			SimulateMsgUpdateServiceBinding(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgSetWithdrawAddress,
			SimulateMsgSetWithdrawAddress(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgDisableServiceBinding,
			SimulateMsgDisableServiceBinding(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgEnableServiceBinding,
			SimulateMsgEnableServiceBinding(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgRefundServiceDeposit,
			SimulateMsgRefundServiceDeposit(ak, bk, k),
		),
	}
}

// SimulateMsgDefineService generates a MsgDefineService with random values.
func SimulateMsgDefineService(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		simAccount, _ := simtypes.RandomAcc(r, accs)

		serviceName := simtypes.RandStringOfLength(r, 20)
		serviceDescription := simtypes.RandStringOfLength(r, 50)
		authorDescription := simtypes.RandStringOfLength(r, 50)
		tags := []string{simtypes.RandStringOfLength(r, 20), simtypes.RandStringOfLength(r, 20)}
		schemas := `{"input":{"type":"object"},"output":{"type":"object"}}`

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		msg := types.NewMsgDefineService(serviceName, serviceDescription, tags, simAccount.Address.String(), authorDescription, schemas)

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeEncodingConfig().TxConfig
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
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgBindService generates a MsgBindService with random values.
func SimulateMsgBindService(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		simAccount, _ := simtypes.RandomAcc(r, accs)

		serviceName := simtypes.RandStringOfLength(r, 20)
		deposit := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 1000000, 2000000)))))
		pricing := fmt.Sprintf(`{"price":"%d%s"`, simtypes.RandIntBetween(r, 100, 1000), sdk.DefaultBondDenom)
		qos := uint64(simtypes.RandIntBetween(r, 10, 100))
		// XXX
		options := "{...}"

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		msg := types.NewMsgBindService(serviceName, simAccount.Address.String(), deposit, pricing, qos, options, simAccount.Address.String())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeEncodingConfig().TxConfig
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
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgUpdateServiceBinding generates a MsgUpdateServiceBinding with random values.
func SimulateMsgUpdateServiceBinding(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		simAccount, _ := simtypes.RandomAcc(r, accs)

		serviceName := simtypes.RandStringOfLength(r, 20)
		deposit := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 100, 1000)))))
		pricing := fmt.Sprintf(`{"price":"%d%s"`, simtypes.RandIntBetween(r, 100, 1000), sdk.DefaultBondDenom)
		qos := uint64(simtypes.RandIntBetween(r, 10, 100))
		// XXX
		options := "{...}"

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		msg := types.NewMsgUpdateServiceBinding(serviceName, simAccount.Address.String(), deposit, pricing, qos, options, simAccount.Address.String())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeEncodingConfig().TxConfig
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
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgSetWithdrawAddress generates a MsgSetWithdrawAddress with random values.
func SimulateMsgSetWithdrawAddress(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		simAccount, _ := simtypes.RandomAcc(r, accs)
		withdrawalAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		msg := types.NewMsgSetWithdrawAddress(simAccount.Address.String(), withdrawalAccount.Address.String())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeEncodingConfig().TxConfig
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
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDisableServiceBinding generates a MsgDisableServiceBinding with random values.
func SimulateMsgDisableServiceBinding(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		simAccount, _ := simtypes.RandomAcc(r, accs)
		serviceName := simtypes.RandStringOfLength(r, 20)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		msg := types.NewMsgDisableServiceBinding(serviceName, simAccount.Address.String(), simAccount.Address.String())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeEncodingConfig().TxConfig
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
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgEnableServiceBinding generates a MsgEnableServiceBinding with random values.
func SimulateMsgEnableServiceBinding(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		simAccount, _ := simtypes.RandomAcc(r, accs)

		serviceName := simtypes.RandStringOfLength(r, 20)
		deposit := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 100, 1000)))))

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		msg := types.NewMsgEnableServiceBinding(serviceName, simAccount.Address.String(), deposit, simAccount.Address.String())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeEncodingConfig().TxConfig
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
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgRefundServiceDeposit generates a MsgRefundServiceDeposit with random values.
func SimulateMsgRefundServiceDeposit(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		simAccount, _ := simtypes.RandomAcc(r, accs)
		serviceName := simtypes.RandStringOfLength(r, 20)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		msg := types.NewMsgRefundServiceDeposit(serviceName, simAccount.Address.String(), simAccount.Address.String())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeEncodingConfig().TxConfig
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
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
