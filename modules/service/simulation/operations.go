package simulation

import (
	"fmt"
	"math/rand"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	cosmossimappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irisnet/irismod/modules/service/keeper"
	"github.com/irisnet/irismod/modules/service/types"
	irishelper "github.com/irisnet/irismod/simapp/helpers"
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
	OpWeightMsgCallService           = "op_weight_msg_call_service"
	OpWeightMsgRespondService        = "op_weight_msg_respond_service"
	OpWeightMsgStartRequestContext   = "op_weight_msg_start_request_context"
	OpWeightMsgPauseRequestContext   = "op_weight_msg_pause_request_context"
	OpWeightMsgKillRequestContext    = "op_weight_msg_kill_request_context"
	OpWeightMsgUpdateRequestContext  = "op_weight_msg_update_request_context"
	OpWeightMsgWithdrawEarnedFees    = "op_weight_msg_withdraw_earned_fees"
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
		weightMsgCallService           int
		weightMsgRespondService        int
		weightMsgPauseRequestContext   int
		weightMsgStartRequestContext   int
		weightMsgKillRequestContext    int
		weightMsgUpdateRequestContext  int
		weightMsgWithdrawEarnedFees    int
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

	appParams.GetOrGenerate(cdc, OpWeightMsgRefundServiceDeposit, &weightMsgRefundServiceDeposit, nil,
		func(_ *rand.Rand) {
			weightMsgRefundServiceDeposit = DefaultWeightMsgRefundServiceDeposit
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgCallService, &weightMsgCallService, nil,
		func(_ *rand.Rand) {
			weightMsgCallService = DefaultWeightMsgCallService
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgRespondService, &weightMsgRespondService, nil,
		func(_ *rand.Rand) {
			weightMsgRespondService = DefaultWeightMsgRespondService
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgPauseRequestContext, &weightMsgPauseRequestContext, nil,
		func(_ *rand.Rand) {
			weightMsgPauseRequestContext = DefaultWeightMsgPauseRequestContext
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgStartRequestContext, &weightMsgStartRequestContext, nil,
		func(_ *rand.Rand) {
			weightMsgStartRequestContext = DefaultWeightMsgStartRequestContext
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgKillRequestContext, &weightMsgKillRequestContext, nil,
		func(_ *rand.Rand) {
			weightMsgKillRequestContext = DefaultWeightMsgKillRequestContext
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgUpdateRequestContext, &weightMsgUpdateRequestContext, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateRequestContext = DefaultWeightMsgUpdateRequestContext
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgWithdrawEarnedFees, &weightMsgWithdrawEarnedFees, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawEarnedFees = DefaultWeightMsgWithdrawEarnedFees
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
		simulation.NewWeightedOperation(
			weightMsgCallService,
			SimulateMsgCallService(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgRespondService,
			SimulateMsgRespondService(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgPauseRequestContext,
			SimulateMsgPauseRequestContext(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgStartRequestContext,
			SimulateMsgStartRequestContext(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgKillRequestContext,
			SimulateMsgKillRequestContext(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgUpdateRequestContext,
			SimulateMsgUpdateRequestContext(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgWithdrawEarnedFees,
			SimulateMsgWithdrawEarnedFees(ak, bk, k),
		),
	}
}

// SimulateMsgDefineService generates a MsgDefineService with random values.
func SimulateMsgDefineService(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		simAccount, _ := simtypes.RandomAcc(r, accs)

		serviceName := simtypes.RandStringOfLength(r, 70)
		serviceDescription := simtypes.RandStringOfLength(r, 280)
		authorDescription := simtypes.RandStringOfLength(r, 280)
		tags := []string{simtypes.RandStringOfLength(r, 20), simtypes.RandStringOfLength(r, 20)}
		schemas := `{"input":{"type":"object"},"output":{"type":"object"}}`

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		msg := types.NewMsgDefineService(serviceName, serviceDescription, tags, simAccount.Address.String(), authorDescription, schemas)

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
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

		def := GenServiceDefinition(r, k, ctx)
		if def.Size() == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBindService, "def not exsit"), nil, nil
		}

		owner, err := sdk.AccAddressFromBech32(def.Author)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBindService, "invalid owner address"), nil, nil
		}

		acc, found := simtypes.FindAccount(accs, owner)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBindService, "account not found"), nil, nil
		}

		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, owner)

		pricing := fmt.Sprintf(`{"price":"%d%s"}`, simtypes.RandIntBetween(r, 10, 50), sdk.DefaultBondDenom)

		parsedPricing, err := k.ParsePricing(ctx, pricing)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBindService, err.Error()), nil, err
		}

		deposit, err := k.GetMinDeposit(ctx, parsedPricing)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBindService, "invalid minimum deposit"), nil, nil
		}

		// random provider address
		provider, _ := simtypes.RandomAcc(r, accs)

		currentOwner, found := k.GetOwner(ctx, provider.Address)
		if found && !owner.Equals(currentOwner) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBindService, "owner not matching"), nil, nil
		}

		if _, found := k.GetServiceBinding(ctx, def.Name, provider.Address); found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBindService, "service binding already exists"), nil, nil
		}

		qos := uint64(simtypes.RandIntBetween(r, 10, 100))
		options := "{}"
		msg := types.NewMsgBindService(def.Name, provider.Address.String(), deposit, pricing, qos, options, def.Author)

		spendable, hasNeg := spendable.SafeSub(deposit)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBindService, "Insufficient funds"), nil, nil
		}

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			acc.PrivKey,
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

		binding := GenServiceBinding(r, k, ctx)
		if binding.Size() == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateServiceBinding, "binding not exist"), nil, nil
		}
		owner, err := sdk.AccAddressFromBech32(binding.Owner)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateServiceBinding, "invalid owner address"), nil, nil
		}

		acc, found := simtypes.FindAccount(accs, owner)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateServiceBinding, "account not found"), nil, nil
		}
		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, owner)

		pricing := fmt.Sprintf(`{"price":"%d%s"}`, simtypes.RandIntBetween(r, 10, 50), sdk.DefaultBondDenom)
		qos := uint64(simtypes.RandIntBetween(r, 10, 100))
		options := "{}"

		parsedPricing, err := k.ParsePricing(ctx, pricing)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBindService, err.Error()), nil, err
		}

		deposit, err := k.GetMinDeposit(ctx, parsedPricing)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBindService, "invalid minimum deposit"), nil, nil
		}

		spendable, hasNeg := spendable.SafeSub(deposit)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateServiceBinding, "Insufficient funds"), nil, nil
		}

		msg := types.NewMsgUpdateServiceBinding(binding.ServiceName, binding.Provider, deposit, pricing, qos, options, acc.Address.String())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			acc.PrivKey,
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

		withdrawalAccount, _ := simtypes.RandomAcc(r, accs)

		binding := GenServiceBinding(r, k, ctx)
		if binding.Size() == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetWithdrawAddress, "binding not exist"), nil, nil
		}
		owner, err := sdk.AccAddressFromBech32(binding.Owner)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetWithdrawAddress, "invalid owner address"), nil, nil
		}

		acc, found := simtypes.FindAccount(accs, owner)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetWithdrawAddress, "account not found"), nil, nil
		}
		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, owner)

		msg := types.NewMsgSetWithdrawAddress(binding.Owner, withdrawalAccount.Address.String())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			acc.PrivKey,
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

		binding := GenServiceBinding(r, k, ctx)
		if binding.Size() == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDisableServiceBinding, "binding not exist"), nil, nil
		}
		if !binding.Available {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDisableServiceBinding, "binding is disabled"), nil, nil
		}
		owner, err := sdk.AccAddressFromBech32(binding.Owner)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDisableServiceBinding, "invalid owner address"), nil, nil
		}

		acc, found := simtypes.FindAccount(accs, owner)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDisableServiceBinding, "account not found"), nil, nil
		}
		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, owner)

		msg := types.NewMsgDisableServiceBinding(binding.ServiceName, binding.Provider, binding.Owner)

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			acc.PrivKey,
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

		binding := GenServiceBinding(r, k, ctx)
		if binding.Size() == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEnableServiceBinding, "binding not exist"), nil, nil
		}
		if binding.Available {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEnableServiceBinding, "binding is available"), nil, nil
		}
		owner, err := sdk.AccAddressFromBech32(binding.Owner)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEnableServiceBinding, "invalid owner address"), nil, nil
		}

		acc, found := simtypes.FindAccount(accs, owner)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEnableServiceBinding, "account not found"), nil, nil
		}
		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, owner)

		provider, err := sdk.AccAddressFromBech32(binding.Provider)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEnableServiceBinding, "invalid provider address"), nil, nil
		}
		pricing := k.GetPricing(ctx, binding.ServiceName, provider)
		deposit, err := k.GetMinDeposit(ctx, pricing)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBindService, "invalid minimum deposit"), nil, nil
		}

		spendable, hasNeg := spendable.SafeSub(deposit)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateServiceBinding, "Insufficient funds"), nil, nil
		}

		msg := types.NewMsgEnableServiceBinding(binding.ServiceName, binding.Provider, deposit, binding.Owner)

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			acc.PrivKey,
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

		binding := GenServiceBindingDisabled(r, k, ctx)
		if binding.Size() == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRefundServiceDeposit, "binding not exist"), nil, nil
		}

		owner, err := sdk.AccAddressFromBech32(binding.Owner)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRefundServiceDeposit, "invalid owner address"), nil, nil
		}

		acc, found := simtypes.FindAccount(accs, owner)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRefundServiceDeposit, "account not found"), nil, nil
		}
		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, owner)

		provider, err := sdk.AccAddressFromBech32(binding.Provider)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRefundServiceDeposit, "invalid provider address"), nil, nil
		}

		refundableTime := binding.DisabledTime.Add(k.ArbitrationTimeLimit(ctx)).Add(k.ComplaintRetrospect(ctx))

		currentTime := ctx.BlockHeader().Time
		if currentTime.Before(refundableTime) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRefundServiceDeposit, "invalid refundable time"), nil, nil
		}

		msg := types.NewMsgRefundServiceDeposit(binding.ServiceName, provider.String(), owner.String())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			acc.PrivKey,
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

// SimulateMsgCallService generates a MsgCallService with random values.
func SimulateMsgCallService(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)
		binding := GenServiceBinding(r, k, ctx)
		if binding.Size() == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCallService, "binding not exist"), nil, nil
		}
		definition, found := k.GetServiceDefinition(ctx, binding.ServiceName)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCallService, "serviceDefinition not exist"), nil, nil
		}
		providers := GetProviders(definition, k, ctx)
		if len(providers) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCallService, "providers not exist"), nil, nil
		}

		serviceName := binding.ServiceName
		consumer := simAccount.Address.String()
		input := `{"header":{},"body":{}}`
		serviceFeeCap := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 2, 10))))}
		timeout := int64(simtypes.RandIntBetween(r, 1, int(k.MaxRequestTimeout(ctx))))

		repeated := true
		repeatedFrequency := uint64(100)
		repeatedTotal := int64(10)

		msg := types.NewMsgCallService(serviceName, providers, consumer, input,
			serviceFeeCap, timeout, repeated, repeatedFrequency, repeatedTotal)

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		spendable, hasNeg := spendable.SafeSub(serviceFeeCap)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCallService, "Insufficient funds"), nil, nil
		}

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := irishelper.GenTx(
			r,
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

// SimulateMsgRespondService generates a MsgRespondService with random values.
func SimulateMsgRespondService(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		request := GenRequest(r, k, ctx)
		if request.Size() == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRespondService, "request is not exsit"), nil, nil
		}

		provider, err := sdk.AccAddressFromBech32(request.Provider)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRespondService, "invalid address"), nil, nil
		}

		result := `{"code":200,"message":""}`
		output := `{"header":{},"body":{}}`

		msg := types.NewMsgRespondService(request.Id, request.Provider, result, output)

		acc, found := simtypes.FindAccount(accs, provider)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRespondService, "account not found"), nil, nil
		}

		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			acc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, nil
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil

	}
}

// SimulateMsgPauseRequestContext generates a MsgSPauseRequestContext with random values.
func SimulateMsgPauseRequestContext(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// requestContext must be running
		requestContextId := GenRunningContextId(r, k, ctx)
		if len(requestContextId) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPauseRequestContext, "requestContextId not exist"), nil, nil
		}

		requestContext, found := k.GetRequestContext(ctx, requestContextId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPauseRequestContext, "requestContext not found"), nil, nil
		}
		if len(requestContext.ModuleName) > 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPauseRequestContext, "not authorized operation"), nil, nil
		}
		consumer, err := sdk.AccAddressFromBech32(requestContext.Consumer)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPauseRequestContext, "invalid address"), nil, nil
		}

		msg := types.NewMsgPauseRequestContext(requestContextId.String(), consumer.String())

		acc, found := simtypes.FindAccount(accs, consumer)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPauseRequestContext, "account not found"), nil, nil
		}

		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			acc.PrivKey,
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

// SimulateMsgStartRequestContext generates a MsgStartRequestContext with random values.
func SimulateMsgStartRequestContext(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		requestContextId := GenPausedRequestContextId(r, k, ctx)
		if len(requestContextId) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStartRequestContext, "requestContextId not exist"), nil, nil
		}

		requestContext, found := k.GetRequestContext(ctx, requestContextId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStartRequestContext, "requestContext not found"), nil, nil
		}

		if !requestContext.Repeated {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStartRequestContext, "requestContext non repeated"), nil, nil
		}

		if len(requestContext.ModuleName) > 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStartRequestContext, "not authorized operation"), nil, nil
		}
		consumer, err := sdk.AccAddressFromBech32(requestContext.Consumer)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStartRequestContext, "invalid address"), nil, nil
		}

		msg := types.NewMsgStartRequestContext(requestContextId.String(), consumer.String())

		acc, found := simtypes.FindAccount(accs, consumer)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStartRequestContext, "account not found"), nil, nil
		}

		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			acc.PrivKey,
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

// SimulateMsgKillRequestContext generates a MsgKillRequestContext with random values.
func SimulateMsgKillRequestContext(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		requestContextId := GenRequestContextId(r, k, ctx)
		if len(requestContextId) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgKillRequestContext, "requestContextId not exist"), nil, nil
		}

		requestContext, found := k.GetRequestContext(ctx, requestContextId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgKillRequestContext, "requestContext not found"), nil, nil
		}

		if !requestContext.Repeated {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgKillRequestContext, "requestContext non repeated"), nil, nil
		}
		if len(requestContext.ModuleName) > 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgKillRequestContext, "not authorized operation"), nil, nil
		}
		consumer, err := sdk.AccAddressFromBech32(requestContext.Consumer)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgKillRequestContext, "invalid address"), nil, nil
		}

		msg := types.NewMsgKillRequestContext(requestContextId.String(), consumer.String())

		acc, found := simtypes.FindAccount(accs, consumer)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgKillRequestContext, "account not found"), nil, nil
		}

		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			acc.PrivKey,
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

// SimulateMsgUpdateRequestContext generates a MsgUpdateRequestContext with random values.
func SimulateMsgUpdateRequestContext(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		requestContextId := GenRequestContextId(r, k, ctx)
		if len(requestContextId) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateRequestContext, "requestContextId not exist"), nil, nil
		}

		requestContext, found := k.GetRequestContext(ctx, requestContextId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateRequestContext, "request not found"), nil, nil
		}

		if requestContext.State == types.COMPLETED {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateRequestContext, "request context completed"), nil, nil
		}
		if len(requestContext.ModuleName) > 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateRequestContext, "not authorized operation"), nil, nil
		}
		consumer, err := sdk.AccAddressFromBech32(requestContext.Consumer)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateRequestContext, "invalid address"), nil, nil
		}

		definition, found := k.GetServiceDefinition(ctx, requestContext.ServiceName)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateRequestContext, "definition not found"), nil, nil
		}
		providers := GetProviders(definition, k, ctx)

		serviceFeeCap := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 2, 10))))}
		timeout := r.Int63n(k.MaxRequestTimeout(ctx))
		repeatedFrequency := uint64(0)
		repeatedTotal := int64(0)

		msg := types.NewMsgUpdateRequestContext(requestContextId.String(), providers, serviceFeeCap, timeout, repeatedFrequency, repeatedTotal, consumer.String())

		acc, found := simtypes.FindAccount(accs, consumer)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgKillRequestContext, "account not found"), nil, nil
		}

		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		spendable, hasNeg := spendable.SafeSub(serviceFeeCap)

		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCallService, "Insufficient funds"), nil, nil
		}

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			acc.PrivKey,
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

// SimulateMsgWithdrawEarnedFees generates a MsgWithdrawEarnedFees with random values.
func SimulateMsgWithdrawEarnedFees(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		binding := GenServiceBinding(r, k, ctx)
		if binding.Size() == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawEarnedFees, "binding not found"), nil, nil
		}

		owner, err := sdk.AccAddressFromBech32(binding.Owner)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawEarnedFees, "invalid address"), nil, nil
		}
		acc, found := simtypes.FindAccount(accs, owner)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWithdrawEarnedFees, "account not found"), nil, nil
		}

		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		msg := types.NewMsgWithdrawEarnedFees(binding.Owner, binding.Provider)

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := cosmossimappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			acc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err := app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, nil
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil

	}
}

// GenServiceDefinition randomized serviceDefinition
func GenServiceDefinition(r *rand.Rand, k keeper.Keeper, ctx sdk.Context) types.ServiceDefinition {
	var definitions []types.ServiceDefinition
	k.IterateServiceDefinitions(
		ctx,
		func(definition types.ServiceDefinition) bool {
			definitions = append(definitions, definition)
			return false
		},
	)
	if len(definitions) > 0 {
		return definitions[r.Intn(len(definitions))]
	}
	return types.ServiceDefinition{}
}

// GenServiceBinding randomized serviceBinding
func GenServiceBinding(r *rand.Rand, k keeper.Keeper, ctx sdk.Context) types.ServiceBinding {
	var bindings []types.ServiceBinding
	k.IterateServiceBindings(
		ctx,
		func(binding types.ServiceBinding) bool {
			bindings = append(bindings, binding)
			return false
		},
	)
	if len(bindings) > 0 {
		return bindings[r.Intn(len(bindings))]
	}
	return types.ServiceBinding{}
}

// GenServiceBindingDisabled randomized serviceBindingDisabled
func GenServiceBindingDisabled(r *rand.Rand, k keeper.Keeper, ctx sdk.Context) types.ServiceBinding {
	var bindings []types.ServiceBinding
	k.IterateServiceBindings(
		ctx,
		func(binding types.ServiceBinding) bool {
			if !binding.Available {
				bindings = append(bindings, binding)
				return false
			}
			return false
		},
	)
	if len(bindings) > 0 {
		return bindings[r.Intn(len(bindings))]
	}
	return types.ServiceBinding{}
}

// GenRequestContextId randomized requestContext
func GenRequestContextId(r *rand.Rand, k keeper.Keeper, ctx sdk.Context) tmbytes.HexBytes {
	var requestIds []tmbytes.HexBytes
	k.IterateRequestContexts(
		ctx,
		func(requestContextID tmbytes.HexBytes, requestContext types.RequestContext) bool {
			requestIds = append(requestIds, requestContextID)
			return false
		},
	)
	if len(requestIds) > 0 {
		return requestIds[r.Intn(len(requestIds))]
	}
	return tmbytes.HexBytes{}
}

// GenRequest randomized request
func GenRequest(r *rand.Rand, k keeper.Keeper, ctx sdk.Context) types.Request {
	requestContextId := GenContextId(r, k, ctx)
	requestContext, found := k.GetRequestContext(ctx, requestContextId)
	if !found {
		return types.Request{}
	}
	var providerAddrs []sdk.AccAddress
	if len(requestContext.Providers) == 0 {
		return types.Request{}
	}
	for _, provider := range requestContext.Providers {
		providerAddr, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return types.Request{}
		}
		providerAddrs = append(providerAddrs, providerAddr)
	}
	providerRequests := make(map[string][]string)
	requestIds := k.InitiateRequests(ctx, requestContextId, providerAddrs, providerRequests)

	if len(requestIds) == 0 {
		return types.Request{}
	}
	requestId := requestIds[r.Intn(len(requestIds))]
	request, found := k.GetRequest(ctx, requestId)
	if !found {
		return types.Request{}
	}
	return request
}

// GenContextId randomized contextId
func GenContextId(r *rand.Rand, k keeper.Keeper, ctx sdk.Context) tmbytes.HexBytes {
	var requestIds []tmbytes.HexBytes
	k.IterateRequestContexts(
		ctx,
		func(requestContextID tmbytes.HexBytes, requestContext types.RequestContext) bool {
			requestIds = append(requestIds, requestContextID)
			return false
		},
	)
	if len(requestIds) > 0 {
		return requestIds[r.Intn(len(requestIds))]
	}
	return tmbytes.HexBytes{}
}

// GenRunningContextId randomized runningContextId
func GenRunningContextId(r *rand.Rand, k keeper.Keeper, ctx sdk.Context) tmbytes.HexBytes {
	var requestIds []tmbytes.HexBytes
	k.IterateRequestContexts(
		ctx,
		func(requestContextID tmbytes.HexBytes, requestContext types.RequestContext) bool {
			if requestContext.State == types.RUNNING {
				requestIds = append(requestIds, requestContextID)
				return false
			}
			return false
		},
	)
	if len(requestIds) > 0 {
		return requestIds[r.Intn(len(requestIds))]
	}
	return tmbytes.HexBytes{}
}

// GenPausedRequestContextId randomized pausedRequestContextId
func GenPausedRequestContextId(r *rand.Rand, k keeper.Keeper, ctx sdk.Context) tmbytes.HexBytes {
	var requestIds []tmbytes.HexBytes
	k.IterateRequestContexts(
		ctx,
		func(requestContextID tmbytes.HexBytes, requestContext types.RequestContext) bool {
			if requestContext.State == types.PAUSED {
				requestIds = append(requestIds, requestContextID)
				return false
			}
			return false
		},
	)
	if len(requestIds) > 0 {
		return requestIds[r.Intn(len(requestIds))]
	}
	return tmbytes.HexBytes{}
}

func GetProviders(definition types.ServiceDefinition, k keeper.Keeper, ctx sdk.Context) (providers []string) {
	if definition.Size() == 0 {
		return
	}

	owner, err := sdk.AccAddressFromBech32(definition.Author)
	if err != nil {
		return
	}

	bindings := k.GetOwnerServiceBindings(ctx, owner, definition.Name)
	if len(bindings) > 0 {
		for _, binding := range bindings {
			providers = append(providers, binding.Provider)
		}
	}
	return
}
