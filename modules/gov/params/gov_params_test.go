package govparams

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
	"fmt"
	"github.com/irisnet/irishub/types"
)

func defaultContext(key sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	return ctx
}

func TestDepositProcedureParam(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	ctx := defaultContext(skey)
	paramKeeper := params.NewKeeper(wire.NewCodec(), skey)

	p1deposit, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 10, "iris"))
	p2Deposit, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 200, "iris"))
	p1 := DepositProcedure{
		MinDeposit:       sdk.Coins{p1deposit},
		MaxDepositPeriod: 1440}

	p2 := DepositProcedure{
		MinDeposit:       sdk.Coins{p2Deposit},
		MaxDepositPeriod: 1440}

	DepositProcedureParameter.SetReadWriter(paramKeeper.Setter())
	find := DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	DepositProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	require.Equal(t, DepositProcedureParameter.ToJson(), "{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":1440}")

	DepositProcedureParameter.Update(ctx, "{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"200000000000000000000\"}],\"max_deposit_period\":1440}")

	require.NotEqual(t, p1, DepositProcedureParameter.Value)
	require.Equal(t, p2, DepositProcedureParameter.Value)

	result := DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"atom\",\"amount\":\"200000000000000000000\"}],\"max_deposit_period\":1440}")
	require.Error(t, result)

	DepositProcedureParameter.InitGenesis(p2)
	require.Equal(t, p2, DepositProcedureParameter.Value)
	DepositProcedureParameter.InitGenesis(p1)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, DepositProcedureParameter.Value)

}

func TestDepositProcedureParamValid(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	ctx := defaultContext(skey)
	paramKeeper := params.NewKeeper(wire.NewCodec(), skey)

	p1deposit, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 10, "iris"))
	p1 := DepositProcedure{
		MinDeposit:       sdk.Coins{p1deposit},
		MaxDepositPeriod: 1440}

	DepositProcedureParameter.SetReadWriter(paramKeeper.Setter())
	find := DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	DepositProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	result := DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":1440}")
	require.NoError(t, result)
	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000000\"}],\"max_deposit_period\":1440}")
	require.Error(t, result)
	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"200000000000000000\"}],\"max_deposit_period\":1440}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-att\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":1440}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":1440}")
	require.NoError(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":1}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":1440000}")
	require.Error(t, result)

}

func TestVotingProcedureParam(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	ctx := defaultContext(skey)
	paramKeeper := params.NewKeeper(wire.NewCodec(), skey)


	p1 := VotingProcedure{
		VotingPeriod:1000,
	}

	p2 := VotingProcedure{
		VotingPeriod: 2000,
	}

	VotingProcedureParameter.SetReadWriter(paramKeeper.Setter())
	find := VotingProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	VotingProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, VotingProcedureParameter.Value)

	require.Equal(t, VotingProcedureParameter.ToJson(), "{\"voting_period\":1000}")

	VotingProcedureParameter.Update(ctx, "{\"voting_period\":2000}")

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
	ctx := defaultContext(skey)
	paramKeeper := params.NewKeeper(wire.NewCodec(), skey)


	p1 := TallyingProcedure{
		Threshold:         sdk.NewRat(1, 2),
		Veto:              sdk.NewRat(1, 3),
		GovernancePenalty: sdk.NewRat(1, 100),
	}

	p2 := TallyingProcedure{
        Threshold:         sdk.NewRat(1, 2),
        Veto:              sdk.NewRat(1, 3),
        GovernancePenalty: sdk.NewRat(1, 50),
    }

	TallyingProcedureParameter.SetReadWriter(paramKeeper.Setter())
	find := TallyingProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	TallyingProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, TallyingProcedureParameter.Value)
	require.Equal(t, TallyingProcedureParameter.ToJson(), "{\"threshold\":\"1/2\",\"veto\":\"1/3\",\"governance_penalty\":\"1/100\"}")

	TallyingProcedureParameter.Update(ctx, "{\"threshold\":\"0.5\",\"veto\":\"1/3\",\"governance_penalty\":\"1/50\"}")

	require.NotEqual(t, p1, TallyingProcedureParameter.Value)
	require.Equal(t, p2, TallyingProcedureParameter.Value)

	result := TallyingProcedureParameter.Valid("{\"threshold\":\"2/1\",\"veto\":\"1/3\",\"governance_penalty\":\"1/100\"}")
	require.Error(t, result)

	result = TallyingProcedureParameter.Valid("{\"threshold\":\"abcd\",\"veto\":\"1/3\",\"governance_penalty\":\"1/100\"}")
	require.Error(t, result)

	TallyingProcedureParameter.InitGenesis(p2)
	require.Equal(t, p2, TallyingProcedureParameter.Value)
	TallyingProcedureParameter.InitGenesis(p1)
	require.Equal(t, p1, TallyingProcedureParameter.Value)

	TallyingProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, TallyingProcedureParameter.Value)

}