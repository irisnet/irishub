package keeper_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/keeper"
	"github.com/irisnet/irismod/modules/htlc/types"
	"github.com/irisnet/irismod/simapp"
)

type AssetTestSuite struct {
	suite.Suite

	cdc    codec.JSONMarshaler
	ctx    sdk.Context
	keeper *keeper.Keeper
	app    *simapp.SimApp
}

func (suite *AssetTestSuite) SetupTest() {
	app := simapp.SetupWithGenesisHTLC(NewHTLTGenesis(TestDeputy))
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: 1, Time: tmtime.Now()})

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.keeper = &app.HTLCKeeper
	suite.app = app

	suite.setTestParams()
}

func (suite *AssetTestSuite) setTestParams() {
	params := suite.keeper.GetParams(suite.ctx)
	params.AssetParams[0].SupplyLimit.Limit = sdk.NewInt(50)
	params.AssetParams[1].SupplyLimit.Limit = sdk.NewInt(100)
	params.AssetParams[1].SupplyLimit.TimeBasedLimit = sdk.NewInt(15)
	suite.keeper.SetParams(suite.ctx, params)

	bnbSupply := types.NewAssetSupply(
		c("htltbnb", 5),
		c("htltbnb", 5),
		c("htltbnb", 40),
		c("htltbnb", 0),
		time.Duration(0),
	)
	incSupply := types.NewAssetSupply(
		c("htltinc", 10),
		c("htltinc", 5),
		c("htltinc", 5),
		c("htltinc", 0),
		time.Duration(0),
	)
	suite.keeper.SetAssetSupply(suite.ctx, bnbSupply, bnbSupply.IncomingSupply.Denom)
	suite.keeper.SetAssetSupply(suite.ctx, incSupply, incSupply.IncomingSupply.Denom)
	suite.keeper.SetPreviousBlockTime(suite.ctx, suite.ctx.BlockTime())
}

func TestAssetTestSuite(t *testing.T) {
	suite.Run(t, new(AssetTestSuite))
}

func (suite *AssetTestSuite) TestIncrementCurrentAssetSupply() {
	type args struct {
		coin sdk.Coin
	}
	testCases := []struct {
		name       string
		args       args
		expectPass bool
	}{{
		"normal",
		args{coin: c("htltbnb", 5)},
		true,
	}, {
		"equal limit",
		args{coin: c("htltbnb", 10)},
		true,
	}, {
		"exceeds limit",
		args{coin: c("htltbnb", 11)},
		false,
	}, {
		"unsupported asset",
		args{coin: c("htltxyz", 5)},
		false,
	}, {
		"unsupported asset",
		args{coin: c("xyz", 5)},
		false,
	}}

	for _, tc := range testCases {
		suite.SetupTest()
		suite.Run(tc.name, func() {
			preSupply, found := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)
			err := suite.keeper.IncrementCurrentAssetSupply(suite.ctx, tc.args.coin)
			postSupply, _ := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)

			if tc.expectPass {
				suite.True(found)
				suite.NoError(err)
				suite.Equal(preSupply.CurrentSupply.Add(tc.args.coin), postSupply.CurrentSupply)
			} else {
				suite.Error(err)
				suite.Equal(preSupply.CurrentSupply, postSupply.CurrentSupply)
			}
		})
	}
}

func (suite *AssetTestSuite) TestIncrementTimeLimitedCurrentAssetSupply() {
	type args struct {
		coin           sdk.Coin
		expectedSupply types.AssetSupply
	}
	type errArgs struct {
		expectPass bool
		contains   string
	}
	testCases := []struct {
		name    string
		args    args
		errArgs errArgs
	}{{
		"normal",
		args{
			coin: c("htltinc", 5),
			expectedSupply: types.AssetSupply{
				IncomingSupply:           c("htltinc", 10),
				OutgoingSupply:           c("htltinc", 5),
				CurrentSupply:            c("htltinc", 10),
				TimeLimitedCurrentSupply: c("htltinc", 5),
				TimeElapsed:              time.Duration(0)},
		},
		errArgs{
			expectPass: true,
			contains:   "",
		},
	}, {
		"over limit",
		args{
			coin:           c("htltinc", 16),
			expectedSupply: types.AssetSupply{},
		},
		errArgs{
			expectPass: false,
			contains:   "asset supply over limit for current time period",
		},
	}}
	for _, tc := range testCases {
		suite.SetupTest()
		suite.Run(
			tc.name,
			func() {
				err := suite.keeper.IncrementCurrentAssetSupply(suite.ctx, tc.args.coin)
				if tc.errArgs.expectPass {
					suite.Require().NoError(err)
					supply, _ := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)
					suite.Require().Equal(tc.args.expectedSupply, supply)
				} else {
					suite.Require().Error(err)
					suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
				}
			},
		)
	}
}

func (suite *AssetTestSuite) TestDecrementCurrentAssetSupply() {
	type args struct {
		coin sdk.Coin
	}
	testCases := []struct {
		name       string
		args       args
		expectPass bool
	}{{
		"normal",
		args{coin: c("htltbnb", 30)},
		true,
	}, {
		"equal current",
		args{coin: c("htltbnb", 40)},
		true,
	}, {
		"exceeds current",
		args{coin: c("htltbnb", 41)},
		false,
	}, {
		"unsupported asset",
		args{coin: c("htltxyz", 30)},
		false,
	}}

	for _, tc := range testCases {
		suite.SetupTest()
		suite.Run(
			tc.name,
			func() {
				preSupply, found := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)
				err := suite.keeper.DecrementCurrentAssetSupply(suite.ctx, tc.args.coin)
				postSupply, _ := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)

				if tc.expectPass {
					suite.True(found)
					suite.NoError(err)
					suite.True(preSupply.CurrentSupply.Sub(tc.args.coin).IsEqual(postSupply.CurrentSupply))
				} else {
					suite.Error(err)
					suite.Equal(preSupply.CurrentSupply, postSupply.CurrentSupply)
				}
			},
		)
	}
}

func (suite *AssetTestSuite) TestIncrementIncomingAssetSupply() {
	type args struct {
		coin sdk.Coin
	}
	testCases := []struct {
		name       string
		args       args
		expectPass bool
	}{{
		"normal",
		args{coin: c("htltbnb", 2)},
		true,
	}, {
		"incoming + current = limit",
		args{coin: c("htltbnb", 5)},
		true,
	}, {
		"incoming + current > limit",
		args{coin: c("htltbnb", 6)},
		false,
	}, {
		"unsupported asset",
		args{coin: c("htltxyz", 2)},
		false,
	}}

	for _, tc := range testCases {
		suite.SetupTest()
		suite.Run(
			tc.name,
			func() {
				preSupply, found := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)
				err := suite.keeper.IncrementIncomingAssetSupply(suite.ctx, tc.args.coin)
				postSupply, _ := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)

				if tc.expectPass {
					suite.True(found)
					suite.NoError(err)
					suite.Equal(preSupply.IncomingSupply.Add(tc.args.coin), postSupply.IncomingSupply)
				} else {
					suite.Error(err)
					suite.Equal(preSupply.IncomingSupply, postSupply.IncomingSupply)
				}
			},
		)
	}
}

func (suite *AssetTestSuite) TestIncrementTimeLimitedIncomingAssetSupply() {
	type args struct {
		coin           sdk.Coin
		expectedSupply types.AssetSupply
	}
	type errArgs struct {
		expectPass bool
		contains   string
	}
	testCases := []struct {
		name    string
		args    args
		errArgs errArgs
	}{{
		"normal",
		args{
			coin: c("htltinc", 5),
			expectedSupply: types.AssetSupply{
				IncomingSupply:           c("htltinc", 15),
				OutgoingSupply:           c("htltinc", 5),
				CurrentSupply:            c("htltinc", 5),
				TimeLimitedCurrentSupply: c("htltinc", 0),
				TimeElapsed:              time.Duration(0)},
		},
		errArgs{
			expectPass: true,
			contains:   "",
		},
	}, {
		"over limit",
		args{
			coin:           c("htltinc", 6),
			expectedSupply: types.AssetSupply{},
		},
		errArgs{
			expectPass: false,
			contains:   "asset supply over limit for current time period",
		},
	}}

	for _, tc := range testCases {
		suite.SetupTest()
		suite.Run(
			tc.name,
			func() {
				err := suite.keeper.IncrementIncomingAssetSupply(suite.ctx, tc.args.coin)
				if tc.errArgs.expectPass {
					suite.Require().NoError(err)
					supply, _ := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)
					suite.Require().Equal(tc.args.expectedSupply, supply)
				} else {
					suite.Require().Error(err)
					suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
				}
			},
		)
	}
}

func (suite *AssetTestSuite) TestDecrementIncomingAssetSupply() {
	type args struct {
		coin sdk.Coin
	}
	testCases := []struct {
		name       string
		args       args
		expectPass bool
	}{{
		"normal",
		args{coin: c("htltbnb", 4)},
		true,
	}, {
		"equal incoming",
		args{coin: c("htltbnb", 5)},
		true,
	}, {
		"exceeds incoming",
		args{coin: c("htltbnb", 6)},
		false,
	}, {
		"unsupported asset",
		args{coin: c("htltxyz", 4)},
		false,
	}}

	for _, tc := range testCases {
		suite.SetupTest()
		suite.Run(
			tc.name,
			func() {
				preSupply, found := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)
				err := suite.keeper.DecrementIncomingAssetSupply(suite.ctx, tc.args.coin)
				postSupply, _ := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)

				if tc.expectPass {
					suite.True(found)
					suite.NoError(err)
					suite.True(preSupply.IncomingSupply.Sub(tc.args.coin).IsEqual(postSupply.IncomingSupply))
				} else {
					suite.Error(err)
					suite.Equal(preSupply.IncomingSupply, postSupply.IncomingSupply)
				}
			},
		)
	}
}

func (suite *AssetTestSuite) TestIncrementOutgoingAssetSupply() {
	type args struct {
		coin sdk.Coin
	}
	testCases := []struct {
		name       string
		args       args
		expectPass bool
	}{{
		"normal",
		args{coin: c("htltbnb", 30)},
		true,
	}, {
		"outgoing + amount = current",
		args{coin: c("htltbnb", 35)},
		true,
	}, {
		"outoing + amount > current",
		args{coin: c("htltbnb", 36)},
		false,
	}, {
		"unsupported asset",
		args{coin: c("htltxyz", 30)},
		false,
	}}

	for _, tc := range testCases {
		suite.SetupTest()
		suite.Run(
			tc.name,
			func() {
				preSupply, found := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)
				err := suite.keeper.IncrementOutgoingAssetSupply(suite.ctx, tc.args.coin)
				postSupply, _ := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)

				if tc.expectPass {
					suite.True(found)
					suite.NoError(err)
					suite.Equal(preSupply.OutgoingSupply.Add(tc.args.coin), postSupply.OutgoingSupply)
				} else {
					suite.Error(err)
					suite.Equal(preSupply.OutgoingSupply, postSupply.OutgoingSupply)
				}
			},
		)
	}
}

func (suite *AssetTestSuite) TestDecrementOutgoingAssetSupply() {
	type args struct {
		coin sdk.Coin
	}
	testCases := []struct {
		name       string
		args       args
		expectPass bool
	}{{
		"normal",
		args{coin: c("htltbnb", 4)},
		true,
	}, {
		"equal outgoing",
		args{coin: c("htltbnb", 5)},
		true,
	}, {
		"exceeds outgoing",
		args{coin: c("htltbnb", 6)},
		false,
	}, {
		"unsupported asset",
		args{coin: c("htltxyz", 4)},
		false,
	}}

	for _, tc := range testCases {
		suite.SetupTest()
		suite.Run(
			tc.name,
			func() {
				preSupply, found := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)
				err := suite.keeper.DecrementOutgoingAssetSupply(suite.ctx, tc.args.coin)
				postSupply, _ := suite.keeper.GetAssetSupply(suite.ctx, tc.args.coin.Denom)

				if tc.expectPass {
					suite.True(found)
					suite.NoError(err)
					suite.True(preSupply.OutgoingSupply.Sub(tc.args.coin).IsEqual(postSupply.OutgoingSupply))
				} else {
					suite.Error(err)
					suite.Equal(preSupply.OutgoingSupply, postSupply.OutgoingSupply)
				}
			},
		)
	}
}

func (suite *AssetTestSuite) TestUpdateTimeBasedSupplyLimits() {
	type args struct {
		asset          string
		duration       time.Duration
		expectedSupply types.AssetSupply
	}
	type errArgs struct {
		expectPanic bool
		contains    string
	}
	testCases := []struct {
		name    string
		args    args
		errArgs errArgs
	}{{
		"rate-limited increment time",
		args{
			asset:          "htltinc",
			duration:       time.Second,
			expectedSupply: types.NewAssetSupply(c("htltinc", 10), c("htltinc", 5), c("htltinc", 5), c("htltinc", 0), time.Second),
		},
		errArgs{
			expectPanic: false,
			contains:    "",
		},
	}, {
		"rate-limited increment time half",
		args{
			asset:          "htltinc",
			duration:       time.Minute * 30,
			expectedSupply: types.NewAssetSupply(c("htltinc", 10), c("htltinc", 5), c("htltinc", 5), c("htltinc", 0), time.Minute*30),
		},
		errArgs{
			expectPanic: false,
			contains:    "",
		},
	}, {
		"rate-limited period change",
		args{
			asset:          "htltinc",
			duration:       time.Hour + time.Second,
			expectedSupply: types.NewAssetSupply(c("htltinc", 10), c("htltinc", 5), c("htltinc", 5), c("htltinc", 0), time.Duration(0)),
		},
		errArgs{
			expectPanic: false,
			contains:    "",
		},
	}, {
		"rate-limited period change exact",
		args{
			asset:          "htltinc",
			duration:       time.Hour,
			expectedSupply: types.NewAssetSupply(c("htltinc", 10), c("htltinc", 5), c("htltinc", 5), c("htltinc", 0), time.Duration(0)),
		},
		errArgs{
			expectPanic: false,
			contains:    "",
		},
	}, {
		"rate-limited period change big",
		args{
			asset:          "htltinc",
			duration:       time.Hour * 4,
			expectedSupply: types.NewAssetSupply(c("htltinc", 10), c("htltinc", 5), c("htltinc", 5), c("htltinc", 0), time.Duration(0)),
		},
		errArgs{
			expectPanic: false,
			contains:    "",
		},
	}, {
		"non rate-limited increment time",
		args{
			asset:          "htltbnb",
			duration:       time.Second,
			expectedSupply: types.NewAssetSupply(c("htltbnb", 5), c("htltbnb", 5), c("htltbnb", 40), c("htltbnb", 0), time.Duration(0)),
		},
		errArgs{
			expectPanic: false,
			contains:    "",
		},
	}, {
		"new asset increment time",
		args{
			asset:          "htltlol",
			duration:       time.Second,
			expectedSupply: types.NewAssetSupply(c("htltlol", 0), c("htltlol", 0), c("htltlol", 0), c("htltlol", 0), time.Second),
		},
		errArgs{
			expectPanic: false,
			contains:    "",
		},
	}}

	for _, tc := range testCases {
		suite.SetupTest()
		suite.Run(
			tc.name,
			func() {
				newParams := types.Params{
					AssetParams: []types.AssetParam{{
						Denom: "htltbnb",
						SupplyLimit: types.SupplyLimit{
							Limit:          sdk.NewInt(350000000000000),
							TimeLimited:    false,
							TimeBasedLimit: sdk.ZeroInt(),
							TimePeriod:     time.Hour,
						},
						Active:        true,
						DeputyAddress: TestDeputy.String(),
						FixedFee:      sdk.NewInt(1000),
						MinSwapAmount: sdk.OneInt(),
						MaxSwapAmount: sdk.NewInt(1000000000000),
						MinBlockLock:  MinTimeLock,
						MaxBlockLock:  MaxTimeLock,
					}, {
						Denom: "htltinc",
						SupplyLimit: types.SupplyLimit{
							Limit:          sdk.NewInt(100),
							TimeLimited:    true,
							TimeBasedLimit: sdk.NewInt(10),
							TimePeriod:     time.Hour,
						},
						Active:        false,
						DeputyAddress: TestDeputy.String(),
						FixedFee:      sdk.NewInt(1000),
						MinSwapAmount: sdk.OneInt(),
						MaxSwapAmount: sdk.NewInt(1000000000000),
						MinBlockLock:  MinTimeLock,
						MaxBlockLock:  MaxTimeLock,
					}, {
						Denom: "htltlol",
						SupplyLimit: types.SupplyLimit{
							Limit:          sdk.NewInt(100),
							TimeLimited:    true,
							TimeBasedLimit: sdk.NewInt(10),
							TimePeriod:     time.Hour,
						},
						Active:        false,
						DeputyAddress: TestDeputy.String(),
						FixedFee:      sdk.NewInt(1000),
						MinSwapAmount: sdk.OneInt(),
						MaxSwapAmount: sdk.NewInt(1000000000000),
						MinBlockLock:  MinTimeLock,
						MaxBlockLock:  MaxTimeLock,
					}},
				}
				suite.keeper.SetParams(suite.ctx, newParams)
				suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(tc.args.duration))
				suite.NotPanics(
					func() {
						suite.keeper.UpdateTimeBasedSupplyLimits(suite.ctx)
					},
				)
				if !tc.errArgs.expectPanic {
					supply, found := suite.keeper.GetAssetSupply(suite.ctx, tc.args.asset)
					suite.Require().True(found)
					suite.Require().Equal(tc.args.expectedSupply, supply)
				}
			},
		)
	}
}
