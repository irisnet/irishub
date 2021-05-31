package rest_test

import (
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

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

func (s *IntegrationTestSuite) TestRest() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	baseURL := val.APIAddress

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

	expectedContents := farmtypes.FarmPoolEntry{
		Name:               farmPool,
		Creator:            creator.String(),
		StartHeight:        uint64(startHeight),
		EndHeight:          uint64(startHeight + 101),
		Destructible:       destructible,
		Expired:            false,
		TotalLpTokenLocked: sdk.NewCoin(lpTokenDenom, sdk.ZeroInt()),
		TotalReward:        totalReward,
		RemainingReward:    totalReward,
		RewardPerBlock:     rewardPerBlock,
	}

	respType = proto.Message(&farmtypes.QueryPoolsResponse{})
	queryPoolURL := fmt.Sprintf("%s/irismod/farm/pools", baseURL)
	resp, err := rest.GetRequest(queryPoolURL)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(resp, respType))
	pools := respType.(*farmtypes.QueryPoolsResponse)
	s.Require().EqualValues(expectedContents, *pools.List[0])

	_, err = s.network.WaitForHeight(startHeight)
	s.Require().NoError(err)

	respType = proto.Message(&sdk.TxResponse{})
	lpToken := sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))
	bz, err = testutil.StakeExec(clientCtx,
		creator.String(),
		farmPool,
		lpToken.String(),
		globalFlags...,
	)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	expectFarmer := farmtypes.LockedInfo{
		PoolName:      farmPool,
		Locked:        lpToken,
		PendingReward: sdk.Coins{},
	}

	queryFarmerRespType := proto.Message(&farmtypes.QueryFarmerResponse{})
	queryFarmInfoURL := fmt.Sprintf("%s/irismod/farm/farmers/%s", baseURL, creator.String())
	resp, err = rest.GetRequest(queryFarmInfoURL)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(resp, queryFarmerRespType))
	farmer := queryFarmerRespType.(*farmtypes.QueryFarmerResponse)
	s.Require().EqualValues(expectFarmer, *farmer.List[0])
}

func (s *IntegrationTestSuite) LatestHeight() int64 {
	height, err := s.network.LatestHeight()
	s.Require().NoError(err)
	return height
}
