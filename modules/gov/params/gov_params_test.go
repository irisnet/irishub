package govparams

import (
	"fmt"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
	"time"
)

func defaultContext(key sdk.StoreKey, tkeyParams *sdk.TransientStoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	return ctx
}

func TestInitGenesisParameter(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	Denom := "iris"
	IrisCt := types.NewDefaultCoinType(Denom)
	minDeposit, err := IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", 10, Denom))
	require.NoError(t, err)

	p1 := DepositProcedure{
		MinDeposit:       sdk.Coins{minDeposit},
		MaxDepositPeriod: time.Duration(172800) * time.Second}

	minDeposit, err = IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", 20, Denom))
	require.NoError(t, err)

	p2 := DepositProcedure{
		MinDeposit:       sdk.Coins{minDeposit},
		MaxDepositPeriod: time.Duration(172800) * time.Second}

	subspace := paramKeeper.Subspace("Gov").WithTypeTable(
		params.NewTypeTable(
			DepositProcedureParameter.GetStoreKey(), DepositProcedure{},
			VotingProcedureParameter.GetStoreKey(), VotingProcedure{},
			TallyingProcedureParameter.GetStoreKey(), TallyingProcedure{},
		))
	params.SetParamReadWriter(subspace, &DepositProcedureParameter, &DepositProcedureParameter)
	params.InitGenesisParameter(&DepositProcedureParameter, ctx, nil)

	require.Equal(t, p1, DepositProcedureParameter.Value)
	require.Equal(t, DepositProcedureParameter.ToJson(""), "{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":172800000000000}")

	params.InitGenesisParameter(&DepositProcedureParameter, ctx, p2)
	require.Equal(t, p1, DepositProcedureParameter.Value)
}

func TestRegisterParamMapping(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	Denom := "iris"
	IrisCt := types.NewDefaultCoinType(Denom)
	minDeposit, err := IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", 10, Denom))
	require.NoError(t, err)

	p1 := DepositProcedure{
		MinDeposit:       sdk.Coins{minDeposit},
		MaxDepositPeriod: time.Duration(172800) * time.Second}

	minDeposit, err = IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", 30, Denom))
	require.NoError(t, err)

	p2 := DepositProcedure{
		MinDeposit:       sdk.Coins{minDeposit},
		MaxDepositPeriod: time.Duration(172800) * time.Second}

	subspace := paramKeeper.Subspace("Gov").WithTypeTable(
		params.NewTypeTable(
			DepositProcedureParameter.GetStoreKey(), DepositProcedure{},
			VotingProcedureParameter.GetStoreKey(), VotingProcedure{},
			TallyingProcedureParameter.GetStoreKey(), TallyingProcedure{},
		))
	params.SetParamReadWriter(subspace, &DepositProcedureParameter, &DepositProcedureParameter)
	params.RegisterGovParamMapping(&DepositProcedureParameter)
	params.InitGenesisParameter(&DepositProcedureParameter, ctx, nil)

	require.Equal(t, params.ParamMapping["Gov/"+string(DepositProcedureParameter.GetStoreKey())].ToJson(""), "{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Equal(t, p1, DepositProcedureParameter.Value)

	params.ParamMapping["Gov/"+string(DepositProcedureParameter.GetStoreKey())].Update(ctx, "{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"30000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, DepositProcedureParameter.Value)
}

func TestDepositProcedureParam(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	p1deposit, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 10, "iris"))
	p2Deposit, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 200, "iris"))
	p1 := DepositProcedure{
		MinDeposit:       sdk.Coins{p1deposit},
		MaxDepositPeriod: time.Duration(172800) * time.Second}

	p2 := DepositProcedure{
		MinDeposit:       sdk.Coins{p2Deposit},
		MaxDepositPeriod: time.Duration(172800) * time.Second}

	subspace := paramKeeper.Subspace("Gov").WithTypeTable(
		params.NewTypeTable(
			DepositProcedureParameter.GetStoreKey(), DepositProcedure{},
			VotingProcedureParameter.GetStoreKey(), VotingProcedure{},
			TallyingProcedureParameter.GetStoreKey(), TallyingProcedure{},
		))

	DepositProcedureParameter.SetReadWriter(subspace)
	find := DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	DepositProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	require.Equal(t, DepositProcedureParameter.ToJson(""), "{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":172800000000000}")

	DepositProcedureParameter.Update(ctx, "{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"200000000000000000000\"}],\"max_deposit_period\":172800000000000}")

	require.NotEqual(t, p1, DepositProcedureParameter.Value)
	require.Equal(t, p2, DepositProcedureParameter.Value)

	result := DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"atom\",\"amount\":\"200000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"20000000000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"200000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-att\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"20000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.NoError(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":1}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Error(t, result)

	DepositProcedureParameter.InitGenesis(p2)
	require.Equal(t, p2, DepositProcedureParameter.Value)
	DepositProcedureParameter.InitGenesis(p1)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, DepositProcedureParameter.Value)

}

func TestVotingProcedureParam(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	p1 := VotingProcedure{
		VotingPeriod: time.Duration(172800) * time.Second,
	}

	p2 := VotingProcedure{
		VotingPeriod: time.Duration(192800) * time.Second,
	}

	subspace := paramKeeper.Subspace("Gov").WithTypeTable(
		params.NewTypeTable(
			DepositProcedureParameter.GetStoreKey(), DepositProcedure{},
			VotingProcedureParameter.GetStoreKey(), VotingProcedure{},
			TallyingProcedureParameter.GetStoreKey(), TallyingProcedure{},
		))

	VotingProcedureParameter.SetReadWriter(subspace)
	find := VotingProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	VotingProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, VotingProcedureParameter.Value)

	require.Equal(t, VotingProcedureParameter.ToJson(""), "{\"voting_period\":172800000000000}")

	VotingProcedureParameter.Update(ctx, "{\"voting_period\":192800000000000}")

	require.NotEqual(t, p1, VotingProcedureParameter.Value)
	require.Equal(t, p2, VotingProcedureParameter.Value)

	result := VotingProcedureParameter.Valid("{\"voting_period\":400000}")
	require.Error(t, result)

	VotingProcedureParameter.InitGenesis(p2)
	require.Equal(t, p2, VotingProcedureParameter.Value)
	VotingProcedureParameter.InitGenesis(p1)
	require.Equal(t, p1, VotingProcedureParameter.Value)

	VotingProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, VotingProcedureParameter.Value)

}

func TestTallyingProcedureParam(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	p1 := TallyingProcedure{
		Threshold:     sdk.NewDecWithPrec(5, 1),
		Veto:          sdk.NewDecWithPrec(334, 3),
		Participation: sdk.NewDecWithPrec(667, 3),
	}

	p2 := TallyingProcedure{
		Threshold:     sdk.NewDecWithPrec(5, 1),
		Veto:          sdk.NewDecWithPrec(334, 3),
		Participation: sdk.NewDecWithPrec(2, 2),
	}

	subspace := paramKeeper.Subspace("Gov").WithTypeTable(
		params.NewTypeTable(
			DepositProcedureParameter.GetStoreKey(), DepositProcedure{},
			VotingProcedureParameter.GetStoreKey(), VotingProcedure{},
			TallyingProcedureParameter.GetStoreKey(), TallyingProcedure{},
		))

	TallyingProcedureParameter.SetReadWriter(subspace)
	find := TallyingProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	TallyingProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, TallyingProcedureParameter.Value)
	require.Equal(t, "{\"threshold\":\"0.5000000000\",\"veto\":\"0.3340000000\",\"participation\":\"0.6670000000\"}", TallyingProcedureParameter.ToJson(""))

	TallyingProcedureParameter.Update(ctx, "{\"threshold\":\"0.5\",\"veto\":\"0.3340000000\",\"participation\":\"0.0200000000\"}")

	require.NotEqual(t, p1, TallyingProcedureParameter.Value)
	require.Equal(t, p2, TallyingProcedureParameter.Value)

	result := TallyingProcedureParameter.Valid("{\"threshold\":\"1/1\",\"veto\":\"1/3\",\"participation\":\"1/100\"}")
	require.Error(t, result)

	result = TallyingProcedureParameter.Valid("{\"threshold\":\"abcd\",\"veto\":\"1/3\",\"participation\":\"1/100\"}")
	require.Error(t, result)

	TallyingProcedureParameter.InitGenesis(p2)
	require.Equal(t, p2, TallyingProcedureParameter.Value)
	TallyingProcedureParameter.InitGenesis(p1)
	require.Equal(t, p1, TallyingProcedureParameter.Value)

	TallyingProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, TallyingProcedureParameter.Value)

}
