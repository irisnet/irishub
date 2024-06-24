package cli_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/suite"

// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	coinswaptypes "mods.irisnet.org/modules/coinswap/types"
// 	tokentypes "mods.irisnet.org/modules/token/types/v1"
// 	"mods.irisnet.org/simapp"
// 	farmcli "mods.irisnet.org/modules/farm/client/cli"
// 	"mods.irisnet.org/modules/farm/client/testutil"
// 	farmtypes "mods.irisnet.org/modules/farm/types"
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

// func (s *IntegrationTestSuite) TestFarm() {
// 	val := s.network.Validators[0]
// 	clientCtx := val.ClientCtx

// 	s.Init()

// 	// ---------------------------------------------------------------------------

// 	creator := val.Address
// 	description := "iris-atom farm pool"
// 	startHeight := s.LatestHeight() + 2
// 	rewardPerBlock := sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10)))
// 	totalReward := sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(1000)))
// 	editable := true
// 	lptDenom := "lpt-1"

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
// 		fmt.Sprintf("--%s=%s", farmcli.FlagLPTokenDenom, lptDenom),
// 		fmt.Sprintf("--%s=%s", farmcli.FlagTotalReward, totalReward),
// 		fmt.Sprintf("--%s=%v", farmcli.FlagEditable, editable),
// 	}

// 	args = append(args, globalFlags...)
// 	txResult := testutil.CreateFarmPoolExec(
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
// 	expectedContents := &farmtypes.FarmPoolEntry{
// 		Id:              poolId,
// 		Creator:         creator.String(),
// 		Description:     description,
// 		StartHeight:     startHeight,
// 		EndHeight:       startHeight + 100,
// 		Editable:        editable,
// 		Expired:         false,
// 		TotalLptLocked:  sdk.NewCoin(lptDenom, sdk.ZeroInt()),
// 		TotalReward:     totalReward,
// 		RemainingReward: totalReward,
// 		RewardPerBlock:  rewardPerBlock,
// 	}

// 	respType := testutil.QueryFarmPoolExec(s.T(), s.network, val.ClientCtx, poolId)
// 	s.Require().EqualValues(expectedContents, respType.Pool)

// 	reward := sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(1000)))
// 	args = []string{
// 		fmt.Sprintf("--%s=%v", farmcli.FlagAdditionalReward, reward.String()),
// 	}
// 	args = append(args, globalFlags...)
// 	txResult = testutil.AppendRewardExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		creator.String(),
// 		poolId,
// 		args...,
// 	)
// 	s.Require().EqualValues(txResult.Code, 0, txResult.Log)

// 	lpToken := sdk.NewCoin(lptDenom, sdk.NewInt(100))
// 	txResult = testutil.StakeExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		creator.String(),
// 		poolId,
// 		lpToken.String(),
// 		globalFlags...,
// 	)
// 	s.Require().EqualValues(txResult.Code, 0, txResult.Log)
// 	beginHeight := txResult.Height

// 	unstakeLPToken := sdk.NewCoin(lptDenom, sdk.NewInt(50))
// 	txResult = testutil.UnstakeExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		creator.String(),
// 		poolId,
// 		unstakeLPToken.String(),
// 		globalFlags...,
// 	)
// 	s.Require().EqualValues(txResult.Code, 0, txResult.Log)
// 	endHeight := txResult.Height

// 	rewardGot := s.network.GetAttribute(
// 		farmtypes.EventTypeUnstake,
// 		farmtypes.AttributeValueReward,
// 		txResult.Events,
// 	)
// 	expectedReward := rewardPerBlock.MulInt(sdk.NewInt(endHeight - beginHeight))
// 	s.Require().Equal(expectedReward.String(), rewardGot)

// 	txResult = testutil.HarvestExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		creator.String(),
// 		poolId,
// 		globalFlags...,
// 	)
// 	s.Require().EqualValues(txResult.Code, 0, txResult.Log)
// 	endHeight1 := txResult.Height

// 	rewardGot = s.network.GetAttribute(
// 		farmtypes.EventTypeHarvest,
// 		farmtypes.AttributeValueReward,
// 		txResult.Events,
// 	)
// 	expectedReward = rewardPerBlock.MulInt(sdk.NewInt(endHeight1 - endHeight))
// 	s.Require().Equal(expectedReward.String(), rewardGot)

// 	queryFarmerArgs := []string{
// 		fmt.Sprintf("--%s=%s", farmcli.FlagFarmPool, poolId),
// 	}

// 	leftlpToken := lpToken.Sub(unstakeLPToken)
// 	response := testutil.QueryFarmerExec(
// 		s.T(),
// 		s.network,
// 		val.ClientCtx, creator.String(), queryFarmerArgs...)
// 	s.Require().EqualValues(leftlpToken, response.List[0].Locked)

// 	txResult = testutil.DestroyExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		creator.String(),
// 		poolId,
// 		globalFlags...,
// 	)
// 	s.Require().EqualValues(txResult.Code, 0, txResult.Log)
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
