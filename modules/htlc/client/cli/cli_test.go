package cli_test

// import (
// 	"fmt"
// 	"math/rand"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/suite"

// 	tmbytes "github.com/cometbft/cometbft/libs/bytes"

// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/irisnet/irismod/simapp"
// 	htlccli "github.com/irisnet/irismod/htlc/client/cli"
// 	htlctestutil "github.com/irisnet/irismod/htlc/client/testutil"
// 	htlctypes "github.com/irisnet/irismod/htlc/types"
// )

// const (
// 	BNB_DENOM   = "htltbnb"
// 	DEPUTY_ADDR = "cosmos1kznrznww4pd6gx0zwrpthjk68fdmqypjpkj5hp"
// )

// var (
// 	Deputy               sdk.AccAddress
// 	MinTimeLock          uint64 = 50
// 	MaxTimeLock          uint64 = 60
// 	ReceiverOnOtherChain        = "ReceiverOnOtherChain"
// 	SenderOnOtherChain          = "SenderOnOtherChain"
// )

// const DeputyArmor = `-----BEGIN TENDERMINT PRIVATE KEY-----
// salt: C3586B75587D2824187D2CDA22B6AFB6
// type: secp256k1
// kdf: bcrypt

// 1+15OrCKgjnwym1zO3cjo/SGe3PPqAYChQ5wMHjdUbTZM7mWsH3/ueL6swgjzI3b
// DDzEQAPXBQflzNW6wbne9IfT651zCSm+j1MWaGk=
// =wEHs
// -----END TENDERMINT PRIVATE KEY-----`

// type IntegrationTestSuite struct {
// 	suite.Suite

// 	network simapp.Network
// }

// func c(denom string, amount int64) sdk.Coin {
// 	return sdk.NewInt64Coin(denom, amount)
// }

// func cs(coins ...sdk.Coin) sdk.Coins {
// 	return sdk.NewCoins(coins...)
// }

// func ts(minOffset int) uint64 {
// 	return uint64(time.Now().Add(time.Duration(minOffset) * time.Minute).Unix())
// }

// func (s *IntegrationTestSuite) SetupSuite() {
// 	s.T().Log("setting up integration test suite")

// 	cfg := simapp.NewConfig()
// 	cfg.NumValidators = 4

// 	Deputy, _ = sdk.AccAddressFromBech32(DEPUTY_ADDR)
// 	cfg.GenesisState[htlctypes.ModuleName] = cfg.Codec.MustMarshalJSON(NewHTLTGenesis(Deputy))
// 	s.network = simapp.SetupNetworkWithConfig(s.T(), cfg)
// }

// func (s *IntegrationTestSuite) TearDownSuite() {
// 	s.T().Log("tearing down integration test suite")
// 	s.network.Cleanup()
// }

// func TestIntegrationTestSuite(t *testing.T) {
// 	suite.Run(t, new(IntegrationTestSuite))
// }

// func (s *IntegrationTestSuite) TestHTLC() {
// 	// ---------------------------------------------------------------
// 	ctx := s.network.Validators[0].ClientCtx
// 	err := ctx.Keyring.ImportPrivKey("deputy", DeputyArmor, "1234567890")
// 	s.Require().NoError(err)

// 	args := []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	_ = simapp.MsgSendExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		s.network.Validators[0].Address,
// 		Deputy,
// 		cs(c(sdk.DefaultBondDenom, 50000000)),
// 		args...,
// 	)

// 	// ---------------------------------------------------------------

// 	type htlcArgs struct {
// 		sender             sdk.AccAddress
// 		receiver           sdk.AccAddress
// 		receiverOtherChain string
// 		senderOtherChain   string
// 		amount             sdk.Coins
// 		secret             tmbytes.HexBytes
// 		timestamp          uint64
// 		timeLock           uint64
// 		transfer           bool
// 		direction          htlctypes.SwapDirection
// 	}
// 	testCases := []struct {
// 		name string
// 		args htlcArgs
// 		pass bool
// 	}{{
// 		"valid htlc",
// 		htlcArgs{
// 			sender:             s.network.Validators[0].Address,
// 			receiver:           s.network.Validators[1].Address,
// 			receiverOtherChain: ReceiverOnOtherChain,
// 			senderOtherChain:   SenderOnOtherChain,
// 			amount:             cs(c(sdk.DefaultBondDenom, 1000)),
// 			secret:             GenerateRandomSecret(),
// 			timestamp:          uint64(1580000000),
// 			timeLock:           uint64(50),
// 			transfer:           false,
// 			direction:          htlctypes.None,
// 		},
// 		true,
// 	}, {
// 		"valid incoming htlt",
// 		htlcArgs{
// 			sender:             Deputy,
// 			receiver:           s.network.Validators[0].Address,
// 			receiverOtherChain: ReceiverOnOtherChain,
// 			senderOtherChain:   SenderOnOtherChain,
// 			amount:             cs(c(BNB_DENOM, 10000)),
// 			secret:             GenerateRandomSecret(),
// 			timestamp:          ts(0),
// 			timeLock:           MinTimeLock,
// 			transfer:           true,
// 			direction:          htlctypes.Incoming,
// 		},
// 		true,
// 	}, {
// 		"valid outgoing htlt",
// 		htlcArgs{
// 			sender:             s.network.Validators[0].Address,
// 			receiver:           Deputy,
// 			receiverOtherChain: ReceiverOnOtherChain,
// 			senderOtherChain:   SenderOnOtherChain,
// 			amount:             cs(c(BNB_DENOM, 5000)),
// 			secret:             GenerateRandomSecret(),
// 			timestamp:          ts(0),
// 			timeLock:           MinTimeLock,
// 			transfer:           true,
// 			direction:          htlctypes.Outgoing,
// 		},
// 		true,
// 	}}

// 	// ---------------------------------------------------------------
// 	// HTLC
// 	// ---------------------------------------------------------------

// 	args = []string{
// 		fmt.Sprintf("--%s=%s", htlccli.FlagTo, testCases[0].args.receiver),
// 		fmt.Sprintf("--%s=%s", htlccli.FlagAmount, testCases[0].args.amount),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			htlccli.FlagReceiverOnOtherChain,
// 			testCases[0].args.receiverOtherChain,
// 		),
// 		fmt.Sprintf("--%s=%s", htlccli.FlagSenderOnOtherChain, testCases[0].args.senderOtherChain),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			htlccli.FlagHashLock,
// 			tmbytes.HexBytes(htlctypes.GetHashLock(testCases[0].args.secret, testCases[0].args.timestamp)).
// 				String(),
// 		),
// 		fmt.Sprintf("--%s=%d", htlccli.FlagTimeLock, testCases[0].args.timeLock),
// 		fmt.Sprintf("--%s=%d", htlccli.FlagTimestamp, testCases[0].args.timestamp),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			htlccli.FlagTransfer,
// 			strconv.FormatBool(testCases[0].args.transfer),
// 		),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult := htlctestutil.CreateHTLCExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		testCases[0].args.sender.String(),
// 		args...,
// 	)

// 	// ---------------------------------------------------------------

// 	expectedhtlc := htlctypes.HTLC{
// 		Id: htlctypes.GetID(testCases[0].args.sender, testCases[0].args.receiver, testCases[0].args.amount, htlctypes.GetHashLock(testCases[0].args.secret, testCases[0].args.timestamp)).
// 			String(),
// 		Sender:               testCases[0].args.sender.String(),
// 		To:                   testCases[0].args.receiver.String(),
// 		ReceiverOnOtherChain: ReceiverOnOtherChain,
// 		SenderOnOtherChain:   SenderOnOtherChain,
// 		Amount:               testCases[0].args.amount,
// 		Secret:               "",
// 		HashLock: tmbytes.HexBytes(htlctypes.GetHashLock(testCases[0].args.secret, testCases[0].args.timestamp)).
// 			String(),
// 		Timestamp:        testCases[0].args.timestamp,
// 		ExpirationHeight: uint64(txResult.Height) + testCases[0].args.timeLock,
// 		State:            htlctypes.Open,
// 		ClosedBlock:      0,
// 		Transfer:         testCases[0].args.transfer,
// 		Direction:        testCases[0].args.direction,
// 	}
// 	respType := htlctestutil.QueryHTLCExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		expectedhtlc.Id,
// 	)
// 	s.Require().Equal(expectedhtlc.String(), respType.String())

// 	// ---------------------------------------------------------------

// 	args = []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = htlctestutil.ClaimHTLCExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		testCases[0].args.sender.String(),
// 		expectedhtlc.Id,
// 		testCases[0].args.secret.String(),
// 		args...,
// 	)

// 	respType = htlctestutil.QueryHTLCExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		expectedhtlc.Id,
// 	)
// 	s.Require().Equal(htlctypes.Completed.String(), respType.State.String())

// 	balance := simapp.QueryBalanceExec(
// 		s.T(),
// 		s.network,
// 		ctx, testCases[0].args.receiver.String(),
// 		sdk.DefaultBondDenom,
// 	)
// 	s.Require().Equal("400001000stake", balance.String())

// 	// ---------------------------------------------------------------
// 	// HTLT INCOMING
// 	// ---------------------------------------------------------------

// 	args = []string{
// 		fmt.Sprintf("--%s=%s", htlccli.FlagTo, testCases[1].args.receiver),
// 		fmt.Sprintf("--%s=%s", htlccli.FlagAmount, testCases[1].args.amount),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			htlccli.FlagReceiverOnOtherChain,
// 			testCases[1].args.receiverOtherChain,
// 		),
// 		fmt.Sprintf("--%s=%s", htlccli.FlagSenderOnOtherChain, testCases[1].args.senderOtherChain),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			htlccli.FlagHashLock,
// 			tmbytes.HexBytes(htlctypes.GetHashLock(testCases[1].args.secret, testCases[1].args.timestamp)).
// 				String(),
// 		),
// 		fmt.Sprintf("--%s=%d", htlccli.FlagTimeLock, testCases[1].args.timeLock),
// 		fmt.Sprintf("--%s=%d", htlccli.FlagTimestamp, testCases[1].args.timestamp),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			htlccli.FlagTransfer,
// 			strconv.FormatBool(testCases[1].args.transfer),
// 		),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = htlctestutil.CreateHTLCExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		testCases[1].args.sender.String(),
// 		args...,
// 	)

// 	// ---------------------------------------------------------------

// 	expectedhtlt := htlctypes.HTLC{
// 		Id: htlctypes.GetID(testCases[1].args.sender, testCases[1].args.receiver, testCases[1].args.amount, htlctypes.GetHashLock(testCases[1].args.secret, testCases[1].args.timestamp)).
// 			String(),
// 		Sender:               testCases[1].args.sender.String(),
// 		To:                   testCases[1].args.receiver.String(),
// 		ReceiverOnOtherChain: ReceiverOnOtherChain,
// 		SenderOnOtherChain:   SenderOnOtherChain,
// 		Amount:               testCases[1].args.amount,
// 		Secret:               "",
// 		HashLock: tmbytes.HexBytes(htlctypes.GetHashLock(testCases[1].args.secret, testCases[1].args.timestamp)).
// 			String(),
// 		Timestamp:        testCases[1].args.timestamp,
// 		ExpirationHeight: uint64(txResult.Height) + testCases[1].args.timeLock,
// 		State:            htlctypes.Open,
// 		ClosedBlock:      0,
// 		Transfer:         testCases[1].args.transfer,
// 		Direction:        testCases[1].args.direction,
// 	}
// 	respType = htlctestutil.QueryHTLCExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		expectedhtlt.Id,
// 	)
// 	s.Require().Equal(expectedhtlt.String(), respType.String())

// 	// ---------------------------------------------------------------

// 	args = []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = htlctestutil.ClaimHTLCExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		testCases[1].args.sender.String(),
// 		expectedhtlt.Id,
// 		testCases[1].args.secret.String(),
// 		args...,
// 	)

// 	respType = htlctestutil.QueryHTLCExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		expectedhtlc.Id,
// 	)
// 	s.Require().Equal(htlctypes.Completed.String(), respType.State.String())

// 	// ---------------------------------------------------------------
// 	// HTLT OUTGOING
// 	// ---------------------------------------------------------------

// 	args = []string{
// 		fmt.Sprintf("--%s=%s", htlccli.FlagTo, testCases[2].args.receiver),
// 		fmt.Sprintf("--%s=%s", htlccli.FlagAmount, testCases[2].args.amount),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			htlccli.FlagReceiverOnOtherChain,
// 			testCases[2].args.receiverOtherChain,
// 		),
// 		fmt.Sprintf("--%s=%s", htlccli.FlagSenderOnOtherChain, testCases[2].args.senderOtherChain),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			htlccli.FlagHashLock,
// 			tmbytes.HexBytes(htlctypes.GetHashLock(testCases[2].args.secret, testCases[2].args.timestamp)).
// 				String(),
// 		),
// 		fmt.Sprintf("--%s=%d", htlccli.FlagTimeLock, testCases[2].args.timeLock),
// 		fmt.Sprintf("--%s=%d", htlccli.FlagTimestamp, testCases[2].args.timestamp),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			htlccli.FlagTransfer,
// 			strconv.FormatBool(testCases[2].args.transfer),
// 		),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = htlctestutil.CreateHTLCExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		testCases[2].args.sender.String(),
// 		args...,
// 	)

// 	// ---------------------------------------------------------------

// 	expectedhtlt = htlctypes.HTLC{
// 		Id: htlctypes.GetID(testCases[2].args.sender, testCases[2].args.receiver, testCases[2].args.amount, htlctypes.GetHashLock(testCases[2].args.secret, testCases[2].args.timestamp)).
// 			String(),
// 		Sender:               testCases[2].args.sender.String(),
// 		To:                   testCases[2].args.receiver.String(),
// 		ReceiverOnOtherChain: ReceiverOnOtherChain,
// 		SenderOnOtherChain:   SenderOnOtherChain,
// 		Amount:               testCases[2].args.amount,
// 		Secret:               "",
// 		HashLock: tmbytes.HexBytes(htlctypes.GetHashLock(testCases[2].args.secret, testCases[2].args.timestamp)).
// 			String(),
// 		Timestamp:        testCases[2].args.timestamp,
// 		ExpirationHeight: uint64(txResult.Height) + testCases[2].args.timeLock,
// 		State:            htlctypes.Open,
// 		ClosedBlock:      0,
// 		Transfer:         testCases[2].args.transfer,
// 		Direction:        testCases[2].args.direction,
// 	}

// 	respType = htlctestutil.QueryHTLCExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		expectedhtlc.Id,
// 	)
// 	s.Require().Equal(htlctypes.Completed.String(), respType.State.String())

// 	// ---------------------------------------------------------------

// 	args = []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = htlctestutil.ClaimHTLCExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		testCases[2].args.sender.String(),
// 		expectedhtlt.Id,
// 		testCases[2].args.secret.String(),
// 		args...,
// 	)

// 	respType = htlctestutil.QueryHTLCExec(
// 		s.T(),
// 		s.network,
// 		ctx,
// 		expectedhtlc.Id,
// 	)
// 	s.Require().Equal(htlctypes.Completed.String(), respType.State.String())

// 	// ---------------------------------------------------------------
// }

// func NewHTLTGenesis(deputyAddress sdk.AccAddress) *htlctypes.GenesisState {
// 	return &htlctypes.GenesisState{
// 		Params: htlctypes.Params{
// 			AssetParams: []htlctypes.AssetParam{
// 				{
// 					Denom: "htltbnb",
// 					SupplyLimit: htlctypes.SupplyLimit{
// 						Limit:          sdk.NewInt(350000000000000),
// 						TimeLimited:    false,
// 						TimeBasedLimit: sdk.ZeroInt(),
// 						TimePeriod:     time.Hour,
// 					},
// 					Active:        true,
// 					DeputyAddress: deputyAddress.String(),
// 					FixedFee:      sdk.NewInt(1000),
// 					MinSwapAmount: sdk.OneInt(),
// 					MaxSwapAmount: sdk.NewInt(1000000000000),
// 					MinBlockLock:  MinTimeLock,
// 					MaxBlockLock:  MaxTimeLock,
// 				},
// 			},
// 		},
// 		Htlcs: []htlctypes.HTLC{},
// 		Supplies: []htlctypes.AssetSupply{
// 			htlctypes.NewAssetSupply(
// 				sdk.NewCoin("htltbnb", sdk.ZeroInt()),
// 				sdk.NewCoin("htltbnb", sdk.ZeroInt()),
// 				sdk.NewCoin("htltbnb", sdk.ZeroInt()),
// 				sdk.NewCoin("htltbnb", sdk.ZeroInt()),
// 				time.Duration(0),
// 			),
// 		},
// 		PreviousBlockTime: htlctypes.DefaultPreviousBlockTime,
// 	}
// }

// func GenerateRandomSecret() tmbytes.HexBytes {
// 	bytes := make([]byte, 32)
// 	if _, err := rand.Read(bytes); err != nil {
// 		panic(err.Error())
// 	}
// 	return bytes
// }
