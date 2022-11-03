package keeper_test

import (
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irismod/modules/farm/keeper"
	"github.com/irisnet/irismod/modules/farm/types"
	"github.com/irisnet/irismod/simapp"
)

var (
	testInitCoinAmt     = sdk.NewInt(100000000_000_000)
	testPoolDescription = "USDT/IRIS Farm Pool"
	testBeginHeight     = int64(1)
	testLPTokenDenom    = sdk.DefaultBondDenom
	testRewardPerBlock  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000)))
	testTotalReward     = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000_000_000)))
	testDestructible    = true

	testCreator sdk.AccAddress
	testFarmer1 sdk.AccAddress
	testFarmer2 sdk.AccAddress
	testFarmer3 sdk.AccAddress
	testFarmer4 sdk.AccAddress
	testFarmer5 sdk.AccAddress

	isCheckTx = false
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    codec.BinaryCodec
	ctx    sdk.Context
	keeper *keeper.Keeper
	app    *simapp.SimApp
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(suite.T(), isCheckTx)
	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{Height: 1})
	suite.app = app
	suite.keeper = &app.FarmKeeper
	suite.keeper.SetParams(suite.ctx, types.DefaultParams())
	suite.setTestAddrs()
}

func (suite *KeeperTestSuite) setTestAddrs() {
	testAddrs := simapp.AddTestAddrs(suite.app, suite.ctx, 6, testInitCoinAmt)

	testCreator = testAddrs[0]
	testFarmer1 = testAddrs[1]
	testFarmer2 = testAddrs[2]
	testFarmer3 = testAddrs[3]
	testFarmer4 = testAddrs[4]
	testFarmer5 = testAddrs[5]
}

func (suite *KeeperTestSuite) TestCreatePool() {
	ctx := suite.ctx
	pool, err := suite.keeper.CreatePool(
		ctx,
		testPoolDescription,
		testLPTokenDenom,
		testBeginHeight,
		testRewardPerBlock,
		testTotalReward,
		testDestructible,
		testCreator,
	)
	suite.Require().NoError(err)

	suite.Require().NotEmpty(pool.Id)
	suite.Require().Equal(testPoolDescription, pool.Description)
	suite.Require().Equal(testLPTokenDenom, pool.TotalLptLocked.Denom)
	suite.Require().Equal(testBeginHeight, pool.StartHeight)
	suite.Require().Equal(testDestructible, pool.Editable)
	suite.Require().Equal(testCreator.String(), pool.Creator)

	//check reward rules
	rules := suite.keeper.GetRewardRules(ctx, pool.Id)
	suite.Require().Len(rules, len(testRewardPerBlock))

	for _, r := range rules {
		suite.Require().Equal(testTotalReward.AmountOf(r.Reward), r.RemainingReward)
		suite.Require().Equal(testTotalReward.AmountOf(r.Reward), r.TotalReward)
		suite.Require().Equal(testRewardPerBlock.AmountOf(r.Reward), r.RewardPerBlock)
		suite.Require().Equal(sdk.ZeroDec(), r.RewardPerShare)
	}

	pool.Rules = rules
	endHeight, _ := pool.ExpiredHeight()
	suite.Require().Equal(endHeight, pool.EndHeight)

	//check queue
	suite.keeper.IteratorExpiredPool(ctx, pool.EndHeight, func(pool types.FarmPool) {
		suite.Require().Equal(pool.Id, pool.Id)
	})

	//check balance
	expectedBal := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, testInitCoinAmt)).
		Sub(sdk.NewCoins(suite.keeper.CreatePoolFee(ctx))...).
		Sub(testTotalReward...)
	actualBal := suite.app.BankKeeper.GetAllBalances(ctx, testCreator)
	suite.Require().Equal(expectedBal, actualBal)
}

func (suite *KeeperTestSuite) TestDestroyPool() {
	ctx := suite.ctx
	pool, err := suite.keeper.CreatePool(ctx,
		testPoolDescription,
		testLPTokenDenom,
		testBeginHeight,
		testRewardPerBlock,
		testTotalReward,
		testDestructible,
		testCreator,
	)
	suite.Require().NoError(err)

	newCtx := suite.app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		Height: 10,
	})
	_, err = suite.keeper.DestroyPool(newCtx, pool.Id, testCreator)
	suite.Require().NoError(err)

	//check farm pool
	p, exist := suite.keeper.GetPool(newCtx, pool.Id)
	suite.Require().True(exist)

	suite.Require().EqualValues(newCtx.BlockHeight(), p.LastHeightDistrRewards)
	suite.Require().EqualValues(newCtx.BlockHeight(), p.EndHeight)

	//check balance
	expectedBal := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, testInitCoinAmt)).
		Sub(sdk.NewCoins(suite.keeper.CreatePoolFee(ctx))...)
	actualBal := suite.app.BankKeeper.GetAllBalances(ctx, testCreator)
	suite.Require().Equal(expectedBal, actualBal)

	rewardAdded := sdk.NewCoins(
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10_000_000)),
		sdk.NewCoin("uiris", sdk.NewInt(10_000_000)),
	)
	err = suite.keeper.AdjustPool(newCtx,
		p.Id,
		rewardAdded,
		nil,
		testCreator,
	)
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestAppendReward() {
	ctx := suite.ctx
	pool, err := suite.keeper.CreatePool(
		ctx,
		testPoolDescription,
		testLPTokenDenom,
		testBeginHeight,
		testRewardPerBlock,
		testTotalReward,
		testDestructible,
		testCreator,
	)
	suite.Require().NoError(err)

	//check farm pool
	p, exist := suite.keeper.GetPool(ctx, pool.Id)
	suite.Require().True(exist)

	//panic with adding new token as reward
	rewardAdded := sdk.NewCoins(
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10_000_000)),
		sdk.NewCoin("uiris", sdk.NewInt(10_000_000)),
	)
	err = suite.keeper.AdjustPool(ctx,
		p.Id,
		rewardAdded,
		nil,
		testCreator,
	)
	suite.Require().Error(err)

	rewardAdded = sdk.NewCoins(
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10_000_000)),
	)
	err = suite.keeper.AdjustPool(ctx,
		p.Id,
		rewardAdded,
		nil,
		testCreator,
	)
	suite.Require().NoError(err)

	//check farm pool
	pool2, exist := suite.keeper.GetPool(ctx, p.Id)
	suite.Require().True(exist)
	suite.Require().EqualValues(p.EndHeight+10, pool2.EndHeight)

	//check reward rules
	rules := suite.keeper.GetRewardRules(ctx, pool.Id)
	suite.Require().Len(rules, len(testRewardPerBlock))

	for _, r := range rules {
		suite.Equal(
			testTotalReward.AmountOf(r.Reward).Add(rewardAdded.AmountOf(r.Reward)),
			r.RemainingReward,
		)
		suite.Equal(
			testTotalReward.AmountOf(r.Reward).Add(rewardAdded.AmountOf(r.Reward)),
			r.TotalReward,
		)
	}
}

func (suite *KeeperTestSuite) TestStake() {
	ctx := suite.ctx
	pool, err := suite.keeper.CreatePool(
		ctx,
		testPoolDescription,
		testLPTokenDenom,
		testBeginHeight,
		testRewardPerBlock,
		testTotalReward,
		testDestructible,
		testCreator,
	)
	suite.Require().NoError(err)

	lpToken := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))
	type args struct {
		height         int64
		stakeCoin      sdk.Coin
		locked         sdk.Int
		expectReward   sdk.Coins
		debt           sdk.Coins
		rewardPerShare sdk.Dec
	}

	var testcase = []args{
		{
			height:         100,
			stakeCoin:      lpToken,
			locked:         lpToken.Amount.MulRaw(1),
			expectReward:   nil,
			debt:           nil,
			rewardPerShare: sdk.ZeroDec(),
		},
		{
			height:         200,
			stakeCoin:      lpToken,
			locked:         lpToken.Amount.MulRaw(2),
			expectReward:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))),
			debt:           sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(200_000_000))),
			rewardPerShare: sdk.NewDecFromIntWithPrec(sdk.NewInt(1), 0),
		},
		{
			height:         300,
			stakeCoin:      lpToken,
			locked:         lpToken.Amount.MulRaw(3),
			expectReward:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))),
			debt:           sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(450_000_000))),
			rewardPerShare: sdk.NewDecFromIntWithPrec(sdk.NewInt(15), 1),
		},
		{
			height:         400,
			stakeCoin:      lpToken,
			locked:         lpToken.Amount.MulRaw(4),
			expectReward:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(99999999))),
			debt:           sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(733_333_333))),
			rewardPerShare: sdk.NewDecFromIntWithPrec(sdk.NewInt(1_833_333_333_333_333_333), 18),
		},
	}

	for _, tc := range testcase {
		suite.AssertStake(
			pool.Id,
			tc.height,
			tc.stakeCoin,
			tc.locked,
			tc.expectReward,
			tc.debt,
			tc.rewardPerShare)
	}
}

func (suite *KeeperTestSuite) TestUnstake() {
	ctx := suite.ctx
	pool, err := suite.keeper.CreatePool(
		ctx,
		testPoolDescription,
		testLPTokenDenom,
		testBeginHeight,
		testRewardPerBlock,
		testTotalReward,
		testDestructible,
		testCreator,
	)
	suite.Require().NoError(err)

	lpToken := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))
	suite.AssertStake(pool.Id, 100,
		lpToken,
		lpToken.Amount,
		nil,
		nil,
		sdk.ZeroDec(),
	)
	suite.AssertUnstake(pool.Id, 200,
		lpToken,
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))),
		nil,
		sdk.NewDecFromIntWithPrec(sdk.NewInt(1), 0),
		true,
	)
	suite.AssertStake(pool.Id, 300,
		lpToken,
		lpToken.Amount,
		nil,
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))),
		sdk.NewDecFromIntWithPrec(sdk.NewInt(1), 0),
	)
	suite.AssertUnstake(pool.Id, 400,
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(50_000_000)),
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))),
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))),
		sdk.NewDecFromIntWithPrec(sdk.NewInt(2), 0),
		false,
	)
	suite.AssertUnstake(pool.Id, 500,
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(50_000_000)),
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))),
		nil,
		sdk.NewDecFromIntWithPrec(sdk.NewInt(4), 0),
		true,
	)
}

func (suite *KeeperTestSuite) TestHarvest() {
	ctx := suite.ctx
	pool, err := suite.keeper.CreatePool(
		ctx,
		testPoolDescription,
		testLPTokenDenom,
		testBeginHeight,
		testRewardPerBlock,
		testTotalReward,
		testDestructible,
		testCreator,
	)
	suite.Require().NoError(err)

	lpToken := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))
	suite.AssertStake(pool.Id, 100, lpToken, lpToken.Amount, nil, nil, sdk.ZeroDec())

	type args struct {
		index          int64
		height         int64
		expectReward   sdk.Coins
		debt           sdk.Coins
		rewardPerShare sdk.Dec
	}

	var testcase = []args{
		{
			index:          1,
			height:         200,
			expectReward:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))),
			debt:           sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))),
			rewardPerShare: sdk.NewDecFromIntWithPrec(sdk.NewInt(1), 0),
		},
		{
			index:          2,
			height:         300,
			expectReward:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))),
			debt:           sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(200_000_000))),
			rewardPerShare: sdk.NewDecFromIntWithPrec(sdk.NewInt(2), 0),
		},
		{
			index:          3,
			height:         400,
			expectReward:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))),
			debt:           sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(300_000_000))),
			rewardPerShare: sdk.NewDecFromIntWithPrec(sdk.NewInt(3), 0),
		},
		{
			index:          4,
			height:         500,
			expectReward:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000))),
			debt:           sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(400_000_000))),
			rewardPerShare: sdk.NewDecFromIntWithPrec(sdk.NewInt(4), 0),
		},
	}

	for _, tc := range testcase {
		suite.AssertHarvest(pool.Id, tc.index,
			tc.height,
			tc.expectReward,
			tc.debt,
			tc.rewardPerShare)
	}
}

func (suite *KeeperTestSuite) AssertStake(
	poolId string,
	height int64,
	stakeCoin sdk.Coin,
	locked sdk.Int,
	expectReward, debt sdk.Coins,
	rewardPerShare sdk.Dec,
) {
	ctx := suite.app.BaseApp.NewContext(isCheckTx, tmproto.Header{Height: height})
	reward, err := suite.keeper.Stake(ctx, poolId, stakeCoin, testFarmer1)

	suite.Require().NoError(err)
	suite.Require().Equal(expectReward, reward)

	info, exist := suite.keeper.GetFarmInfo(ctx, poolId, testFarmer1.String())
	suite.Require().True(exist)
	suite.Require().Equal(debt, info.RewardDebt)
	suite.Require().Equal(locked, info.Locked)

	//check reward rules again
	rules := suite.keeper.GetRewardRules(ctx, poolId)
	suite.Require().Len(rules, len(testRewardPerBlock))
	for _, r := range rules {
		suite.Require().Equal(rewardPerShare, r.RewardPerShare)
	}
}

func (suite *KeeperTestSuite) AssertUnstake(
	poolId string,
	height int64,
	unstakeCoin sdk.Coin,
	expectReward, expectDebt sdk.Coins,
	rewardPerShare sdk.Dec,
	unstakeAll bool,
) {
	ctx := suite.app.BaseApp.NewContext(isCheckTx, tmproto.Header{Height: height})

	//check farm pool
	poolSrc, _ := suite.keeper.GetPool(ctx, poolId)
	//check farm information
	farmInfoSrc, _ := suite.keeper.GetFarmInfo(ctx, poolId, testFarmer1.String())

	reward, err := suite.keeper.Unstake(ctx, poolId, unstakeCoin, testFarmer1)
	suite.Require().NoError(err)
	suite.Require().Equal(expectReward, reward)

	//check farm information
	farmInfo, exist := suite.keeper.GetFarmInfo(ctx, poolId, testFarmer1.String())
	if unstakeAll {
		suite.Require().False(exist)
	} else {
		suite.Require().True(exist)
		suite.Require().Equal(farmInfoSrc.Locked.Sub(unstakeCoin.Amount), farmInfo.Locked)
		suite.Require().Equal(expectDebt, farmInfo.RewardDebt)
	}

	//check farm pool
	pool, exist := suite.keeper.GetPool(ctx, poolId)
	suite.Require().True(exist)
	suite.Require().Equal(pool.TotalLptLocked.String(), poolSrc.TotalLptLocked.Sub(unstakeCoin).String())

	//check reward rules again
	rules := suite.keeper.GetRewardRules(ctx, poolId)
	suite.Require().Len(rules, len(testRewardPerBlock))
	for _, r := range rules {
		suite.Require().Equal(rewardPerShare, r.RewardPerShare)
	}
}

func (suite *KeeperTestSuite) AssertHarvest(
	poolId string,
	index int64,
	height int64,
	expectReward sdk.Coins,
	debt sdk.Coins,
	rewardPerShare sdk.Dec,
) {
	ctx := suite.app.BaseApp.NewContext(isCheckTx, tmproto.Header{Height: height})
	reward, err := suite.keeper.Harvest(ctx, poolId, testFarmer1)

	suite.Require().NoError(err)
	suite.Require().Equal(expectReward, reward)

	info, exist := suite.keeper.GetFarmInfo(ctx, poolId, testFarmer1.String())
	suite.Require().True(exist)
	suite.Require().Equal(debt, info.RewardDebt)

	//check reward rules again
	rules := suite.keeper.GetRewardRules(ctx, poolId)
	suite.Require().Len(rules, len(testRewardPerBlock))
	for _, r := range rules {
		suite.Require().Equal(rewardPerShare, r.RewardPerShare)
	}
}
