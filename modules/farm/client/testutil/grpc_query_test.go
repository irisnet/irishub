package testutil_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"

	farmcli "github.com/irisnet/irismod/modules/farm/client/cli"
	farmtestutil "github.com/irisnet/irismod/modules/farm/client/testutil"
	farmtypes "github.com/irisnet/irismod/modules/farm/types"
	"github.com/irisnet/irismod/simapp"
)

type IntegrationTestSuite struct {
	suite.Suite

	network simapp.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = simapp.SetupNetwork(s.T())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestRest() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	baseURL := val.APIAddress

	// ---------------------------------------------------------------------------

	creator := val.Address
	description := "iris-atom farm pool"
	startHeight := s.LatestHeight() + 1
	rewardPerBlock := sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10)))
	lpTokenDenom := s.network.BondDenom
	totalReward := sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(1000)))
	editable := true

	globalFlags := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	args := []string{
		fmt.Sprintf("--%s=%s", farmcli.FlagDescription, description),
		fmt.Sprintf("--%s=%d", farmcli.FlagStartHeight, startHeight),
		fmt.Sprintf("--%s=%s", farmcli.FlagRewardPerBlock, rewardPerBlock),
		fmt.Sprintf("--%s=%s", farmcli.FlagLPTokenDenom, lpTokenDenom),
		fmt.Sprintf("--%s=%s", farmcli.FlagTotalReward, totalReward),
		fmt.Sprintf("--%s=%v", farmcli.FlagEditable, editable),
	}

	args = append(args, globalFlags...)
	txResult := farmtestutil.CreateFarmPoolExec(
		s.T(),
		s.network,
		clientCtx,
		creator.String(),
		args...,
	)

	poolId := s.network.GetAttribute(
		farmtypes.EventTypeCreatePool,
		farmtypes.AttributeValuePoolId,
		txResult.Events,
	)
	expectedContents := farmtypes.FarmPoolEntry{
		Id:              poolId,
		Description:     description,
		Creator:         creator.String(),
		StartHeight:     startHeight,
		EndHeight:       startHeight + 100,
		Editable:        editable,
		Expired:         false,
		TotalLptLocked:  sdk.NewCoin(lpTokenDenom, sdk.ZeroInt()),
		TotalReward:     totalReward,
		RemainingReward: totalReward,
		RewardPerBlock:  rewardPerBlock,
	}

	respType := proto.Message(&farmtypes.QueryFarmPoolsResponse{})
	queryPoolURL := fmt.Sprintf("%s/irismod/farm/pools", baseURL)
	resp, err := testutil.GetRequest(queryPoolURL)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	result := respType.(*farmtypes.QueryFarmPoolsResponse)
	s.Require().EqualValues(expectedContents, *result.Pools[0])

	_, err = s.network.WaitForHeight(startHeight)
	s.Require().NoError(err)
	s.network.WaitForNextBlock()

	lpToken := sdk.NewCoin(s.network.BondDenom, sdk.NewInt(100))
	txResult = farmtestutil.StakeExec(
		s.T(),
		s.network,
		clientCtx,
		creator.String(),
		poolId,
		lpToken.String(),
		globalFlags...,
	)

	expectFarmer := farmtypes.LockedInfo{
		PoolId:        poolId,
		Locked:        lpToken,
		PendingReward: sdk.Coins{},
	}

	queryFarmerRespType := proto.Message(&farmtypes.QueryFarmerResponse{})
	queryFarmInfoURL := fmt.Sprintf("%s/irismod/farm/farmers/%s", baseURL, creator.String())
	resp, err = testutil.GetRequest(queryFarmInfoURL)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, queryFarmerRespType))
	farmer := queryFarmerRespType.(*farmtypes.QueryFarmerResponse)

	if farmer.Height-txResult.Height > 0 {
		expectFarmer.PendingReward = rewardPerBlock.MulInt(
			sdk.NewInt((farmer.Height - txResult.Height)),
		)
	}
	s.Require().EqualValues(expectFarmer, *farmer.List[0])
}

func (s *IntegrationTestSuite) LatestHeight() int64 {
	height, err := s.network.LatestHeight()
	s.Require().NoError(err)
	return height
}
