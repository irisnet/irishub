package keeper_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/guardian"
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

	cdc    *codec.Codec
	ctx    sdk.Context
	keeper guardian.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(false, abci.Header{})
	suite.keeper = app.GuardianKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestAddProfiler() {
	profiler := guardian.NewGuardian("test", guardian.Genesis, addrs[0], addrs[1])

	suite.keeper.AddProfiler(suite.ctx, profiler)
	addedProfilerI, found := suite.keeper.GetProfiler(suite.ctx, addrs[0])
	suite.True(found)
	addedProfiler, ok := addedProfilerI.(guardian.Guardian)
	suite.True(ok)
	suite.True(profiler.Equal(addedProfiler))

	profilersIterator := suite.keeper.ProfilersIterator(suite.ctx)
	defer profilersIterator.Close()
	var profilers []guardian.Guardian
	for ; profilersIterator.Valid(); profilersIterator.Next() {
		var profiler guardian.Guardian
		suite.cdc.MustUnmarshalBinaryLengthPrefixed(profilersIterator.Value(), &profiler)
		profilers = append(profilers, profiler)
	}
	suite.Equal(2, len(profilers))
	suite.Contains(profilers, profiler)
}

func (suite *KeeperTestSuite) TestDeleteProfiler() {
	profiler := guardian.NewGuardian("test", guardian.Genesis, addrs[0], addrs[1])

	suite.keeper.AddProfiler(suite.ctx, profiler)
	addedProfilerI, found := suite.keeper.GetProfiler(suite.ctx, addrs[0])
	suite.True(found)
	addedProfiler, ok := addedProfilerI.(guardian.Guardian)
	suite.True(ok)
	suite.True(profiler.Equal(addedProfiler))

	suite.keeper.DeleteProfiler(suite.ctx, profiler.Address)

	_, found = suite.keeper.GetProfiler(suite.ctx, addrs[0])
	suite.False(found)
}

func (suite *KeeperTestSuite) TestAddTrustee() {
	trustee := guardian.NewGuardian("test", guardian.Genesis, addrs[0], addrs[1])
	suite.keeper.AddTrustee(suite.ctx, trustee)
	addedTrusteeI, found := suite.keeper.GetTrustee(suite.ctx, addrs[0])
	suite.True(found)
	addedTrustee, ok := addedTrusteeI.(guardian.Guardian)
	suite.True(ok)
	suite.True(trustee.Equal(addedTrustee))

	trusteesIterator := suite.keeper.TrusteesIterator(suite.ctx)
	defer trusteesIterator.Close()
	var trustees []guardian.Guardian
	for ; trusteesIterator.Valid(); trusteesIterator.Next() {
		var trustee guardian.Guardian
		suite.cdc.MustUnmarshalBinaryLengthPrefixed(trusteesIterator.Value(), &trustee)
		trustees = append(trustees, trustee)
	}
	suite.Equal(2, len(trustees))
	suite.Contains(trustees, trustee)
}

func (suite *KeeperTestSuite) TestDeleteTrustee() {
	trustee := guardian.NewGuardian("test", guardian.Genesis, addrs[0], addrs[1])
	suite.keeper.AddTrustee(suite.ctx, trustee)
	addedTrusteeI, found := suite.keeper.GetTrustee(suite.ctx, addrs[0])
	suite.True(found)
	addedTrustee, ok := addedTrusteeI.(guardian.Guardian)
	suite.True(ok)
	suite.True(trustee.Equal(addedTrustee))

	suite.keeper.DeleteTrustee(suite.ctx, trustee.Address)
	_, found = suite.keeper.GetProfiler(suite.ctx, trustee.Address)
	suite.False(found)

}

func (suite *KeeperTestSuite) TestQueryProfilers() {
	profiler := guardian.NewGuardian("test", guardian.Genesis, addrs[0], addrs[1])
	suite.keeper.AddProfiler(suite.ctx, profiler)

	var profilers []guardian.Guardian
	querier := guardian.NewQuerier(suite.keeper)
	res, sdkErr := querier(suite.ctx, []string{guardian.QueryProfilers}, abci.RequestQuery{})
	suite.NoError(sdkErr)

	err := suite.cdc.UnmarshalJSON(res, &profilers)
	suite.NoError(err)
	suite.Len(profilers, 2)
	suite.Contains(profilers, profiler)
}

func (suite *KeeperTestSuite) TestQueryTrustees() {
	trustee := guardian.NewGuardian("test", guardian.Genesis, addrs[0], addrs[1])
	suite.keeper.AddTrustee(suite.ctx, trustee)

	var trustees []guardian.Guardian
	querier := guardian.NewQuerier(suite.keeper)
	res, sdkErr := querier(suite.ctx, []string{guardian.QueryTrustees}, abci.RequestQuery{})
	suite.NoError(sdkErr)

	err := suite.cdc.UnmarshalJSON(res, &trustees)
	suite.NoError(err)
	suite.Len(trustees, 2)
	suite.Contains(trustees, trustee)
}

func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd
}
