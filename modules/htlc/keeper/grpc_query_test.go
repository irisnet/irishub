package keeper_test

import (
	gocontext "context"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/suite"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/keeper"
	"github.com/irisnet/irismod/modules/htlc/types"
	"github.com/irisnet/irismod/simapp"
)

type QueryTestSuite struct {
	suite.Suite

	cdc    codec.JSONCodec
	ctx    sdk.Context
	keeper *keeper.Keeper
	app    *simapp.SimApp

	queryClient types.QueryClient
	addrs       []sdk.AccAddress
	htlcIDs     []tmbytes.HexBytes
	isHTLCID    map[string]bool
}

func (suite *QueryTestSuite) SetupTest() {
	app := simapp.SetupWithGenesisHTLC(NewHTLTGenesis(TestDeputy))

	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: 1, Time: tmtime.Now()})
	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.keeper = &app.HTLCKeeper
	suite.app = app

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.keeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	suite.setTestHTLCs()
}

func (suite *QueryTestSuite) setTestHTLCs() {
	_, addrs := GeneratePrivKeyAddressPairs(11)
	suite.addrs = addrs

	var htlcIDs []tmbytes.HexBytes
	isHTLCID := make(map[string]bool)
	for i := 0; i < 10; i++ {
		timeLock := MinTimeLock
		amount := cs(c("htltbnb", 100))
		timestamp := ts(0)
		randomSecret, _ := GenerateRandomSecret()
		randomHashLock := types.GetHashLock(randomSecret, timestamp)

		id, err := suite.keeper.CreateHTLC(
			suite.ctx,
			TestDeputy,
			suite.addrs[i],
			ReceiverOnOtherChain,
			SenderOnOtherChain,
			amount,
			randomHashLock,
			timestamp,
			timeLock,
			true,
		)
		suite.Nil(err)

		htlcIDs = append(htlcIDs, id)
		isHTLCID[hex.EncodeToString(id)] = true
	}
	suite.htlcIDs = htlcIDs
	suite.isHTLCID = isHTLCID
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}

func (suite *QueryTestSuite) TestQueryAssetSupply() {
	supplyResp, err := suite.queryClient.AssetSupply(gocontext.Background(), &types.QueryAssetSupplyRequest{Denom: "htltbnb"})
	suite.Require().NoError(err)

	expected, found := suite.keeper.GetAssetSupply(suite.ctx, "htltbnb")
	suite.Require().True(found)
	suite.Equal(expected, *supplyResp.AssetSupply)
}

func (suite *QueryTestSuite) TestQueryAssetSupplies() {
	suppliesResp, err := suite.queryClient.AssetSupplies(gocontext.Background(), &types.QueryAssetSuppliesRequest{})
	suite.Require().NoError(err)

	expected := suite.keeper.GetAllAssetSupplies(suite.ctx)
	suite.Equal(expected, suppliesResp.AssetSupplies)
}

func (suite *QueryTestSuite) TestQueryHTLC() {
	htlcResp, err := suite.queryClient.HTLC(gocontext.Background(), &types.QueryHTLCRequest{Id: suite.htlcIDs[0].String()})
	suite.Require().NoError(err)

	expected, found := suite.keeper.GetHTLC(suite.ctx, suite.htlcIDs[0])
	suite.Require().True(found)
	suite.Equal(expected, *htlcResp.Htlc)
}

func (suite *QueryTestSuite) TestQueryParams() {
	paramsResp, err := suite.queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.Require().NoError(err)

	expected := suite.keeper.GetParams(suite.ctx)
	suite.Equal(expected, paramsResp.Params)
}
