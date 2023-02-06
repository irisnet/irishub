package cli_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	farmcli "github.com/irisnet/irismod/modules/farm/client/cli"
	"github.com/irisnet/irismod/modules/farm/client/testutil"
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

func (s *IntegrationTestSuite) TestFarm() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	// ---------------------------------------------------------------------------

	creator := val.Address
	description := "iris-atom farm pool"
	startHeight := s.LatestHeight() + 1
	rewardPerBlock := sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10)))
	totalReward := sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(1000)))
	editable := true

	globalFlags := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}

	args := []string{
		fmt.Sprintf("--%s=%s", farmcli.FlagDescription, description),
		fmt.Sprintf("--%s=%d", farmcli.FlagStartHeight, startHeight),
		fmt.Sprintf("--%s=%s", farmcli.FlagRewardPerBlock, rewardPerBlock),
		fmt.Sprintf("--%s=%s", farmcli.FlagLPTokenDenom, s.network.BondDenom),
		fmt.Sprintf("--%s=%s", farmcli.FlagTotalReward, totalReward),
		fmt.Sprintf("--%s=%v", farmcli.FlagEditable, editable),
	}

	args = append(args, globalFlags...)
	txResult := testutil.CreateFarmPoolExec(
		s.T(),
		s.network,
		clientCtx,
		creator.String(),
		args...,
	)

	poolId := s.network.GetAttribute(farmtypes.EventTypeCreatePool, farmtypes.AttributeValuePoolId, txResult.Events)
	expectedContents := &farmtypes.FarmPoolEntry{
		Id:              poolId,
		Creator:         creator.String(),
		Description:     description,
		StartHeight:     startHeight,
		EndHeight:       startHeight + 100,
		Editable:        editable,
		Expired:         false,
		TotalLptLocked:  sdk.NewCoin(s.network.BondDenom, sdk.ZeroInt()),
		TotalReward:     totalReward,
		RemainingReward: totalReward,
		RewardPerBlock:  rewardPerBlock,
	}

	respType := testutil.QueryFarmPoolExec(s.T(), s.network, val.ClientCtx, poolId)
	s.Require().EqualValues(expectedContents, respType.Pool)

	reward := sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(1000)))
	args = []string{
		fmt.Sprintf("--%s=%v", farmcli.FlagAdditionalReward, reward.String()),
	}
	args = append(args, globalFlags...)
	txResult = testutil.AppendRewardExec(
		s.T(),
		s.network,
		clientCtx,
		creator.String(),
		poolId,
		args...,
	)

	lpToken := sdk.NewCoin(s.network.BondDenom, sdk.NewInt(100))
	txResult = testutil.StakeExec(
		s.T(),
		s.network,
		clientCtx,
		creator.String(),
		poolId,
		lpToken.String(),
		globalFlags...,
	)
	beginHeight := txResult.Height

	unstakeLPToken := sdk.NewCoin(s.network.BondDenom, sdk.NewInt(50))
	txResult = testutil.UnstakeExec(
		s.T(),
		s.network,
		clientCtx,
		creator.String(),
		poolId,
		unstakeLPToken.String(),
		globalFlags...,
	)
	endHeight := txResult.Height

	rewardGot := s.network.GetAttribute(farmtypes.EventTypeUnstake, farmtypes.AttributeValueReward, txResult.Events)
	expectedReward := rewardPerBlock.MulInt(sdk.NewInt(endHeight - beginHeight))
	s.Require().Equal(expectedReward.String(), rewardGot)

	txResult = testutil.HarvestExec(
		s.T(),
		s.network,
		clientCtx,
		creator.String(),
		poolId,
		globalFlags...,
	)
	endHeight1 := txResult.Height

	rewardGot = s.network.GetAttribute(farmtypes.EventTypeHarvest, farmtypes.AttributeValueReward, txResult.Events)
	expectedReward = rewardPerBlock.MulInt(sdk.NewInt(endHeight1 - endHeight))
	s.Require().Equal(expectedReward.String(), rewardGot)

	queryFarmerArgs := []string{
		fmt.Sprintf("--%s=%s", farmcli.FlagFarmPool, poolId),
	}

	leftlpToken := lpToken.Sub(unstakeLPToken)
	response := testutil.QueryFarmerExec(
		s.T(),
		s.network,
		val.ClientCtx, creator.String(), queryFarmerArgs...)
	s.Require().EqualValues(leftlpToken, response.List[0].Locked)

	txResult = testutil.DestroyExec(
		s.T(),
		s.network,
		clientCtx,
		creator.String(),
		poolId,
		globalFlags...,
	)
}

func (s *IntegrationTestSuite) LatestHeight() int64 {
	height, err := s.network.LatestHeight()
	s.Require().NoError(err)
	return height
}
