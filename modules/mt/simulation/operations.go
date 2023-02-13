package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/irisnet/irismod/modules/mt/keeper"
	mt "github.com/irisnet/irismod/modules/mt/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgIssueDenom    = "op_weight_msg_issue_denom"
	OpWeightMsgMintMT        = "op_weight_msg_mint_mt"
	OpWeightMsgEditMT        = "op_weight_msg_edit_mt_tokenData"
	OpWeightMsgTransferMT    = "op_weight_msg_transfer_mt"
	OpWeightMsgBurnMT        = "op_weight_msg_transfer_burn_mt"
	OpWeightMsgTransferDenom = "op_weight_msg_transfer_denom"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONCodec,
	k keeper.Keeper,
	ak mt.AccountKeeper,
	bk mt.BankKeeper,
) simulation.WeightedOperations {
	var weightIssueDenom, weightMint, weightEdit, weightBurn, weightTransfer, weightTransferDenom int

	appParams.GetOrGenerate(
		cdc, OpWeightMsgIssueDenom, &weightIssueDenom, nil,
		func(_ *rand.Rand) {
			weightIssueDenom = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgMintMT, &weightMint, nil,
		func(_ *rand.Rand) {
			weightMint = 100
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgEditMT, &weightEdit, nil,
		func(_ *rand.Rand) {
			weightEdit = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgTransferMT, &weightTransfer, nil,
		func(_ *rand.Rand) {
			weightTransfer = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgBurnMT, &weightBurn, nil,
		func(_ *rand.Rand) {
			weightBurn = 10
		},
	)
	appParams.GetOrGenerate(
		cdc, OpWeightMsgTransferDenom, &weightTransferDenom, nil,
		func(_ *rand.Rand) {
			weightTransferDenom = 10
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightTransfer,
			SimulateMsgTransferMT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightIssueDenom,
			SimulateMsgIssueDenom(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMint,
			SimulateMsgMintMT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightEdit,
			SimulateMsgEditMT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightBurn,
			SimulateMsgBurnMT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightTransferDenom,
			SimulateMsgTransferDenom(k, ak, bk),
		),
	}
}

// SimulateMsgIssueDenom simulates issue an denom
func SimulateMsgIssueDenom(k keeper.Keeper, ak mt.AccountKeeper, bk mt.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		sender, _ := simtypes.RandomAcc(r, accs)
		denomName := simtypes.RandStringOfLength(r, 10)
		denomData := simtypes.RandStringOfLength(r, 10)

		senderAcc := ak.GetAccount(ctx, sender.Address)
		spendableCoins := bk.SpendableCoins(ctx, sender.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeIssueDenom, err.Error()), nil, err
		}

		spendLimit := spendableCoins.Sub(fees...)
		if spendLimit == nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeIssueDenom, "spend limit is nil"), nil, nil
		}

		msg := &mt.MsgIssueDenom{
			Name:   denomName,
			Data:   []byte(denomData),
			Sender: sender.Address.String(),
		}

		txCfg := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenSignedMockTx(
			r,
			txCfg,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{senderAcc.GetAccountNumber()},
			[]uint64{senderAcc.GetSequence()},
			sender.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeIssueDenom, "unable to generate mock tx"), nil, err
		}

		_, _, err = app.SimDeliver(txCfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, sdk.MsgTypeURL(msg), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgMintMT simulates mint an MT
func SimulateMsgMintMT(k keeper.Keeper, ak mt.AccountKeeper, bk mt.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		collection, ok := randCollection(ctx, r, k)
		if !ok {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeMintMT, "not fetch a collection"), nil, nil
		}

		mtr, denomID, ok := randMTWithCollection(ctx, collection, r, k)
		if !ok {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeMintMT, "not fetch an mt"), nil, nil
		}

		owner := collection.Denom.Owner
		senderAddr, err := sdk.AccAddressFromBech32(owner)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeMintMT, err.Error()), nil, err
		}

		senderAcc := ak.GetAccount(ctx, senderAddr)
		sender, ok := simtypes.FindAccount(accs, senderAcc.GetAddress())
		if !ok {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeMintMT, "owner(sender) not found"), nil, err
		}
		recipient, _ := simtypes.RandomAcc(r, accs)

		spendableCoins := bk.SpendableCoins(ctx, sender.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeMintMT, err.Error()), nil, err
		}

		spendLimit := spendableCoins.Sub(fees...)
		if spendLimit == nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeMintMT, "spend limit is nil"), nil, nil
		}

		msg := &mt.MsgMintMT{
			Id:        mtr.Id,
			DenomId:   denomID,
			Amount:    uint64(simtypes.RandIntBetween(r, 1, 100)),
			Data:      nil,
			Sender:    sender.Address.String(),
			Recipient: recipient.Address.String(),
		}

		txCfg := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenSignedMockTx(
			r,
			txCfg,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{senderAcc.GetAccountNumber()},
			[]uint64{senderAcc.GetSequence()},
			sender.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeEditMT, "unable to generate mock tx"), nil, err
		}

		_, _, err = app.SimDeliver(txCfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, sdk.MsgTypeURL(msg), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgEditMT simulates an edit tokenData transaction
func SimulateMsgEditMT(k keeper.Keeper, ak mt.AccountKeeper, bk mt.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		data := simtypes.RandStringOfLength(r, 10)

		mtr, denomID, owner, ok := randMT(ctx, r, k)
		if !ok {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeEditMT, "not fetch an mt"), nil, nil
		}
		ownerAddr := sdk.MustAccAddressFromBech32(owner)
		senderAcc := ak.GetAccount(ctx, ownerAddr)
		spendableCoins := bk.SpendableCoins(ctx, ownerAddr)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeEditMT, err.Error()), nil, err
		}

		spendLimit := spendableCoins.Sub(fees...)
		if spendLimit == nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeEditMT, "spend limit is nil"), nil, nil
		}

		amt := k.GetBalance(ctx, denomID, mtr.Id, ownerAddr)
		if amt == 0 {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeEditMT, "sender doesn't have this mt"), nil, nil
		}

		msg := &mt.MsgEditMT{
			Id:      mtr.Id,
			DenomId: denomID,
			Data:    []byte(data),
			Sender:  ownerAddr.String(),
		}

		sender, ok := simtypes.FindAccount(accs, senderAcc.GetAddress())
		if !ok {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeEditMT, "owner(sender) not found"), nil, nil
		}

		txCfg := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenSignedMockTx(
			r,
			txCfg,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{senderAcc.GetAccountNumber()},
			[]uint64{senderAcc.GetSequence()},
			sender.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeEditMT, "unable to generate mock tx"), nil, err
		}

		_, _, err = app.SimDeliver(txCfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, sdk.MsgTypeURL(msg), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgTransferMT simulates the transfer of an MT
func SimulateMsgTransferMT(k keeper.Keeper, ak mt.AccountKeeper, bk mt.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error,
	) {

		mtr, denomID, owner, ok := randMT(ctx, r, k)
		if !ok {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransfer, "not fetch an mt"), nil, nil
		}

		ownerAddr := sdk.MustAccAddressFromBech32(owner)
		senderAcc := ak.GetAccount(ctx, ownerAddr)
		sender, ok := simtypes.FindAccount(accs, senderAcc.GetAddress())
		if !ok {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransfer, "owner(sender) not found"), nil, nil
		}

		recipient, _ := simtypes.RandomAcc(r, accs)
		if sender.Address.Equals(recipient.Address) {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransfer, "sender and recipient are same"), nil, nil
		}

		spendableCoins := bk.SpendableCoins(ctx, sender.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransfer, err.Error()), nil, err
		}

		spendLimit := spendableCoins.Sub(fees...)
		if spendLimit == nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransfer, "spend limit is nil"), nil, nil
		}

		amt := k.GetBalance(ctx, denomID, mtr.Id, sender.Address)
		if amt <= 1 {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransfer, "sender doesn't have enough mt balances "), nil, nil
		}

		amt = uint64(simtypes.RandIntBetween(r, 1, int(amt)))
		msg := &mt.MsgTransferMT{
			Id:        mtr.Id,
			DenomId:   denomID,
			Amount:    amt,
			Sender:    sender.Address.String(),
			Recipient: recipient.Address.String(),
		}

		txCfg := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenSignedMockTx(
			r,
			txCfg,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{senderAcc.GetAccountNumber()},
			[]uint64{senderAcc.GetSequence()},
			sender.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransfer, "unable to generate mock tx"), nil, err
		}

		_, _, err = app.SimDeliver(txCfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, sdk.MsgTypeURL(msg), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgBurnMT simulates a burn of an existing MT
func SimulateMsgBurnMT(k keeper.Keeper, ak mt.AccountKeeper, bk mt.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		mtr, denomID, owner, ok := randMT(ctx, r, k)
		if !ok {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeBurnMT, "not fetch an mt"), nil, nil
		}

		ownerAddr := sdk.MustAccAddressFromBech32(owner)
		senderAcc := ak.GetAccount(ctx, ownerAddr)
		sender, ok := simtypes.FindAccount(accs, senderAcc.GetAddress())
		if !ok {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeBurnMT, "not fetch an mt"), nil, nil
		}

		spendableCoins := bk.SpendableCoins(ctx, sender.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeBurnMT, err.Error()), nil, err
		}

		spendLimit := spendableCoins.Sub(fees...)
		if spendLimit == nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeBurnMT, "spend limit is nil"), nil, nil
		}

		amt := k.GetBalance(ctx, denomID, mtr.Id, sender.Address)
		if amt <= 1 {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeBurnMT, "sender doesn't have enough mt balances "), nil, nil
		}

		amt = uint64(simtypes.RandIntBetween(r, 1, int(amt))) // unsafe conversion
		msg := &mt.MsgBurnMT{
			Id:      mtr.Id,
			DenomId: denomID,
			Amount:  amt,
			Sender:  sender.Address.String(),
		}

		txCfg := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenSignedMockTx(
			r,
			txCfg,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{senderAcc.GetAccountNumber()},
			[]uint64{senderAcc.GetSequence()},
			sender.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeBurnMT, "unable to generate mock tx"), nil, err
		}

		_, _, err = app.SimDeliver(txCfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, sdk.MsgTypeURL(msg), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgTransferDenom simulates the transfer of a Denom
func SimulateMsgTransferDenom(k keeper.Keeper, ak mt.AccountKeeper, bk mt.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		collection, ok := randCollection(ctx, r, k)
		if !ok {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransferDenom, "collection not found"), nil, err
		}
		owner := collection.Denom.Owner
		senderAddr, err := sdk.AccAddressFromBech32(owner)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransferDenom, err.Error()), nil, err
		}

		senderAcc := ak.GetAccount(ctx, senderAddr)
		sender, ok := simtypes.FindAccount(accs, senderAcc.GetAddress())
		if !ok {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransferDenom, "owner(sender) not found"), nil, err
		}
		recipient, _ := simtypes.RandomAcc(r, accs)

		spendableCoins := bk.SpendableCoins(ctx, sender.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransferDenom, err.Error()), nil, err
		}

		spendLimit := spendableCoins.Sub(fees...)
		if spendLimit == nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransferDenom, "spend limit is nil"), nil, nil
		}

		msg := &mt.MsgTransferDenom{
			Id:        collection.Denom.Id,
			Sender:    sender.Address.String(),
			Recipient: recipient.Address.String(),
		}

		txCfg := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenSignedMockTx(
			r,
			txCfg,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{senderAcc.GetAccountNumber()},
			[]uint64{senderAcc.GetSequence()},
			sender.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, mt.EventTypeTransferDenom, "unable to generate mock tx"), nil, err
		}

		_, _, err = app.SimDeliver(txCfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(mt.ModuleName, sdk.MsgTypeURL(msg), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// randCollection randomly returns a Collection
func randCollection(ctx sdk.Context, r *rand.Rand, k keeper.Keeper) (mt.Collection, bool) {
	collection := mt.Collection{}

	denoms := k.GetDenoms(ctx)
	if len(denoms) == 0 {
		return collection, false
	}

	denom := denoms[r.Intn(len(denoms))]
	mts := k.GetMTs(ctx, denom.Id)
	rmts := make([]mt.MT, len(mts))
	for i := 0; i < len(mts); i++ {
		var ok bool
		if rmts[i], ok = mts[i].(mt.MT); !ok {
			return collection, false
		}
	}

	collection.Denom = &denom
	collection.Mts = rmts

	return collection, true
}

// randMT randomly returns an MT
func randMT(ctx sdk.Context, r *rand.Rand, k keeper.Keeper) (mt.MT, string, string, bool) {
	collection, ok := randCollection(ctx, r, k)
	if !ok {
		return mt.MT{}, "", "", false
	}

	mts := collection.Mts
	if len(mts) == 0 {
		return mt.MT{}, "", "", false
	}

	idx := r.Intn(len(mts))

	return mts[idx], collection.Denom.Id, collection.Denom.Owner, true
}

// randMTWithCollection randomly returns an MT but with collection specified.
func randMTWithCollection(ctx sdk.Context, collection mt.Collection, r *rand.Rand, k keeper.Keeper) (mt.MT, string, bool) {
	if collection.Denom == nil || len(collection.Mts) == 0 {
		return mt.MT{}, "", false
	}

	idx := r.Intn(len(collection.Mts))
	return collection.Mts[idx], collection.Denom.Id, true
}
