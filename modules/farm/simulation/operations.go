package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/farm/keeper"
	"github.com/irisnet/irismod/modules/farm/types"

	"github.com/cosmos/cosmos-sdk/codec"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreatePool   = "op_weight_msg_create_pool"
	OpWeightMsgAppendReward = "op_weight_msg_append_reward"
	OpWeightMsgStake        = "op_weight_msg_stake"
	OpWeightMsgUnStake      = "op_weight_msg_unStake"
	OpWeightMsgHarvest      = "op_weight_msg_harvest"
	OpWeightMsgDestroyPool  = "op_weight_msg_destroy_pool"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONMarshaler,
	k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simulation.WeightedOperations {

	var (
		weightMsgCreatePool   int
		weightMsgAppendReward int
		weightMsgStake        int
		weightMsgUnStake      int
		weightMsgHarvest      int
		weightMsgDestroyPool  int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgCreatePool, &weightMsgCreatePool, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePool = 30
		},
	)
	appParams.GetOrGenerate(cdc, OpWeightMsgAppendReward, &weightMsgAppendReward, nil,
		func(_ *rand.Rand) {
			weightMsgAppendReward = 30
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
			weightMsgAppendReward,
			SimulateMsgAppendReward(k, ak, bk),
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

func SimulateMsgCreatePool(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		if spendable.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "spendable is zero"), nil, nil
		}

		_, hasNeg := spendable.SafeSub(sdk.NewCoins(k.CreatePoolFee(ctx)))

		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "Insufficient funds"), nil, nil
		}

		totalReward := GenTotalReward(r, spendable)
		lpTokenDenom := GenLpToken(r, spendable)
		editable := GenDestructible(r)
		startHeight := GenStartHeight(r, ctx)
		rewardPerBlock := GenRewardPerBlock(r, totalReward)

		if rewardPerBlock.Amount.LT(sdk.ZeroInt()) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "rewardPerBlock less than zeroInt"), nil, nil
		}

		if totalReward.Amount.LT(rewardPerBlock.Amount) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "totalReward less than rewardPerBlock"), nil, nil
		}

		balance, hasNeg := spendable.SafeSub(sdk.NewCoins(totalReward))

		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "Insufficient funds"), nil, nil
		}

		name := GenFarmPoolName(r)
		if _, exist := k.GetPool(ctx, name); exist {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreatePool, "farm pool is exist"), nil, nil
		}

		msg := &types.MsgCreatePool{
			Name:           name,
			Description:    GenDescription(r),
			LpTokenDenom:   lpTokenDenom.Denom,
			StartHeight:    startHeight,
			RewardPerBlock: sdk.Coins{sdk.NewCoin(rewardPerBlock.Denom, rewardPerBlock.Amount)},
			TotalReward:    sdk.NewCoins(totalReward),
			Editable:       editable,
			Creator:        simAccount.Address.String(),
		}

		fees, err := simtypes.RandomFees(r, ctx, balance)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeCreatePool, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
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

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, nil
		}
		keeper.RewardInvariant(k)
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil

	}
}

func SimulateMsgAppendReward(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		farmPool, exist := genRandomFarmPool(ctx, k, r)
		if !exist {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAdjustPool, "farm pool is not exist"), nil, nil
		}

		if k.Expired(ctx, farmPool) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAdjustPool, "farmPool has expired"), nil, nil
		}

		creator, err := sdk.AccAddressFromBech32(farmPool.Creator)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAdjustPool, "invalid address"), nil, err
		}

		simAccount, found := simtypes.FindAccount(accs, creator)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAdjustPool, "unable to find account"), nil, nil
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		if spendable.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAdjustPool, "Insufficient funds"), nil, nil
		}

		rules := k.GetRewardRules(ctx, farmPool.Name)
		amount := GenAppendReward(r, rules, spendable)
		if amount.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAdjustPool, "Insufficient funds"), nil, nil
		}

		// Need to subtract the appendReward balance
		balance, hasNeg := spendable.SafeSub(amount)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAdjustPool, "Insufficient funds"), nil, nil
		}

		msg := &types.MsgAdjustPool{
			PoolName:         farmPool.Name,
			AdditionalReward: amount,
			Creator:          farmPool.Creator,
		}

		fees, err := simtypes.RandomFees(r, ctx, balance)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeAppendReward, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
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

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}
		keeper.RewardInvariant(k)
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateMsgStake(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		farmPool, exist := genRandomFarmPool(ctx, k, r)
		if !exist {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStake, "farm pool is not exist"), nil, nil
		}

		if farmPool.StartHeight > ctx.BlockHeight() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStake, "the farm activity has not yet started"), nil, nil
		}

		if k.Expired(ctx, farmPool) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStake, "the farm activity has ended"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		if spendable.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStake, "spendable is zero"), nil, nil
		}

		amount := GenStake(r, farmPool, spendable)
		if amount.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStake, "The sender does not have the specified lpToken"), nil, nil
		}

		// Need to subtract the stake balance
		balance, hasNeg := spendable.SafeSub(sdk.Coins{amount})
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStake, "Insufficient funds"), nil, nil
		}

		msg := &types.MsgStake{
			PoolName: farmPool.Name,
			Amount:   amount,
			Sender:   account.GetAddress().String(),
		}

		fees, err := simtypes.RandomFees(r, ctx, balance)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeStake, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
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

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		keeper.RewardInvariant(k)
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateMsgUnStake(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)

		farmInfo, exist := genRandomFarmInfo(ctx, k, r, account.GetAddress().String())
		if !exist {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnstake, "farmer not found in pool"), nil, nil
		}

		farmPool, exist := k.GetPool(ctx, farmInfo.PoolName)
		if !exist {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnstake, "farm pool is not exist"), nil, nil
		}

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		if spendable.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnstake, "spendable is zero"), nil, nil
		}

		unStake := GenUnStake(r, farmPool, farmInfo)
		if unStake.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnstake, "the sender does not have the specified lpToken"), nil, nil
		}

		if farmInfo.Locked.LT(unStake.Amount) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnstake, "the lp token unStaked must be less than staked"), nil, nil
		}

		amount := farmInfo.Locked
		msg := &types.MsgUnstake{
			PoolName: farmPool.Name,
			Amount:   sdk.NewCoin(farmPool.TotalLpTokenLocked.Denom, amount),
			Sender:   account.GetAddress().String(),
		}

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeUnstake, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
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

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}
		keeper.RewardInvariant(k)
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
func SimulateMsgHarvest(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		simAccount, _ := simtypes.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)

		farmInfo, exist := genRandomFarmInfo(ctx, k, r, account.GetAddress().String())
		if !exist {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnstake, "farmer not found in pool"), nil, nil
		}

		farmPool, exist := k.GetPool(ctx, farmInfo.PoolName)
		if !exist {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUnstake, "farm pool is not exist"), nil, nil
		}

		if k.Expired(ctx, farmPool) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgHarvest, "farm pool has expired"), nil, nil
		}

		msg := &types.MsgHarvest{
			PoolName: farmPool.Name,
			Sender:   account.GetAddress().String(),
		}

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeHarvest, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
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

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}
		keeper.RewardInvariant(k)
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
func SimulateMsgDestroyPool(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		farmPool, exist := genRandomFarmPool(ctx, k, r)
		if !exist {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStake, "farm pool is not exist"), nil, nil
		}

		creator, err := sdk.AccAddressFromBech32(farmPool.Creator)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAdjustPool, "invalid address"), nil, err
		}

		simAccount, found := simtypes.FindAccount(accs, creator)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAdjustPool, "unable to find account"), nil, nil
		}

		if !farmPool.Editable {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDestroyPool, "farm pool is not editable"), nil, nil
		}

		if k.Expired(ctx, farmPool) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDestroyPool, "farm pool has expired"), nil, nil
		}

		msg := &types.MsgDestroyPool{
			PoolName: farmPool.Name,
			Creator:  simAccount.Address.String(),
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		if spendable.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAdjustPool, "Insufficient funds"), nil, nil
		}

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeDestroyPool, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
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

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}
		keeper.RewardInvariant(k)
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// GenTotalReward randomized totalReward
func GenTotalReward(r *rand.Rand, spendableCoin sdk.Coins) sdk.Coin {
	token := spendableCoin[r.Intn(len(spendableCoin))]
	return sdk.NewCoin(token.Denom, simtypes.RandomAmount(r, token.Amount))
}

// GenLpToken randomized lpToken
func GenLpToken(r *rand.Rand, spendableCoin sdk.Coins) sdk.Coin {
	token := spendableCoin[r.Intn(len(spendableCoin))]
	return sdk.NewCoin(token.Denom, simtypes.RandomAmount(r, token.Amount))
}

// GenStartHeight randomized startHeight
func GenStartHeight(r *rand.Rand, ctx sdk.Context) int64 {
	curHeight := int(ctx.BlockHeight())
	return int64(r.Intn(curHeight) + curHeight)
}

// GenRewardPerBlock randomized rewardPerBlock
func GenRewardPerBlock(r *rand.Rand, coin sdk.Coin) sdk.Coin {
	return sdk.NewCoin(coin.Denom, simtypes.RandomAmount(r, coin.Amount))
}

// GenRewardRule randomized rewardRule
func GenRewardRule(r *rand.Rand, rules types.RewardRules) types.RewardRule {
	return rules[r.Intn(len(rules))]
}

// GenAppendReward randomized appendReward
func GenAppendReward(r *rand.Rand, rules types.RewardRules, spendable sdk.Coins) sdk.Coins {
	rule := GenRewardRule(r, rules)
	for _, coin := range spendable {
		if coin.Denom != rule.Reward {
			break
		}
		return sdk.Coins{sdk.NewCoin(coin.Denom, simtypes.RandomAmount(r, coin.Amount))}
	}
	return sdk.Coins{}
}

// GenStake randomized stake
func GenStake(r *rand.Rand, pool types.FarmPool, spendable sdk.Coins) sdk.Coin {
	for _, coin := range spendable {
		if coin.Denom != pool.TotalLpTokenLocked.Denom {
			break
		}
		return sdk.NewCoin(pool.TotalLpTokenLocked.Denom, simtypes.RandomAmount(r, coin.Amount))
	}
	return sdk.NewCoin(pool.TotalLpTokenLocked.Denom, sdk.ZeroInt())
}

// GenUnStake randomized unStake
func GenUnStake(r *rand.Rand, pool types.FarmPool, info types.FarmInfo) sdk.Coin {
	return sdk.NewCoin(pool.TotalLpTokenLocked.Denom, simtypes.RandomAmount(r, info.Locked))
}

// GenFarmPoolName randomized farmPoolName
func GenFarmPoolName(r *rand.Rand) string {
	return simtypes.RandStringOfLength(r, 10)
}

// GenDestructible randomized editable
func GenDestructible(r *rand.Rand) bool {
	return r.Int()%2 == 0
}

// GenDescription randomized editable
func GenDescription(r *rand.Rand) string {
	return simtypes.RandStringOfLength(r, 100)
}

// genRandomFarmPool randomized farmPoolName
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
func genRandomFarmInfo(ctx sdk.Context, k keeper.Keeper, r *rand.Rand, addr string) (types.FarmInfo, bool) {
	var farmInfos []types.FarmInfo

	k.IteratorFarmInfo(ctx, addr, func(farmInfo types.FarmInfo) {
		farmInfos = append(farmInfos, farmInfo)
	})

	if len(farmInfos) > 0 {
		return farmInfos[r.Intn(len(farmInfos))], true
	}
	return types.FarmInfo{}, false
}
