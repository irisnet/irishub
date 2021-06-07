package cli_test

import (
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	farmcli "github.com/irisnet/irismod/modules/farm/client/cli"
	"github.com/irisnet/irismod/modules/farm/client/testutil"
	farmtypes "github.com/irisnet/irismod/modules/farm/types"
	"github.com/irisnet/irismod/simapp"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := simapp.NewConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
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
	rewardPerBlock := sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)))
	lpTokenDenom := s.cfg.BondDenom
	totalReward := sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1000)))
	destructible := true
	farmPool := "iris-atom"

	globalFlags := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	args := []string{
		fmt.Sprintf("--%s=%s", farmcli.FlagDescription, description),
		fmt.Sprintf("--%s=%d", farmcli.FlagStartHeight, startHeight),
		fmt.Sprintf("--%s=%s", farmcli.FlagRewardPerBlock, rewardPerBlock),
		fmt.Sprintf("--%s=%s", farmcli.FlagLPTokenDenom, lpTokenDenom),
		fmt.Sprintf("--%s=%s", farmcli.FlagTotalReward, totalReward),
		fmt.Sprintf("--%s=%v", farmcli.FlagDestructible, destructible),
	}

	args = append(args, globalFlags...)
	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := testutil.CreateFarmPoolExec(clientCtx,
		creator.String(),
		farmPool,
		args...,
	)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&farmtypes.QueryPoolsResponse{})
	expectedContents := farmtypes.FarmPoolEntry{
		Name:               farmPool,
		Creator:            creator.String(),
		Description:        description,
		StartHeight:        uint64(startHeight),
		EndHeight:          uint64(startHeight + 100),
		Destructible:       destructible,
		Expired:            false,
		TotalLpTokenLocked: sdk.NewCoin(lpTokenDenom, sdk.ZeroInt()),
		TotalReward:        totalReward,
		RemainingReward:    totalReward,
		RewardPerBlock:     rewardPerBlock,
	}

	bz, err = testutil.QueryFarmPoolExec(val.ClientCtx, farmPool)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	pools := respType.(*farmtypes.QueryPoolsResponse)
	s.Require().EqualValues(expectedContents, *pools.List[0])

	respType = proto.Message(&sdk.TxResponse{})
	reward := sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1000)))
	bz, err = testutil.AppendRewardExec(clientCtx,
		creator.String(),
		farmPool,
		reward.String(),
		globalFlags...,
	)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	_, err = s.network.WaitForHeight(startHeight)
	s.Require().NoError(err)

	lpToken := sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))
	bz, err = testutil.StakeExec(clientCtx,
		creator.String(),
		farmPool,
		lpToken.String(),
		globalFlags...,
	)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	s.Require().Equal(expectedCode, txResp.Code)

	unstakeLPToken := sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(50))
	bz, err = testutil.UnstakeExec(clientCtx,
		creator.String(),
		farmPool,
		unstakeLPToken.String(),
		globalFlags...,
	)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	s.Require().Equal(expectedCode, txResp.Code)
	rewardGot := gjson.Get(txResp.RawLog, "0.events.2.attributes.3.value").String()
	s.Require().Equal(rewardPerBlock.String(), rewardGot)

	bz, err = testutil.HarvestExec(clientCtx,
		creator.String(),
		farmPool,
		globalFlags...,
	)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	s.Require().Equal(expectedCode, txResp.Code)
	rewardGot = gjson.Get(txResp.RawLog, "0.events.0.attributes.2.value").String()
	s.Require().Equal(rewardPerBlock.String(), rewardGot)

	queryFarmerArgs := []string{
		fmt.Sprintf("--%s=%s", farmcli.FlagFarmPool, farmPool),
	}
	expectFarmer := farmtypes.LockedInfo{
		PoolName:      farmPool,
		Locked:        lpToken.Sub(unstakeLPToken),
		PendingReward: sdk.Coins{},
	}

	queryFarmerRespType := proto.Message(&farmtypes.QueryFarmerResponse{})
	bz, err = testutil.QueryFarmerExec(val.ClientCtx, creator.String(), queryFarmerArgs...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), queryFarmerRespType))
	result := queryFarmerRespType.(*farmtypes.QueryFarmerResponse)
	s.Require().EqualValues(expectFarmer, *result.List[0])

	bz, err = testutil.DestroyExec(clientCtx,
		creator.String(),
		farmPool,
		globalFlags...,
	)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	s.Require().Equal(expectedCode, txResp.Code)

}

func (s *IntegrationTestSuite) LatestHeight() int64 {
	height, err := s.network.LatestHeight()
	s.Require().NoError(err)
	return height
}
