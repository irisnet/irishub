package keeper_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/keeper"
	"github.com/irisnet/irismod/modules/htlc/types"
	"github.com/irisnet/irismod/simapp"
)

type ParamsTestSuite struct {
	suite.Suite

	cdc    codec.BinaryCodec
	ctx    sdk.Context
	keeper *keeper.Keeper
	app    *simapp.SimApp
}

func (suite *ParamsTestSuite) SetupTest() {
	app := simapp.SetupWithGenesisHTLC(NewHTLTGenesis(TestDeputy))
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: 1, Time: tmtime.Now()})

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.keeper = &app.HTLCKeeper
	suite.app = app

}

func TestParamsTestSuite(t *testing.T) {
	suite.Run(t, new(AssetTestSuite))
}

func (suite *ParamsTestSuite) TestGetSetAsset() {
	asset, err := suite.keeper.GetAsset(suite.ctx, "htltbnb")
	suite.Require().NoError(err)
	suite.NotPanics(func() { suite.keeper.SetAsset(suite.ctx, asset) })
	_, err = suite.keeper.GetAsset(suite.ctx, "htltdne")
	suite.Require().Error(err)

	_, err = suite.keeper.GetAsset(suite.ctx, "htltinc")
	suite.Require().NoError(err)
}

func (suite *ParamsTestSuite) TestGetAssets() {
	assets, found := suite.keeper.GetAssets(suite.ctx)
	suite.Require().True(found)
	suite.Require().Equal(2, len(assets))
}

func (suite *ParamsTestSuite) TestGetSetDeputyAddress() {
	asset, err := suite.keeper.GetAsset(suite.ctx, "htltbnb")
	suite.Require().NoError(err)
	asset.DeputyAddress = TestDeputy.String()
	suite.NotPanics(func() { suite.keeper.SetAsset(suite.ctx, asset) })

	asset, err = suite.keeper.GetAsset(suite.ctx, "htltbnb")
	suite.Require().NoError(err)
	suite.Equal(TestDeputy.String(), asset.DeputyAddress)
	addr, err := suite.keeper.GetDeputyAddress(suite.ctx, "htltbnb")
	suite.Require().NoError(err)
	suite.Equal(TestDeputy.String(), addr)

}

func (suite *ParamsTestSuite) TestGetDeputyFixedFee() {
	asset, err := suite.keeper.GetAsset(suite.ctx, "htltbnb")
	suite.Require().NoError(err)
	bnbDeputyFixedFee := asset.FixedFee

	res, err := suite.keeper.GetFixedFee(suite.ctx, asset.Denom)
	suite.Require().NoError(err)
	suite.Equal(bnbDeputyFixedFee, res)
}

func (suite *ParamsTestSuite) TestGetMinMaxSwapAmount() {
	asset, err := suite.keeper.GetAsset(suite.ctx, "htltbnb")
	suite.Require().NoError(err)
	minAmount := asset.MinSwapAmount

	res, err := suite.keeper.GetMinSwapAmount(suite.ctx, asset.Denom)
	suite.Require().NoError(err)
	suite.Equal(minAmount, res)

	maxAmount := asset.MaxSwapAmount
	res, err = suite.keeper.GetMaxSwapAmount(suite.ctx, asset.Denom)
	suite.Require().NoError(err)
	suite.Equal(maxAmount, res)
}

func (suite *ParamsTestSuite) TestGetMinMaxBlockLock() {
	asset, err := suite.keeper.GetAsset(suite.ctx, "htltbnb")
	suite.Require().NoError(err)
	minLock := asset.MinBlockLock

	res, err := suite.keeper.GetMinBlockLock(suite.ctx, asset.Denom)
	suite.Require().NoError(err)
	suite.Equal(minLock, res)

	maxLock := asset.MaxBlockLock
	res, err = suite.keeper.GetMaxBlockLock(suite.ctx, asset.Denom)
	suite.Require().NoError(err)
	suite.Equal(maxLock, res)
}

func (suite *ParamsTestSuite) TestValidateLiveAsset() {
	type args struct {
		coin sdk.Coin
	}
	testCases := []struct {
		name          string
		args          args
		expectedError error
		expectPass    bool
	}{{
		"normal",
		args{coin: c("htltbnb", 1)},
		nil,
		true,
	}, {
		"asset not supported",
		args{coin: c("bad", 1)},
		types.ErrAssetNotSupported,
		false,
	}, {
		"asset not active",
		args{coin: c("htltinc", 1)},
		types.ErrAssetNotActive,
		false,
	}}

	for _, tc := range testCases {
		suite.SetupTest()
		suite.Run(tc.name, func() {
			err := suite.keeper.ValidateLiveAsset(suite.ctx, tc.args.coin)
			if tc.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().True(errors.Is(err, tc.expectedError))
			}
		})
	}
}
