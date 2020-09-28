package keeper_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/guardian/keeper"
	"github.com/irisnet/irishub/modules/guardian/types"
	"github.com/irisnet/irishub/simapp"
)

var (
	pks = []crypto.PubKey{
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB52"),
	}
	addrs = []sdk.AccAddress{
		sdk.AccAddress(pks[0].Address()),
		sdk.AccAddress(pks[1].Address()),
		sdk.AccAddress(pks[2].Address()),
	}
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    *codec.LegacyAmino
	ctx    sdk.Context
	keeper keeper.Keeper
	app    *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.app = app
	suite.cdc = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = app.GuardianKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestAddProfiler() {
	profiler := types.NewGuardian("test", types.Genesis, addrs[0], addrs[1])

	suite.keeper.AddProfiler(suite.ctx, profiler)
	addedProfiler, found := suite.keeper.GetProfiler(suite.ctx, addrs[0])
	suite.True(found)
	suite.True(profiler.Equal(addedProfiler))

	var profilers []types.Guardian
	suite.keeper.IterateProfilers(
		suite.ctx,
		func(profiler types.Guardian) bool {
			profilers = append(profilers, profiler)
			return false
		},
	)

	suite.Equal(1, len(profilers))
	suite.Contains(profilers, profiler)
}

func (suite *KeeperTestSuite) TestDeleteProfiler() {
	profiler := types.NewGuardian("test", types.Genesis, addrs[0], addrs[1])

	suite.keeper.AddProfiler(suite.ctx, profiler)
	addedProfiler, found := suite.keeper.GetProfiler(suite.ctx, addrs[0])
	suite.True(found)
	suite.True(profiler.Equal(addedProfiler))

	suite.keeper.DeleteProfiler(suite.ctx, profiler.Address)

	_, found = suite.keeper.GetProfiler(suite.ctx, addrs[0])
	suite.False(found)
}

func (suite *KeeperTestSuite) TestAddTrustee() {
	trustee := types.NewGuardian("test", types.Genesis, addrs[0], addrs[1])
	suite.keeper.AddTrustee(suite.ctx, trustee)
	addedTrustee, found := suite.keeper.GetTrustee(suite.ctx, addrs[0])
	suite.True(found)
	suite.True(trustee.Equal(addedTrustee))

	var trustees []types.Guardian
	suite.keeper.IterateTrustees(
		suite.ctx,
		func(trustee types.Guardian) bool {
			trustees = append(trustees, trustee)
			return false
		},
	)
	suite.Equal(1, len(trustees))
	suite.Contains(trustees, trustee)
}

func (suite *KeeperTestSuite) TestDeleteTrustee() {
	trustee := types.NewGuardian("test", types.Genesis, addrs[0], addrs[1])
	suite.keeper.AddTrustee(suite.ctx, trustee)
	addedTrustee, found := suite.keeper.GetTrustee(suite.ctx, addrs[0])
	suite.True(found)
	suite.True(trustee.Equal(addedTrustee))

	suite.keeper.DeleteTrustee(suite.ctx, trustee.Address)
	_, found = suite.keeper.GetProfiler(suite.ctx, trustee.Address)
	suite.False(found)

}

func (suite *KeeperTestSuite) TestQueryProfilers() {
	profiler := types.NewGuardian("test", types.Genesis, addrs[0], addrs[1])
	suite.keeper.AddProfiler(suite.ctx, profiler)

	var profilers []types.Guardian
	querier := keeper.NewQuerier(suite.keeper, suite.cdc)
	res, sdkErr := querier(suite.ctx, []string{types.QueryProfilers}, abci.RequestQuery{})
	suite.NoError(sdkErr)

	err := suite.cdc.UnmarshalJSON(res, &profilers)
	suite.NoError(err)
	suite.Len(profilers, 1)
	suite.Contains(profilers, profiler)
}

func (suite *KeeperTestSuite) TestQueryTrustees() {
	trustee := types.NewGuardian("test", types.Genesis, addrs[0], addrs[1])
	suite.keeper.AddTrustee(suite.ctx, trustee)

	var trustees []types.Guardian
	querier := keeper.NewQuerier(suite.keeper, suite.cdc)
	res, sdkErr := querier(suite.ctx, []string{types.QueryTrustees}, abci.RequestQuery{})
	suite.NoError(sdkErr)

	err := suite.cdc.UnmarshalJSON(res, &trustees)
	suite.NoError(err)
	suite.Len(trustees, 1)
	suite.Contains(trustees, trustee)
}

func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd = make(ed25519.PubKey, ed25519.PubKeySize)
	copy(pkEd[:], pkBytes[:])
	return pkEd
}
