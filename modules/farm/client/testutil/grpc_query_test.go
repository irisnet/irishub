package testutil_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/cosmos/gogoproto/proto"
// 	"github.com/stretchr/testify/suite"

// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	"github.com/cosmos/cosmos-sdk/testutil"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
// 	tokentypes "github.com/irisnet/irismod/modules/token/types/v1"
// 	"github.com/irisnet/irismod/simapp"
// 	farmcli "github.com/irisnet/irismod/farm/client/cli"
// 	farmtestutil "github.com/irisnet/irismod/farm/client/testutil"
// 	farmtypes "github.com/irisnet/irismod/farm/types"
// )

// type IntegrationTestSuite struct {
// 	suite.Suite

// 	network simapp.Network
// }

// func (s *IntegrationTestSuite) SetupSuite() {
// 	s.T().Log("setting up integration test suite")

// 	s.network = simapp.SetupNetwork(s.T())
// 	sdk.SetCoinDenomRegex(func() string {
// 		return `[a-zA-Z][a-zA-Z0-9/\-]{2,127}`
// 	})
// }

// func (s *IntegrationTestSuite) TearDownSuite() {
// 	s.T().Log("tearing down integration test suite")
// 	s.network.Cleanup()
// }

// func TestIntegrationTestSuite(t *testing.T) {
// 	suite.Run(t, new(IntegrationTestSuite))
// }

// func (s *IntegrationTestSuite) TestRest() {
// 	val := s.network.Validators[0]
// 	clientCtx := val.ClientCtx
// 	baseURL := val.APIAddress

// 	s.Init()

// 	// ---------------------------------------------------------------------------

// 	creator := val.Address
// 	description := "iris-atom farm pool"
// 	startHeight := s.LatestHeight() + 1
// 	rewardPerBlock := sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10)))
// 	lpTokenDenom := "lpt-1"
// 	totalReward := sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(1000)))
// 	editable := true

// 	globalFlags := []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	args := []string{
// 		fmt.Sprintf("--%s=%s", farmcli.FlagDescription, description),
// 		fmt.Sprintf("--%s=%d", farmcli.FlagStartHeight, startHeight),
// 		fmt.Sprintf("--%s=%s", farmcli.FlagRewardPerBlock, rewardPerBlock),
// 		fmt.Sprintf("--%s=%s", farmcli.FlagLPTokenDenom, lpTokenDenom),
// 		fmt.Sprintf("--%s=%s", farmcli.FlagTotalReward, totalReward),
// 		fmt.Sprintf("--%s=%v", farmcli.FlagEditable, editable),
// 	}

// 	args = append(args, globalFlags...)
// 	txResult := farmtestutil.CreateFarmPoolExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		creator.String(),
// 		args...,
// 	)
// 	s.Require().EqualValues(txResult.Code, 0, txResult.Log)

// 	poolId := s.network.GetAttribute(
// 		farmtypes.EventTypeCreatePool,
// 		farmtypes.AttributeValuePoolId,
// 		txResult.Events,
// 	)
// 	expectedContents := farmtypes.FarmPoolEntry{
// 		Id:              poolId,
// 		Description:     description,
// 		Creator:         creator.String(),
// 		StartHeight:     startHeight,
// 		EndHeight:       startHeight + 100,
// 		Editable:        editable,
// 		Expired:         false,
// 		TotalLptLocked:  sdk.NewCoin(lpTokenDenom, sdk.ZeroInt()),
// 		TotalReward:     totalReward,
// 		RemainingReward: totalReward,
// 		RewardPerBlock:  rewardPerBlock,
// 	}

// 	respType := proto.Message(&farmtypes.QueryFarmPoolsResponse{})
// 	queryPoolURL := fmt.Sprintf("%s/irismod/farm/pools", baseURL)
// 	resp, err := testutil.GetRequest(queryPoolURL)

// 	s.Require().NoError(err)
// 	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
// 	result := respType.(*farmtypes.QueryFarmPoolsResponse)
// 	s.Require().EqualValues(expectedContents, *result.Pools[0])

// 	_, err = s.network.WaitForHeight(startHeight)
// 	s.Require().NoError(err)
// 	s.network.WaitForNextBlock()

// 	lpToken := sdk.NewCoin(lpTokenDenom, sdk.NewInt(100))
// 	txResult = farmtestutil.StakeExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		creator.String(),
// 		poolId,
// 		lpToken.String(),
// 		globalFlags...,
// 	)
// 	s.Require().Equal(uint32(0), txResult.Code, txResult.Log)

// 	expectFarmer := farmtypes.LockedInfo{
// 		PoolId:        poolId,
// 		Locked:        lpToken,
// 		PendingReward: sdk.Coins{},
// 	}

// 	queryFarmerRespType := proto.Message(&farmtypes.QueryFarmerResponse{})
// 	queryFarmInfoURL := fmt.Sprintf("%s/irismod/farm/farmers/%s", baseURL, creator.String())
// 	resp, err = testutil.GetRequest(queryFarmInfoURL)
// 	s.Require().NoError(err)
// 	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, queryFarmerRespType))
// 	farmer := queryFarmerRespType.(*farmtypes.QueryFarmerResponse)

// 	if farmer.Height-txResult.Height > 0 {
// 		expectFarmer.PendingReward = rewardPerBlock.MulInt(
// 			sdk.NewInt((farmer.Height - txResult.Height)),
// 		)
// 	}
// 	s.Require().EqualValues(expectFarmer, *farmer.List[0])
// }

// func (s *IntegrationTestSuite) LatestHeight() int64 {
// 	height, err := s.network.LatestHeight()
// 	s.Require().NoError(err)
// 	return height
// }

// func (s *IntegrationTestSuite) Init() {

// 	val := s.network.Validators[0]
// 	clientCtx := val.ClientCtx

// 	from := val.Address
// 	symbol := "kitty"
// 	name := "Kitty Token"
// 	minUnit := "kitty"
// 	scale := uint32(0)
// 	initialSupply := uint64(100000000)
// 	maxSupply := uint64(200000000)
// 	mintable := true

// 	// issue token
// 	msgIssueToken := &tokentypes.MsgIssueToken{
// 		Symbol:        symbol,
// 		Name:          name,
// 		Scale:         scale,
// 		MinUnit:       minUnit,
// 		InitialSupply: initialSupply,
// 		MaxSupply:     maxSupply,
// 		Mintable:      mintable,
// 		Owner:         from.String(),
// 	}
// 	res := s.network.BlockSendMsgs(s.T(), msgIssueToken)
// 	s.Require().Equal(uint32(0), res.Code, res.Log)

// 	// add liquidity
// 	status, err := clientCtx.Client.Status(context.Background())
// 	s.Require().NoError(err)
// 	deadline := status.SyncInfo.LatestBlockTime.Add(time.Minute)

// 	msgAddLiquidity := &coinswaptypes.MsgAddLiquidity{
// 		MaxToken:         sdk.NewCoin(symbol, sdk.NewInt(1000)),
// 		ExactStandardAmt: sdk.NewInt(1000),
// 		MinLiquidity:     sdk.NewInt(1000),
// 		Deadline:         deadline.Unix(),
// 		Sender:           val.Address.String(),
// 	}
// 	res = s.network.BlockSendMsgs(s.T(), msgAddLiquidity)
// 	s.Require().Equal(uint32(0), res.Code, res.Log)
// }
