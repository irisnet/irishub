package simulation

import (
	"errors"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irisnet/irismod/farm/keeper"
	"github.com/irisnet/irismod/farm/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreatePool  = "op_weight_msg_create_pool"
	OpWeightMsgAdjustPool  = "op_weight_msg_adjust_pool"
	OpWeightMsgStake       = "op_weight_msg_stake"
	OpWeightMsgUnStake     = "op_weight_msg_unStake"
	OpWeightMsgHarvest     = "op_weight_msg_harvest"
	OpWeightMsgDestroyPool = "op_weight_msg_destroy_pool"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec,
	k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper,
) simulation.WeightedOperations {
	var (
		weightMsgCreatePool  int
		weightMsgAdjustPool  int
		weightMsgStake       int
		weightMsgUnStake     int
		weightMsgHarvest     int
		weightMsgDestroyPool int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgCreatePool, &weightMsgCreatePool, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePool = 30
		},
	)
	appParams.GetOrGenerate(cdc, OpWeightMsgAdjustPool, &weightMsgAdjustPool, nil,
		func(_ *rand.Rand) {
			weightMsgAdjustPool = 30
		},
	)
	appParams.GetOrGenerate(cdc, OpWeightMsgStake, &weightMsgStake, nil,
		func(_ *rand.Rand) {
			weightMsgStake = 50
		},
	)
	appParams.GetOrGenerate(cdc, OpWeightMsgStake, &weightMsgStake, nil,
		func(_ *rand.Rand) {
			weightMsgStake = 50
		},
	)
	appParams.GetOrGenerate(cdc, OpWeightMsgUnStake, &weightMsgUnStake, nil,
		func(_ *rand.Rand) {
			weightMsgUnStake = 50
		},
	)
	appParams.GetOrGenerate(cdc, OpWeightMsgHarvest, &weightMsgHarvest, nil,
		func(_ *rand.Rand) {
			weightMsgHarvest = 40
		},
	)
	appParams.GetOrGenerate(cdc, OpWeightMsgDestroyPool, &weightMsgDestroyPool, nil,
		func(_ *rand.Rand) {
			weightMsgDestroyPool = 30
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreatePool,
			SimulateMsgCreatePool(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgAdjustPool,
			SimulateMsgAdjustPool(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgStake,
			SimulateMsgStake(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgUnStake,
			SimulateMsgUnStake(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgHarvest,
			SimulateMsgHarvest(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgDestroyPool,
			SimulateMsgDestroyPool(k, ak, bk),
		),
	}
}

func SimulateMsgCreatePool(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error,
	) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		if spendable.IsZero() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgCreatePool,
				"spendable is zero",
			), nil, nil
		}

		_, hasNeg := spendable.SafeSub(sdk.NewCoins(k.CreatePoolFee(ctx))...)

		if hasNeg {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgCreatePool,
				"Insufficient funds",
			), nil, nil
		}

		totalReward, err := GenTotalReward(r, spendable)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgCreatePool,
				"Insufficient funds",
			), nil, nil
		}

		lpTokenDenom, err := GenLpToken(r, spendable)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgCreatePool,
				"Insufficient funds",
			), nil, nil
		}
		editable := GenDestructible(r)
		startHeight := GenStartHeight(r, ctx)
		rewardPerBlock, err := GenRewardPerBlock(r, totalReward)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgCreatePool,
				"Insufficient funds",
			), nil, nil
		}

		if rewardPerBlock.Amount.LTE(sdk.ZeroInt()) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgCreatePool,
				"rewardPerBlock less than zeroInt",
			), nil, nil
		}

		if totalReward.Amount.LTE(rewardPerBlock.Amount) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgCreatePool,
				"totalReward less than rewardPerBlock",
			), nil, nil
		}

		balance, hasNeg := spendable.SafeSub(sdk.NewCoins(totalReward)...)

		if hasNeg {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgCreatePool,
				"Insufficient funds",
			), nil, nil
		}

		msg := &types.MsgCreatePool{
			Description:    GenDescription(r),
			LptDenom:       lpTokenDenom.Denom,
			StartHeight:    startHeight,
			RewardPerBlock: sdk.Coins{sdk.NewCoin(rewardPerBlock.Denom, rewardPerBlock.Amount)},
			TotalReward:    sdk.NewCoins(totalReward),
			Editable:       editable,
			Creator:        simAccount.Address.String(),
		}

		fees, err := simtypes.RandomFees(r, ctx, balance)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.EventTypeCreatePool,
				err.Error(),
			), nil, err
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
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

		if _, _, err = app.SimDeliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, nil
		}
		keeper.RewardInvariant(k)
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil

	}
}

func SimulateMsgAdjustPool(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error,
	) {
		farmPool, exist := genRandomFarmPool(ctx, k, r)
		if !exist {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"farm pool is not exist",
			), nil, nil
		}

		if k.Expired(ctx, farmPool) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"farmPool has expired",
			), nil, nil
		}

		if !farmPool.Editable {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"farmPool is not editable",
			), nil, nil
		}

		creator, err := sdk.AccAddressFromBech32(farmPool.Creator)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"invalid address",
			), nil, err
		}

		simAccount, found := simtypes.FindAccount(accs, creator)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"unable to find account",
			), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		if spendable.IsZero() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"insufficient funds",
			), nil, nil
		}

		rules := k.GetRewardRules(ctx, farmPool.Id)
		rewardPerBlock, err := GenRewardPerBlock(r, spendable[r.Intn(len(spendable))])
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"Insufficient funds",
			), nil, nil
		}
		if rewardPerBlock.Amount.IsZero() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"insufficient funds",
			), nil, nil
		}

		if rewardPerBlock.Denom != GenRewardRule(r, rules).Reward {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"invalid reward",
			), nil, nil
		}

		amount, err := GenAppendReward(r, rules, spendable)
		if err != nil || amount.IsZero() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"insufficient funds",
			), nil, nil
		}

		// Need to subtract the appendReward balance
		balance, hasNeg := spendable.SafeSub(amount...)
		if hasNeg {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"insufficient funds",
			), nil, nil
		}

		msg := &types.MsgAdjustPool{
			PoolId:           farmPool.Id,
			AdditionalReward: amount,
			RewardPerBlock:   sdk.Coins{sdk.NewCoin(rewardPerBlock.Denom, rewardPerBlock.Amount)},
			Creator:          farmPool.Creator,
		}

		fees, err := simtypes.RandomFees(r, ctx, balance)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.EventTypeAppendReward,
				err.Error(),
			), nil, err
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
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

		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}
		keeper.RewardInvariant(k)
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

func SimulateMsgStake(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error,
	) {
		farmPool, exist := genRandomFarmPool(ctx, k, r)
		if !exist {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgStake,
				"farm pool is not exist",
			), nil, nil
		}

		if farmPool.StartHeight > ctx.BlockHeight() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgStake,
				"the farm activity has not yet started",
			), nil, nil
		}

		if k.Expired(ctx, farmPool) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgStake,
				"the farm activity has ended",
			), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		if spendable.IsZero() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgStake,
				"spendable is zero",
			), nil, nil
		}

		amount, err := GenStake(r, farmPool, spendable)
		if err != nil || amount.IsZero() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgStake,
				"The sender does not have the specified lpToken",
			), nil, nil
		}

		// Need to subtract the stake balance
		balance, hasNeg := spendable.SafeSub(sdk.Coins{amount}...)
		if hasNeg {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgStake,
				"Insufficient funds",
			), nil, nil
		}

		msg := &types.MsgStake{
			PoolId: farmPool.Id,
			Amount: amount,
			Sender: account.GetAddress().String(),
		}

		fees, err := simtypes.RandomFees(r, ctx, balance)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeStake, err.Error()), nil, err
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
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

		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		keeper.RewardInvariant(k)
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

func SimulateMsgUnStake(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error,
	) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)

		farmInfo, exist := genRandomFarmInfo(ctx, k, r, account.GetAddress().String())
		if !exist {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUnstake,
				"farmer not found in pool",
			), nil, nil
		}

		farmPool, exist := k.GetPool(ctx, farmInfo.PoolId)
		if !exist {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUnstake,
				"farm pool is not exist",
			), nil, nil
		}

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		if spendable.IsZero() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUnstake,
				"spendable is zero",
			), nil, nil
		}

		unStake, err := GenUnStake(r, farmPool, farmInfo)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUnstake,
				"Insufficient funds",
			), nil, nil
		}
		if unStake.IsZero() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUnstake,
				"the sender does not have the specified lpToken",
			), nil, nil
		}

		if farmInfo.Locked.LT(unStake.Amount) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUnstake,
				"the lp token unStaked must be less than staked",
			), nil, nil
		}

		amount := farmInfo.Locked
		msg := &types.MsgUnstake{
			PoolId: farmPool.Id,
			Amount: sdk.NewCoin(farmPool.TotalLptLocked.Denom, amount),
			Sender: account.GetAddress().String(),
		}

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeUnstake, err.Error()), nil, err
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
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

		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}
		keeper.RewardInvariant(k)
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

func SimulateMsgHarvest(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error,
	) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)

		farmInfo, exist := genRandomFarmInfo(ctx, k, r, account.GetAddress().String())
		if !exist {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUnstake,
				"farmer not found in pool",
			), nil, nil
		}

		farmPool, exist := k.GetPool(ctx, farmInfo.PoolId)
		if !exist {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUnstake,
				"farm pool is not exist",
			), nil, nil
		}

		if k.Expired(ctx, farmPool) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgHarvest,
				"farm pool has expired",
			), nil, nil
		}

		msg := &types.MsgHarvest{
			PoolId: farmPool.Id,
			Sender: account.GetAddress().String(),
		}

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeHarvest, err.Error()), nil, err
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
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

		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}
		keeper.RewardInvariant(k)
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

func SimulateMsgDestroyPool(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		simtypes.OperationMsg, []simtypes.FutureOperation, error,
	) {
		farmPool, exist := genRandomFarmPool(ctx, k, r)
		if !exist {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgStake,
				"farm pool is not exist",
			), nil, nil
		}

		creator, err := sdk.AccAddressFromBech32(farmPool.Creator)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"invalid address",
			), nil, err
		}

		simAccount, found := simtypes.FindAccount(accs, creator)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"unable to find account",
			), nil, nil
		}

		if !farmPool.Editable {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgDestroyPool,
				"farm pool is not editable",
			), nil, nil
		}

		if k.Expired(ctx, farmPool) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgDestroyPool,
				"farm pool has expired",
			), nil, nil
		}

		msg := &types.MsgDestroyPool{
			PoolId:  farmPool.Id,
			Creator: simAccount.Address.String(),
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		if spendable.IsZero() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAdjustPool,
				"Insufficient funds",
			), nil, nil
		}

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.EventTypeDestroyPool,
				err.Error(),
			), nil, err
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
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

		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}
		keeper.RewardInvariant(k)
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// GenTotalReward randomized totalReward
func GenTotalReward(r *rand.Rand, spendableCoin sdk.Coins) (sdk.Coin, error) {
	token := spendableCoin[r.Intn(len(spendableCoin))]
	amount, err := simtypes.RandPositiveInt(r, token.Amount)
	if err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(token.Denom, amount), nil
}

// GenLpToken randomized lpToken
func GenLpToken(r *rand.Rand, spendableCoin sdk.Coins) (sdk.Coin, error) {
	token := spendableCoin[r.Intn(len(spendableCoin))]
	amount, err := simtypes.RandPositiveInt(r, token.Amount)
	if err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(token.Denom, amount), nil
}

// GenStartHeight randomized startHeight
func GenStartHeight(r *rand.Rand, ctx sdk.Context) int64 {
	curHeight := int(ctx.BlockHeight())
	return int64(r.Intn(curHeight) + curHeight)
}

// GenRewardPerBlock randomized rewardPerBlock
func GenRewardPerBlock(r *rand.Rand, coin sdk.Coin) (sdk.Coin, error) {
	amount, err := simtypes.RandPositiveInt(r, coin.Amount)
	if err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(coin.Denom, amount), nil
}

// GenRewardRule randomized rewardRule
func GenRewardRule(r *rand.Rand, rules types.RewardRules) types.RewardRule {
	return rules[r.Intn(len(rules))]
}

// GenAppendReward randomized appendReward
func GenAppendReward(
	r *rand.Rand,
	rules types.RewardRules,
	spendable sdk.Coins,
) (sdk.Coins, error) {
	rule := GenRewardRule(r, rules)
	for _, coin := range spendable {
		if coin.Denom == rule.Reward {
			amount, err := simtypes.RandPositiveInt(r, coin.Amount)
			if err != nil {
				return nil, err
			}
			return sdk.NewCoins(sdk.NewCoin(coin.Denom, amount)), nil
		}
	}
	return nil, errors.New("no spendable token")
}

// GenStake randomized stake
func GenStake(r *rand.Rand, pool types.FarmPool, spendable sdk.Coins) (sdk.Coin, error) {
	for _, coin := range spendable {
		if coin.Denom == pool.TotalLptLocked.Denom {
			amount, err := simtypes.RandPositiveInt(r, coin.Amount)
			if err != nil {
				return sdk.Coin{}, err
			}
			return sdk.NewCoin(pool.TotalLptLocked.Denom, amount), nil
		}
	}
	return sdk.NewCoin(pool.TotalLptLocked.Denom, sdk.ZeroInt()), nil
}

// GenUnStake randomized unStake
func GenUnStake(r *rand.Rand, pool types.FarmPool, info types.FarmInfo) (sdk.Coin, error) {
	amount, err := simtypes.RandPositiveInt(r, info.Locked)
	if err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(pool.TotalLptLocked.Denom, amount), nil
}

// GenDestructible randomized editable
func GenDestructible(r *rand.Rand) bool {
	return r.Int()%2 == 0
}

// GenDescription randomized editable
func GenDescription(r *rand.Rand) string {
	return simtypes.RandStringOfLength(r, 100)
}

// genRandomFarmPool randomized farmpoolId
func genRandomFarmPool(ctx sdk.Context, k keeper.Keeper, r *rand.Rand) (types.FarmPool, bool) {
	var pools []types.FarmPool
	k.IteratorAllPools(ctx, func(pool types.FarmPool) {
		pools = append(pools, pool)
	})
	if len(pools) > 0 {
		return pools[r.Intn(len(pools))], true
	}
	return types.FarmPool{}, false
}

// genRandomFarmInfo randomized farmInfo
func genRandomFarmInfo(
	ctx sdk.Context,
	k keeper.Keeper,
	r *rand.Rand,
	addr string,
) (types.FarmInfo, bool) {
	var farmInfos []types.FarmInfo

	k.IteratorFarmInfo(ctx, addr, func(farmInfo types.FarmInfo) {
		farmInfos = append(farmInfos, farmInfo)
	})

	if len(farmInfos) > 0 {
		return farmInfos[r.Intn(len(farmInfos))], true
	}
	return types.FarmInfo{}, false
}
