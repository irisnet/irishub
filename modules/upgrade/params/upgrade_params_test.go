package upgradeparams

import (
	"testing"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/codec"
	"github.com/stretchr/testify/require"
	"github.com/irisnet/irishub/store"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/libs/log"
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

func TestCurrentUpgradeProposalIdParameter(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	subspace := paramKeeper.Subspace("Sig").WithTypeTable(params.NewTypeTable(
		CurrentUpgradeProposalIdParameter.GetStoreKey(), uint64((0)),
		ProposalAcceptHeightParameter.GetStoreKey(), int64(0),
		SwitchPeriodParameter.GetStoreKey(), int64(0),
	))

	CurrentUpgradeProposalIdParameter.SetReadWriter(subspace)
	find := CurrentUpgradeProposalIdParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	CurrentUpgradeProposalIdParameter.InitGenesis(nil)
	require.Equal(t, uint64(0), CurrentUpgradeProposalIdParameter.Value)

	CurrentUpgradeProposalIdParameter.LoadValue(ctx)
	require.Equal(t, uint64(0), CurrentUpgradeProposalIdParameter.Value)

	CurrentUpgradeProposalIdParameter.Value = 3
	CurrentUpgradeProposalIdParameter.SaveValue(ctx)

	CurrentUpgradeProposalIdParameter.LoadValue(ctx)
	require.Equal(t, uint64(3), CurrentUpgradeProposalIdParameter.Value)
}

func TestProposalAcceptHeightParameter(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	subspace := paramKeeper.Subspace("Sig").WithTypeTable(params.NewTypeTable(
		CurrentUpgradeProposalIdParameter.GetStoreKey(), uint64((0)),
		ProposalAcceptHeightParameter.GetStoreKey(), int64(0),
		SwitchPeriodParameter.GetStoreKey(), int64(0),
	))

	ProposalAcceptHeightParameter.SetReadWriter(subspace)
	find := ProposalAcceptHeightParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	ProposalAcceptHeightParameter.InitGenesis(nil)
	require.Equal(t, int64(-1), ProposalAcceptHeightParameter.Value)

	ProposalAcceptHeightParameter.LoadValue(ctx)
	require.Equal(t, int64(-1), ProposalAcceptHeightParameter.Value)

	ProposalAcceptHeightParameter.Value = 3
	ProposalAcceptHeightParameter.SaveValue(ctx)

	ProposalAcceptHeightParameter.LoadValue(ctx)
	require.Equal(t, int64(3), ProposalAcceptHeightParameter.Value)
}

func TestSwitchPeriodParameter(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	subspace := paramKeeper.Subspace("Sig").WithTypeTable(params.NewTypeTable(
		CurrentUpgradeProposalIdParameter.GetStoreKey(), uint64((0)),
		ProposalAcceptHeightParameter.GetStoreKey(), int64(0),
		SwitchPeriodParameter.GetStoreKey(), int64(0),
	))

	SwitchPeriodParameter.SetReadWriter(subspace)
	find := SwitchPeriodParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	SwitchPeriodParameter.InitGenesis(int64(12345))
	require.Equal(t, int64(12345), SwitchPeriodParameter.Value)

	SwitchPeriodParameter.LoadValue(ctx)
	require.Equal(t, int64(12345), SwitchPeriodParameter.Value)

	SwitchPeriodParameter.Value = 30
	SwitchPeriodParameter.SaveValue(ctx)

	SwitchPeriodParameter.LoadValue(ctx)
	require.Equal(t, int64(30), SwitchPeriodParameter.Value)
}

func TestUpgradeParameterSetAndGet(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	subspace := paramKeeper.Subspace("Sig").WithTypeTable(params.NewTypeTable(
		CurrentUpgradeProposalIdParameter.GetStoreKey(), uint64((0)),
		ProposalAcceptHeightParameter.GetStoreKey(), int64(0),
		SwitchPeriodParameter.GetStoreKey(), int64(0),
	))

	CurrentUpgradeProposalIdParameter.SetReadWriter(subspace)
	find := CurrentUpgradeProposalIdParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	ProposalAcceptHeightParameter.SetReadWriter(subspace)
	find = ProposalAcceptHeightParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	SwitchPeriodParameter.SetReadWriter(subspace)
	find = SwitchPeriodParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	SetCurrentUpgradeProposalId(ctx,5)
	require.Equal(t,uint64(5),GetCurrentUpgradeProposalId(ctx))
	SetProposalAcceptHeight(ctx,100)
	require.Equal(t, int64(100),GetProposalAcceptHeight(ctx) )

	SetSwitchPeriod(ctx,500000)
	require.Equal(t,int64(500000),GetSwitchPeriod(ctx))
	SetSwitchPeriod(ctx,30000)
	require.Equal(t,int64(30000),GetSwitchPeriod(ctx))
}
