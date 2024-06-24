package simulation

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"mods.irisnet.org/nft/keeper"
	"mods.irisnet.org/nft/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgIssueDenom    = "op_weight_msg_issue_denom"
	OpWeightMsgMintNFT       = "op_weight_msg_mint_nft"
	OpWeightMsgEditNFT       = "op_weight_msg_edit_nft_tokenData"
	OpWeightMsgTransferNFT   = "op_weight_msg_transfer_nft"
	OpWeightMsgBurnNFT       = "op_weight_msg_transfer_burn_nft"
	OpWeightMsgTransferDenom = "op_weight_msg_transfer_denom"
)

var (
	data = []string{
		"{\"key1\":\"value1\",\"key2\":\"value2\"}",
		"{\"irismod:key1\":\"value1\",\"irismod:key2\":\"value2\"}",
	}
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONCodec,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simulation.WeightedOperations {
	var weightIssueDenom, weightMint, weightEdit, weightBurn, weightTransfer, weightTransferDenom int

	appParams.GetOrGenerate(
		cdc, OpWeightMsgIssueDenom, &weightIssueDenom, nil,
		func(_ *rand.Rand) {
			weightIssueDenom = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgMintNFT, &weightMint, nil,
		func(_ *rand.Rand) {
			weightMint = 100
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgEditNFT, &weightEdit, nil,
		func(_ *rand.Rand) {
			weightEdit = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgTransferNFT, &weightTransfer, nil,
		func(_ *rand.Rand) {
			weightTransfer = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgBurnNFT, &weightBurn, nil,
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
			weightIssueDenom,
			SimulateMsgIssueDenom(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMint,
			SimulateMsgMintNFT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightEdit,
			SimulateMsgEditNFT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightTransfer,
			SimulateMsgTransferNFT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightBurn,
			SimulateMsgBurnNFT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightTransferDenom,
			SimulateMsgTransferDenom(k, ak, bk),
		),
	}
}

// SimulateMsgTransferNFT simulates the transfer of an NFT
func SimulateMsgTransferNFT(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		ownerAddr, denom, nftID := randNFT(ctx, k, r, false, true)
		if ownerAddr.Empty() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.EventTypeTransfer,
				"empty account",
			), nil, err
		}

		recipientAccount, _ := simtypes.RandomAcc(r, accs)
		msg := types.NewMsgTransferNFT(
			nftID,
			denom,
			"",
			"",
			"",
			randData(r),                       // tokenData
			ownerAddr.String(),                // sender
			recipientAccount.Address.String(), // recipient
		)
		account := ak.GetAccount(ctx, ownerAddr)

		ownerAccount, found := simtypes.FindAccount(accs, ownerAddr)
		if !found {
			err = fmt.Errorf("account %s not found", msg.Sender)
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.EventTypeTransfer,
				err.Error(),
			), nil, err
		}

		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.EventTypeTransfer,
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
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.EventTypeTransfer,
				err.Error(),
			), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgEditNFT simulates an edit tokenData transaction
func SimulateMsgEditNFT(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		ownerAddr, denom, nftID := randNFT(ctx, k, r, false, true)
		if ownerAddr.Empty() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.EventTypeEditNFT,
				"empty account",
			), nil, err
		}

		msg := types.NewMsgEditNFT(
			nftID,
			denom,
			"",
			simtypes.RandStringOfLength(r, 45), // tokenURI
			simtypes.RandStringOfLength(r, 32), // tokenURI
			randData(r),                        // tokenData
			ownerAddr.String(),
		)

		account := ak.GetAccount(ctx, ownerAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeEditNFT, err.Error()), nil, err
		}

		ownerAccount, found := simtypes.FindAccount(accs, ownerAddr)
		if !found {
			err = fmt.Errorf("account %s not found", ownerAddr)
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeEditNFT, err.Error()), nil, err
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
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeEditNFT, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgMintNFT simulates a mint of an NFT
func SimulateMsgMintNFT(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		randomSender, _ := simtypes.RandomAcc(r, accs)
		randomRecipient, _ := simtypes.RandomAcc(r, accs)

		msg := types.NewMsgMintNFT(
			genNFTID(r, idMinLen, idMaxLen),   // nft ID
			randDenom(ctx, k, r, true, false), // denom
			"",
			simtypes.RandStringOfLength(r, 45), // tokenURI
			simtypes.RandStringOfLength(r, 32), // uriHash
			randData(r),                        // tokenData
			randomSender.Address.String(),      // sender
			randomRecipient.Address.String(),   // recipient
		)

		account := ak.GetAccount(ctx, randomSender.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeMintNFT, err.Error()), nil, err
		}

		simAccount, found := simtypes.FindAccount(accs, randomSender.Address)
		if !found {
			err = fmt.Errorf("account %s not found", msg.Sender)
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeMintNFT, err.Error()), nil, err
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
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeMintNFT, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgBurnNFT simulates a burn of an existing NFT
func SimulateMsgBurnNFT(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		ownerAddr, denom, nftID := randNFT(ctx, k, r, false, false)
		if ownerAddr.Empty() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.EventTypeBurnNFT,
				"empty account",
			), nil, err
		}

		msg := types.NewMsgBurnNFT(ownerAddr.String(), nftID, denom)

		account := ak.GetAccount(ctx, ownerAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeBurnNFT, err.Error()), nil, err
		}

		simAccount, found := simtypes.FindAccount(accs, ownerAddr)
		if !found {
			err = fmt.Errorf("account %s not found", msg.Sender)
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeBurnNFT, err.Error()), nil, err
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
			return simtypes.NoOpMsg(types.ModuleName, types.EventTypeEditNFT, err.Error()), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgTransferDenom simulates the transfer of an denom
func SimulateMsgTransferDenom(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {

		denomID := randDenom(ctx, k, r, false, false)
		denom, err := k.GetDenomInfo(ctx, denomID)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgTransferDenom,
				err.Error(),
			), nil, err
		}

		creator, err := sdk.AccAddressFromBech32(denom.Creator)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgTransferDenom,
				err.Error(),
			), nil, err
		}
		account := ak.GetAccount(ctx, creator)
		owner, found := simtypes.FindAccount(accs, account.GetAddress())
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgTransferDenom,
				"creator not found",
			), nil, nil
		}

		recipient, _ := simtypes.RandomAcc(r, accs)
		msg := types.NewMsgTransferDenom(
			denomID,
			denom.Creator,
			recipient.Address.String(),
		)

		spendable := bk.SpendableCoins(ctx, owner.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgTransferDenom,
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
			owner.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				msg.Type(),
				"unable to generate mock tx",
			), nil, err
		}

		if _, _, err = app.SimDeliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.EventTypeTransfer,
				err.Error(),
			), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgIssueDenom simulates issue an denom
func SimulateMsgIssueDenom(
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {

		denomID := genDenomID(r)
		if k.HasDenom(ctx, denomID) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgTransferDenom,
				"denom exist",
			), nil, nil
		}

		sender, _ := simtypes.RandomAcc(r, accs)
		msg := types.NewMsgIssueDenom(
			denomID,
			strings.ToLower(simtypes.RandStringOfLength(r, 10)),
			"Schema",
			sender.Address.String(),
			simtypes.RandStringOfLength(r, 5),
			genRandomBool(r),
			genRandomBool(r),
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 32),
			randData(r),
		)
		account := ak.GetAccount(ctx, sender.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgTransferDenom,
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
			sender.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				msg.Type(),
				"unable to generate mock tx",
			), nil, err
		}

		if _, _, err = app.SimDeliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.EventTypeTransfer,
				err.Error(),
			), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

func randNFT(
	ctx sdk.Context,
	k keeper.Keeper,
	r *rand.Rand,
	mintable, editable bool,
) (sdk.AccAddress, string, string) {
	var denoms = []string{kitties, doggos}
	res, err := k.Denoms(sdk.UnwrapSDKContext(ctx), &types.QueryDenomsRequest{})

	if err == nil {
		for _, d := range res.Denoms {
			if mintable && !d.MintRestricted {
				denoms = append(denoms, d.Id)
			}

			if editable && !d.UpdateRestricted {
				denoms = append(denoms, d.Id)
			}
		}
	}

	idx := r.Intn(len(denoms))

	rndDenomID := denoms[idx]
	nfts, err := k.GetNFTs(ctx, rndDenomID)
	if err != nil || len(nfts) == 0 {
		return nil, "", ""
	}

	// get random collection from owner's balance
	token := nfts[r.Intn(len(nfts))]
	return token.GetOwner(), rndDenomID, token.GetID()
}

func genDenomID(r *rand.Rand) string {
	len := simtypes.RandIntBetween(r, idMinLen, idMaxLen)
	var denomID string
	for {
		denomID = strings.ToLower(simtypes.RandStringOfLength(r, len))
		if err := types.ValidateDenomID(denomID); err != nil {
			continue
		}

		if err := types.ValidateKeywords(denomID); err != nil {
			continue
		}
		break
	}
	return denomID
}

func genNFTID(r *rand.Rand, min, max int) string {
	n := simtypes.RandIntBetween(r, min, max)
	id := simtypes.RandStringOfLength(r, n)
	return strings.ToLower(id)
}

func randDenom(ctx sdk.Context, k keeper.Keeper, r *rand.Rand, mintable, editable bool) string {
	res, err := k.Denoms(sdk.UnwrapSDKContext(ctx), &types.QueryDenomsRequest{})
	var denoms = []string{kitties, doggos}
	if err != nil {
		i := r.Intn(len(denoms))
		return denoms[i]
	}

	for _, d := range res.Denoms {
		if mintable && !d.MintRestricted {
			denoms = append(denoms, d.Id)
		}

		if editable && !d.UpdateRestricted {
			denoms = append(denoms, d.Id)
		}
	}
	idx := r.Intn(len(denoms))
	return denoms[idx]
}

func randData(r *rand.Rand) string {
	idx := r.Intn(len(data))
	return data[idx]
}

func genRandomBool(r *rand.Rand) bool {
	return r.Int()%2 == 0
}
