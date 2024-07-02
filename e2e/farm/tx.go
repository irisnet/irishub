package farm

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/e2e"
	coinswaptypes "mods.irisnet.org/modules/coinswap/types"
	farmcli "mods.irisnet.org/modules/farm/client/cli"
	farmtypes "mods.irisnet.org/modules/farm/types"
	tokentypes "mods.irisnet.org/modules/token/types/v1"
)

// TxTestSuite is a suite of end-to-end tests for the nft module
type TxTestSuite struct {
	e2e.TestSuite
}

// TestTxCmd tests all tx command in the nft module
func (s *TxTestSuite) TestTxCmd() {
	val := s.Network.Validators[0]
	clientCtx := val.ClientCtx

	s.setup()

	// ---------------------------------------------------------------------------

	creator := val.Address
	description := "iris-atom farm pool"
	startHeight := s.latestHeight() + 2
	rewardPerBlock := sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10)))
	totalReward := sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(1000)))
	editable := true
	lptDenom := "lpt-1"

	globalFlags := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	args := []string{
		fmt.Sprintf("--%s=%s", farmcli.FlagDescription, description),
		fmt.Sprintf("--%s=%d", farmcli.FlagStartHeight, startHeight),
		fmt.Sprintf("--%s=%s", farmcli.FlagRewardPerBlock, rewardPerBlock),
		fmt.Sprintf("--%s=%s", farmcli.FlagLPTokenDenom, lptDenom),
		fmt.Sprintf("--%s=%s", farmcli.FlagTotalReward, totalReward),
		fmt.Sprintf("--%s=%v", farmcli.FlagEditable, editable),
	}

	args = append(args, globalFlags...)
	txResult := CreateFarmPoolExec(
		s.T(),
		s.Network,
		clientCtx,
		creator.String(),
		args...,
	)
	s.Require().EqualValues(txResult.Code, 0, txResult.Log)

	poolID := s.Network.GetAttribute(
		farmtypes.EventTypeCreatePool,
		farmtypes.AttributeValuePoolId,
		txResult.Events,
	)
	expectedContents := &farmtypes.FarmPoolEntry{
		Id:              poolID,
		Creator:         creator.String(),
		Description:     description,
		StartHeight:     startHeight,
		EndHeight:       startHeight + 100,
		Editable:        editable,
		Expired:         false,
		TotalLptLocked:  sdk.NewCoin(lptDenom, sdk.ZeroInt()),
		TotalReward:     totalReward,
		RemainingReward: totalReward,
		RewardPerBlock:  rewardPerBlock,
	}

	respType := QueryFarmPoolExec(s.T(), s.Network, val.ClientCtx, poolID)
	s.Require().EqualValues(expectedContents, respType.Pool)

	reward := sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(1000)))
	args = []string{
		fmt.Sprintf("--%s=%v", farmcli.FlagAdditionalReward, reward.String()),
	}
	args = append(args, globalFlags...)
	txResult = AppendRewardExec(
		s.T(),
		s.Network,
		clientCtx,
		creator.String(),
		poolID,
		args...,
	)
	s.Require().EqualValues(txResult.Code, 0, txResult.Log)

	lpToken := sdk.NewCoin(lptDenom, sdk.NewInt(100))
	txResult = StakeExec(
		s.T(),
		s.Network,
		clientCtx,
		creator.String(),
		poolID,
		lpToken.String(),
		globalFlags...,
	)
	s.Require().EqualValues(txResult.Code, 0, txResult.Log)
	beginHeight := txResult.Height

	unstakeLPToken := sdk.NewCoin(lptDenom, sdk.NewInt(50))
	txResult = UnstakeExec(
		s.T(),
		s.Network,
		clientCtx,
		creator.String(),
		poolID,
		unstakeLPToken.String(),
		globalFlags...,
	)
	s.Require().EqualValues(txResult.Code, 0, txResult.Log)
	endHeight := txResult.Height

	rewardGot := s.Network.GetAttribute(
		farmtypes.EventTypeUnstake,
		farmtypes.AttributeValueReward,
		txResult.Events,
	)
	expectedReward := rewardPerBlock.MulInt(sdk.NewInt(endHeight - beginHeight))
	s.Require().Equal(expectedReward.String(), rewardGot)

	txResult = HarvestExec(
		s.T(),
		s.Network,
		clientCtx,
		creator.String(),
		poolID,
		globalFlags...,
	)
	s.Require().EqualValues(txResult.Code, 0, txResult.Log)
	endHeight1 := txResult.Height

	rewardGot = s.Network.GetAttribute(
		farmtypes.EventTypeHarvest,
		farmtypes.AttributeValueReward,
		txResult.Events,
	)
	expectedReward = rewardPerBlock.MulInt(sdk.NewInt(endHeight1 - endHeight))
	s.Require().Equal(expectedReward.String(), rewardGot)

	queryFarmerArgs := []string{
		fmt.Sprintf("--%s=%s", farmcli.FlagFarmPool, poolID),
	}

	leftlpToken := lpToken.Sub(unstakeLPToken)
	response := QueryFarmerExec(
		s.T(),
		s.Network,
		val.ClientCtx, creator.String(), queryFarmerArgs...)
	s.Require().EqualValues(leftlpToken, response.List[0].Locked)

	txResult = DestroyExec(
		s.T(),
		s.Network,
		clientCtx,
		creator.String(),
		poolID,
		globalFlags...,
	)
	s.Require().EqualValues(txResult.Code, 0, txResult.Log)
}

func (s *TxTestSuite) latestHeight() int64 {
	height, err := s.Network.LatestHeight()
	s.Require().NoError(err)
	return height
}

func (s *TxTestSuite) setup() {
	val := s.Network.Validators[0]
	clientCtx := val.ClientCtx

	from := val.Address
	symbol := "kitty"
	name := "Kitty Token"
	minUnit := "kitty"
	scale := uint32(0)
	initialSupply := uint64(100000000)
	maxSupply := uint64(200000000)
	mintable := true

	// issue token
	msgIssueToken := &tokentypes.MsgIssueToken{
		Symbol:        symbol,
		Name:          name,
		Scale:         scale,
		MinUnit:       minUnit,
		InitialSupply: initialSupply,
		MaxSupply:     maxSupply,
		Mintable:      mintable,
		Owner:         from.String(),
	}
	res := s.Network.BlockSendMsgs(s.T(), msgIssueToken)
	s.Require().Equal(uint32(0), res.Code, res.Log)

	// add liquidity
	status, err := clientCtx.Client.Status(context.Background())
	s.Require().NoError(err)
	deadline := status.SyncInfo.LatestBlockTime.Add(time.Minute)

	msgAddLiquidity := &coinswaptypes.MsgAddLiquidity{
		MaxToken:         sdk.NewCoin(symbol, sdk.NewInt(1000)),
		ExactStandardAmt: sdk.NewInt(1000),
		MinLiquidity:     sdk.NewInt(1000),
		Deadline:         deadline.Unix(),
		Sender:           val.Address.String(),
	}
	res = s.Network.BlockSendMsgs(s.T(), msgAddLiquidity)
	s.Require().Equal(uint32(0), res.Code, res.Log)
}
