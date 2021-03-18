package keeper_test

// import (
// 	"errors"

// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/irisnet/irismod/modules/htlc/types"
// )

// func (suite *KeeperTestSuite) TestGetSetAsset() {
// 	asset, err := suite.keeper.GetAsset(suite.ctx, "htltbnb")
// 	suite.Require().NoError(err)
// 	suite.NotPanics(func() { suite.keeper.SetAsset(suite.ctx, asset) })
// 	_, err = suite.keeper.GetAsset(suite.ctx, "htltdne")
// 	suite.Require().Error(err)

// 	_, err = suite.keeper.GetAsset(suite.ctx, "htltinc")
// 	suite.Require().NoError(err)
// }

// func (suite *KeeperTestSuite) TestGetAssets() {
// 	assets, found := suite.keeper.GetAssets(suite.ctx)
// 	suite.Require().True(found)
// 	suite.Require().Equal(2, len(assets))
// }

// func (suite *KeeperTestSuite) TestGetSetDeputyAddress() {
// 	asset, err := suite.keeper.GetAsset(suite.ctx, "htltbnb")
// 	suite.Require().NoError(err)
// 	asset.DeputyAddress = deputyAddressStr
// 	suite.NotPanics(func() { suite.keeper.SetAsset(suite.ctx, asset) })

// 	asset, err = suite.keeper.GetAsset(suite.ctx, "htltbnb")
// 	suite.Require().NoError(err)
// 	suite.Equal(deputyAddressStr, asset.DeputyAddress)
// 	addr, err := suite.keeper.GetDeputyAddress(suite.ctx, "htltbnb")
// 	suite.Require().NoError(err)
// 	suite.Equal(deputyAddressStr, addr)

// }

// func (suite *KeeperTestSuite) TestGetDeputyFixedFee() {
// 	asset, err := suite.keeper.GetAsset(suite.ctx, "htltbnb")
// 	suite.Require().NoError(err)
// 	bnbDeputyFixedFee := asset.FixedFee

// 	res, err := suite.keeper.GetFixedFee(suite.ctx, asset.Denom)
// 	suite.Require().NoError(err)
// 	suite.Equal(bnbDeputyFixedFee, res)
// }

// func (suite *KeeperTestSuite) TestGetMinMaxSwapAmount() {
// 	asset, err := suite.keeper.GetAsset(suite.ctx, "htltbnb")
// 	suite.Require().NoError(err)
// 	minAmount := asset.MinSwapAmount

// 	res, err := suite.keeper.GetMinSwapAmount(suite.ctx, asset.Denom)
// 	suite.Require().NoError(err)
// 	suite.Equal(minAmount, res)

// 	maxAmount := asset.MaxSwapAmount
// 	res, err = suite.keeper.GetMaxSwapAmount(suite.ctx, asset.Denom)
// 	suite.Require().NoError(err)
// 	suite.Equal(maxAmount, res)
// }

// func (suite *KeeperTestSuite) TestGetMinMaxBlockLock() {
// 	asset, err := suite.keeper.GetAsset(suite.ctx, "htltbnb")
// 	suite.Require().NoError(err)
// 	minLock := asset.MinBlockLock

// 	res, err := suite.keeper.GetMinBlockLock(suite.ctx, asset.Denom)
// 	suite.Require().NoError(err)
// 	suite.Equal(minLock, res)

// 	maxLock := asset.MaxBlockLock
// 	res, err = suite.keeper.GetMaxBlockLock(suite.ctx, asset.Denom)
// 	suite.Require().NoError(err)
// 	suite.Equal(maxLock, res)
// }

// func (suite *KeeperTestSuite) TestValidateLiveAsset() {
// 	type args struct {
// 		coin sdk.Coin
// 	}
// 	testCases := []struct {
// 		name          string
// 		args          args
// 		expectedError error
// 		expectPass    bool
// 	}{
// 		{
// 			"normal",
// 			args{
// 				coin: c("htltbnb", 1),
// 			},
// 			nil,
// 			true,
// 		},
// 		{
// 			"asset not supported",
// 			args{
// 				coin: c("bad", 1),
// 			},
// 			types.ErrAssetNotSupported,
// 			false,
// 		},
// 		{
// 			"asset not active",
// 			args{
// 				coin: c("htltinc", 1),
// 			},
// 			types.ErrAssetNotActive,
// 			false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		suite.SetupTest()
// 		suite.Run(tc.name, func() {
// 			err := suite.keeper.ValidateLiveAsset(suite.ctx, tc.args.coin)

// 			if tc.expectPass {
// 				suite.Require().NoError(err)
// 			} else {
// 				suite.Require().Error(err)
// 				suite.Require().True(errors.Is(err, tc.expectedError))
// 			}
// 		})
// 	}
// }
