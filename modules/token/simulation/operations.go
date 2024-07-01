package simulation

import (
	"fmt"
	"math/rand"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"mods.irisnet.org/modules/token/keeper"
	"mods.irisnet.org/modules/token/types"
	v1 "mods.irisnet.org/modules/token/types/v1"
)

// Simulation operation weights constants
const (
	OpWeightMsgIssueToken         = "op_weight_msg_issue_token"
	OpWeightMsgEditToken          = "op_weight_msg_edit_token"
	OpWeightMsgMintToken          = "op_weight_msg_mint_token"
	OpWeightMsgTransferTokenOwner = "op_weight_msg_transfer_token_owner"
	OpWeightMsgBurnToken          = "op_weight_msg_burn_token"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONCodec,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simulation.WeightedOperations {
	var weightIssue, weightEdit, weightMint, weightTransfer, weightBurn int
	appParams.GetOrGenerate(
		cdc, OpWeightMsgIssueToken, &weightIssue, nil,
		func(_ *rand.Rand) {
			weightIssue = 10
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgEditToken, &weightEdit, nil,
		func(_ *rand.Rand) {
			weightEdit = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgMintToken, &weightMint, nil,
		func(_ *rand.Rand) {
			weightMint = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgTransferTokenOwner, &weightTransfer, nil,
		func(_ *rand.Rand) {
			weightTransfer = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgBurnToken, &weightBurn, nil,
		func(_ *rand.Rand) {
			weightBurn = 50
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightIssue,
			SimulateIssueToken(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightEdit,
			SimulateEditToken(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMint,
			SimulateMintToken(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightTransfer,
			SimulateTransferTokenOwner(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightBurn,
			SimulateBurnToken(k, ak, bk),
		),
	}
}

// SimulateIssueToken tests and runs a single msg issue a new token
func SimulateIssueToken(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainID string,
	) (
		simtypes.OperationMsg,
		[]simtypes.FutureOperation,
		error,
	) {
		token, maxFees := genToken(ctx, r, k, ak, bk, accs)
		msg := v1.NewMsgIssueToken(
			token.Symbol,
			token.MinUnit,
			token.Name,
			token.Scale,
			token.InitialSupply,
			token.MaxSupply,
			token.Mintable,
			token.GetOwner().String(),
		)

		simAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(
					types.ModuleName,
					msg.Type(),
					fmt.Sprintf("account %s not found", token.Owner),
				), nil, fmt.Errorf(
					"account %s not found",
					token.Owner,
				)
		}

		owner, _ := sdk.AccAddressFromBech32(msg.Owner)
		account := ak.GetAccount(ctx, owner)
		fees, err := simtypes.RandomFees(r, ctx, maxFees)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				msg.Type(),
				"unable to generate fees",
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate issue token", nil), nil, nil
	}
}

// SimulateEditToken tests and runs a single msg edit a existed token
func SimulateEditToken(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		token, _, skip := selectOneToken(ctx, k, ak, bk, false)
		if skip {
			return simtypes.NoOpMsg(
				types.ModuleName,
				v1.TypeMsgEditToken,
				"skip edit token",
			), nil, nil
		}
		msg := v1.NewMsgEditToken(
			token.GetName(),
			token.GetSymbol(),
			token.GetMaxSupply(),
			types.True,
			token.GetOwner().String(),
		)

		simAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(
					types.ModuleName,
					msg.Type(),
					fmt.Sprintf("account %s not found", token.GetOwner()),
				), nil, fmt.Errorf(
					"account %s not found",
					token.GetOwner(),
				)
		}

		owner, _ := sdk.AccAddressFromBech32(msg.Owner)
		account := ak.GetAccount(ctx, owner)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				msg.Type(),
				"unable to generate fees",
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate edit token", nil), nil, nil
	}
}

// SimulateMintToken tests and runs a single msg mint a existed token
func SimulateMintToken(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		token, maxFee, skip := selectOneToken(ctx, k, ak, bk, true)
		if skip {
			return simtypes.NoOpMsg(
				types.ModuleName,
				v1.TypeMsgMintToken,
				"skip mint token",
			), nil, nil
		}
		simToAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &v1.MsgMintToken{
			Coin: sdk.Coin{
				Denom:  token.GetMinUnit(),
				Amount: sdkmath.NewIntWithDecimal(100, int(token.GetScale())),
			},
			Receiver: simToAccount.Address.String(),
			Owner:    token.GetOwner().String(),
		}

		ownerAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(
					types.ModuleName,
					msg.Type(),
					fmt.Sprintf("account %s not found", token.GetOwner()),
				), nil, fmt.Errorf(
					"account %s not found",
					token.GetOwner(),
				)
		}

		owner, _ := sdk.AccAddressFromBech32(msg.Owner)
		account := ak.GetAccount(ctx, owner)
		fees, err := simtypes.RandomFees(r, ctx, maxFee)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				msg.Type(),
				"unable to generate fees",
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
			ownerAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				msg.Type(),
				"unable to generate mock tx",
			), nil, err
		}

		if _, _, err = app.SimDeliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate mint token", nil), nil, nil
	}
}

// SimulateTransferTokenOwner tests and runs a single msg transfer to others
func SimulateTransferTokenOwner(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		token, _, skip := selectOneToken(ctx, k, ak, bk, false)
		if skip {
			return simtypes.NoOpMsg(
				types.ModuleName,
				v1.TypeMsgTransferTokenOwner,
				"skip TransferTokenOwner",
			), nil, nil
		}
		simToAccount, _ := simtypes.RandomAcc(r, accs)
		for simToAccount.Address.Equals(token.GetOwner()) {
			simToAccount, _ = simtypes.RandomAcc(r, accs)
		}

		msg := v1.NewMsgTransferTokenOwner(
			token.GetOwner().String(),
			simToAccount.Address.String(),
			token.GetSymbol(),
		)

		simAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(
					types.ModuleName,
					msg.Type(),
					fmt.Sprintf("account %s not found", token.GetOwner()),
				), nil, fmt.Errorf(
					"account %s not found",
					token.GetOwner(),
				)
		}

		srcOwner, _ := sdk.AccAddressFromBech32(msg.SrcOwner)
		account := ak.GetAccount(ctx, srcOwner)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				msg.Type(),
				"unable to generate fees",
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate transfer token", nil), nil, nil
	}
}

// SimulateBurnToken tests and runs a single msg burn a existed token
func SimulateBurnToken(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		token, _, skip := selectOneToken(ctx, k, ak, bk, false)
		if skip {
			return simtypes.NoOpMsg(
				types.ModuleName,
				v1.TypeMsgTransferTokenOwner,
				"skip burnToken",
			), nil, nil
		}

		owner, _ := sdk.AccAddressFromBech32(token.GetOwner().String())
		account := ak.GetAccount(ctx, owner)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		amount := spendable.AmountOf(token.GetMinUnit())
		if !amount.IsPositive() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				v1.TypeMsgBurnToken,
				"Insufficient funds",
			), nil, nil
		}

		amount2 := simtypes.RandomAmount(r, amount)

		spendable, hasNeg := spendable.SafeSub(
			sdk.Coins{sdk.NewCoin(token.GetMinUnit(), amount2)}...)
		if hasNeg {
			return simtypes.NoOpMsg(
				types.ModuleName,
				v1.TypeMsgBurnToken,
				"Insufficient funds",
			), nil, nil
		}

		msg := &v1.MsgBurnToken{
			Coin: sdk.Coin{
				Denom:  token.GetMinUnit(),
				Amount: amount2,
			},
			Sender: token.GetOwner().String(),
		}

		ownerAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(
					types.ModuleName,
					msg.Type(),
					fmt.Sprintf("account %s not found", token.GetOwner()),
				), nil, fmt.Errorf(
					"account %s not found",
					token.GetOwner(),
				)
		}

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				msg.Type(),
				"unable to generate fees",
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
			ownerAccount.PrivKey,
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

		return simtypes.NewOperationMsg(msg, true, "simulate mint token", nil), nil, nil
	}
}

func selectOneToken(
	ctx sdk.Context,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	mint bool,
) (token v1.TokenI, maxFees sdk.Coins, skip bool) {
	tokens := k.GetTokens(ctx, nil)
	if len(tokens) == 0 {
		return token, maxFees, true
	}

	for _, t := range tokens {
		if t.GetSymbol() == v1.GetNativeToken().Symbol {
			continue
		}
		if !mint {
			return t, nil, false
		}

		mintFee, err := k.GetTokenMintFee(ctx, t.GetSymbol())
		if err != nil {
			panic(err)
		}

		account := ak.GetAccount(ctx, t.GetOwner())
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		spendableStake := spendable.AmountOf(v1.GetNativeToken().MinUnit)
		if spendableStake.IsZero() || spendableStake.LT(mintFee.Amount) {
			continue
		}
		maxFees = sdk.NewCoins(
			sdk.NewCoin(v1.GetNativeToken().MinUnit, spendableStake).Sub(mintFee),
		)
		token = t
		return
	}
	return token, maxFees, true
}

func randStringBetween(r *rand.Rand, min, max int) string {
	strLen := simtypes.RandIntBetween(r, min, max)
	randStr := simtypes.RandStringOfLength(r, strLen)
	return randStr
}

func genToken(ctx sdk.Context,
	r *rand.Rand,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	accs []simtypes.Account,
) (v1.Token, sdk.Coins) {
	var token v1.Token
	token = randToken(r, accs)

	for k.HasToken(ctx, token.Symbol) {
		token = randToken(r, accs)
	}

	issueFee, err := k.GetTokenIssueFee(ctx, token.Symbol)
	if err != nil {
		panic(err)
	}

	account, maxFees := filterAccount(ctx, r, ak, bk, accs, issueFee)
	token.Owner = account.String()

	return token, maxFees
}

func filterAccount(
	ctx sdk.Context,
	r *rand.Rand,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	accs []simtypes.Account, fee sdk.Coin,
) (owner sdk.AccAddress, maxFees sdk.Coins) {
loop:
	simAccount, _ := simtypes.RandomAcc(r, accs)
	account := ak.GetAccount(ctx, simAccount.Address)
	spendable := bk.SpendableCoins(ctx, account.GetAddress())
	spendableStake := spendable.AmountOf(v1.GetNativeToken().MinUnit)
	if spendableStake.IsZero() || spendableStake.LT(fee.Amount) {
		goto loop
	}
	owner = account.GetAddress()
	maxFees = sdk.NewCoins(sdk.NewCoin(v1.GetNativeToken().MinUnit, spendableStake).Sub(fee))
	return
}

func randToken(r *rand.Rand, accs []simtypes.Account) v1.Token {
	var symbol, minUint string
	for {
		symbol = randStringBetween(r, types.MinimumSymbolLen, types.MaximumSymbolLen)
		minUint = symbol
		if err := types.ValidateSymbol(symbol); err == nil {
			break
		}
	}

	name := randStringBetween(r, 1, types.MaximumNameLen)
	scale := simtypes.RandIntBetween(r, 1, int(types.MaximumScale))
	initialSupply := r.Int63n(int64(types.MaximumInitSupply))
	maxSupply := 2 * initialSupply
	simAccount, _ := simtypes.RandomAcc(r, accs)

	return v1.Token{
		Symbol:        strings.ToLower(symbol),
		Name:          name,
		Scale:         uint32(scale),
		MinUnit:       strings.ToLower(minUint),
		InitialSupply: uint64(initialSupply),
		MaxSupply:     uint64(maxSupply),
		Mintable:      true,
		Owner:         simAccount.Address.String(),
	}
}
