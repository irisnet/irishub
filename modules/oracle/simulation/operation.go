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

	"github.com/irisnet/irismod/modules/oracle/keeper"
	"github.com/irisnet/irismod/modules/oracle/types"
	irishelpers "github.com/irisnet/irismod/simapp/helpers"
)

const (
	OpWeightMsgCreateFeed = "op_weight_msg_create_feed"
	OpWeightMsgPauseFeed  = "op_weight_msg_pause_feed"
	OpWeightMsgStartFeed  = "op_weight_msg_start_feed"
	OpWeightMsgEditFeed   = "op_weight_msg_edit_feed"
)

func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONMarshaler,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper) simulation.WeightedOperations {

	var weightCreate, weightPause, weightStart, WeightEdit int

	appParams.GetOrGenerate(
		cdc, OpWeightMsgCreateFeed, &weightCreate, nil,
		func(_ *rand.Rand) {
			weightCreate = 50
		},
	)
	appParams.GetOrGenerate(
		cdc, OpWeightMsgPauseFeed, &weightPause, nil,
		func(_ *rand.Rand) {
			weightPause = 50
		},
	)
	appParams.GetOrGenerate(
		cdc, OpWeightMsgStartFeed, &weightStart, nil,
		func(_ *rand.Rand) {
			weightStart = 50
		},
	)
	appParams.GetOrGenerate(
		cdc, OpWeightMsgEditFeed, &WeightEdit, nil,
		func(_ *rand.Rand) {
			WeightEdit = 50
		},
	)
	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightCreate,
			SimulateCreateFeed(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightStart,
			SimulateStartFeed(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightPause,
			SimulatePauseFeed(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			WeightEdit,
			SimulateEditFeed(k, ak, bk),
		),
	}
}

func SimulateCreateFeed(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		feedName := simtypes.RandStringOfLength(r, 10)
		latestHistory := uint64(simtypes.RandIntBetween(r, 1, 100))
		description := simtypes.RandStringOfLength(r, 50)
		creator := simAccount.Address.String()
		serviceName, owner, err := GenServiceDefinition(r, k.GetServiceKeeper(), accs, ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFeed, "Failed to generate service definition"), nil, err
		}
		providers := GenServiceBindingsAndProviders(ctx, serviceName, owner, k.GetServiceKeeper(), accs, r, bk)
		if len(providers) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFeed, "Failed to generate service bindings"), nil, nil
		}
		input := `{"header":{},"body":{}}`
		timeout := int64(simtypes.RandIntBetween(r, 10, 100))
		srvFeeCap := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 2, 10))))}
		repeatedFrequency := uint64(100)
		aggregateFunc := GenAggregateFunc(r)
		valueJsonPath := `{"input":{"type":"object"},"output":{"type":"object"}}`
		responseThreshold := uint32(simtypes.RandIntBetween(r, 1, len(providers)))

		msg := &types.MsgCreateFeed{
			FeedName:          feedName,
			LatestHistory:     latestHistory,
			Description:       description,
			Creator:           creator,
			ServiceName:       serviceName,
			Providers:         providers,
			Input:             input,
			Timeout:           timeout,
			ServiceFeeCap:     srvFeeCap,
			RepeatedFrequency: repeatedFrequency,
			AggregateFunc:     aggregateFunc,
			ValueJsonPath:     valueJsonPath,
			ResponseThreshold: responseThreshold,
		}
		account := ak.GetAccount(ctx, simAccount.Address)

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFeed, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := irishelpers.GenTx(
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, nil
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeCreateFeed, err.Error()), nil, nil
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateStartFeed(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		feed := GenFeed(k, r, ctx)
		if feed.Size() == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStartFeed, "feed not exist"), nil, err
		}

		creator, err := sdk.AccAddressFromBech32(feed.Creator)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStartFeed, "creator invalid address"), nil, err
		}
		acc, found := simtypes.FindAccount(accs, creator)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStartFeed, "account not found"), nil, nil
		}
		account := ak.GetAccount(ctx, acc.Address)

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		msg := &types.MsgStartFeed{
			FeedName: feed.FeedName,
			Creator:  feed.Creator,
		}
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
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

func SimulatePauseFeed(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {

		feed := GenFeed(k, r, ctx)
		if feed.Size() == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPauseFeed, "feed not exist"), nil, err
		}

		creator, err := sdk.AccAddressFromBech32(feed.Creator)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPauseFeed, "creator invalid address"), nil, err
		}
		acc, found := simtypes.FindAccount(accs, creator)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPauseFeed, "account not found"), nil, nil
		}
		account := ak.GetAccount(ctx, acc.Address)

		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		msg := &types.MsgPauseFeed{
			FeedName: feed.FeedName,
			Creator:  feed.Creator,
		}
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
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

func SimulateEditFeed(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		feed := GenFeed(k, r, ctx)
		if feed.Size() == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEditFeed, "feed not exist"), nil, err
		}

		creator, err := sdk.AccAddressFromBech32(feed.Creator)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPauseFeed, "creator invalid address"), nil, err
		}
		acc, found := simtypes.FindAccount(accs, creator)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPauseFeed, "account not found"), nil, nil
		}

		description := simtypes.RandStringOfLength(r, 50)
		latestHistory := uint64(simtypes.RandIntBetween(r, 1, 100))
		providers := GenProviders(r, accs)
		timeout := int64(simtypes.RandIntBetween(r, 10, 100))
		srvFeeCap := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(simtypes.RandIntBetween(r, 2, 10))))}
		repeatedFrequency := uint64(100)
		responseThreshold := uint32(simtypes.RandIntBetween(r, 1, len(providers)))

		msg := &types.MsgEditFeed{
			FeedName:          feed.FeedName,
			Description:       description,
			LatestHistory:     latestHistory,
			Providers:         providers,
			Timeout:           timeout,
			ServiceFeeCap:     srvFeeCap,
			RepeatedFrequency: repeatedFrequency,
			ResponseThreshold: responseThreshold,
			Creator:           feed.Creator,
		}

		account := ak.GetAccount(ctx, acc.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateFeed, err.Error()), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := irishelpers.GenTx(
			r,
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, nil
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeCreateFeed, err.Error()), nil, nil
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil

	}
}

func GenAggregateFunc(r *rand.Rand) string {
	slice := []string{"max", "min", "avg"}
	return slice[r.Intn(len(slice))]
}

func GenServiceDefinition(r *rand.Rand, sk types.ServiceKeeper, accs []simtypes.Account, ctx sdk.Context) (string, string, error) {
	simAccount, _ := simtypes.RandomAcc(r, accs)
	serviceName := simtypes.RandStringOfLength(r, 20)
	serviceDescription := simtypes.RandStringOfLength(r, 50)
	authorDescription := simtypes.RandStringOfLength(r, 50)
	tags := []string{simtypes.RandStringOfLength(r, 20), simtypes.RandStringOfLength(r, 20)}
	schemas := `{"input":{"type":"object"},"output":{"type":"object"}}`
	if err := sk.AddServiceDefinition(ctx, serviceName, serviceDescription, tags, simAccount.Address, authorDescription, schemas); err != nil {
		return "", "", err
	}
	return serviceName, simAccount.Address.String(), nil
}

func GenServiceBindingsAndProviders(ctx sdk.Context, serviceName, owner string, sk types.ServiceKeeper, accs []simtypes.Account, r *rand.Rand, bk types.BankKeeper) (providers []string) {
	ownerAddr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return
	}
	spendable := bk.SpendableCoins(ctx, ownerAddr)
	if spendable.IsZero() {
		return
	}
	token := spendable[r.Intn(len(spendable))]
	if token.IsZero() {
		return
	}
	for i := 0; i < 10; i++ {
		provider, _ := simtypes.RandomAcc(r, accs)
		deposit := sdk.NewCoins(sdk.NewCoin(token.Denom, simtypes.RandomAmount(r, token.Amount)))
		pricing := fmt.Sprintf(`{"price":"%d%s"}`, simtypes.RandIntBetween(r, 100, 1000), sdk.DefaultBondDenom)
		qos := uint64(simtypes.RandIntBetween(r, 10, 100))
		options := "{}"
		err := sk.AddServiceBinding(ctx, serviceName, provider.Address, deposit, pricing, qos, options, ownerAddr)
		if err != nil {
			providers = append(providers, provider.Address.String())
		}
	}
	return
}

func GenFeed(k keeper.Keeper, r *rand.Rand, ctx sdk.Context) types.Feed {
	var feeds []types.Feed
	k.IteratorFeeds(ctx, func(feed types.Feed) {
		feeds = append(feeds, feed)
	})
	if len(feeds) == 0 {
		return types.Feed{}
	}
	return feeds[r.Intn(len(feeds))]
}

func GenProviders(r *rand.Rand, accs []simtypes.Account) (providers []string) {
	for i := 0; i < 10; i++ {
		simAcc, _ := simtypes.RandomAcc(r, accs)
		providers = append(providers, simAcc.Address.String())
	}
	return
}
