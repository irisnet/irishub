package cli_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	htlccli "github.com/irisnet/irismod/modules/htlc/client/cli"
	htlctestutil "github.com/irisnet/irismod/modules/htlc/client/testutil"
	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
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

func (s *IntegrationTestSuite) TestHTLC() {
	val := s.network.Validators[0]

	//------test GetCmdCreateHTLC()-------------
	from := val.Address
	to := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))
	amount := "1000" + sdk.DefaultBondDenom
	receiverOnOtherChain := "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826"
	hashLock := "e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561"
	secretHex := "5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f"
	timeLock := uint64(50)
	timestamp := uint64(1580000000)
	stateOpen := "HTLC_STATE_OPEN"
	stateCompleted := "HTLC_STATE_COMPLETED"
	stateRefunded := "HTLC_STATE_REFUNDED"
	stateExpired := "HTLC_STATE_EXPIRED"

	args := []string{
		fmt.Sprintf("--%s=%s", htlccli.FlagTo, to),
		fmt.Sprintf("--%s=%s", htlccli.FlagAmount, amount),
		fmt.Sprintf("--%s=%s", htlccli.FlagReceiverOnOtherChain, receiverOnOtherChain),
		fmt.Sprintf("--%s=%s", htlccli.FlagHashLock, hashLock),
		fmt.Sprintf("--%s=%d", htlccli.FlagTimeLock, timeLock),
		fmt.Sprintf("--%s=%d", htlccli.FlagTimestamp, timestamp),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := htlctestutil.CreateHTLCExec(val.ClientCtx, from.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	//------test GetCmdQueryHTLC()-------------

	respType = proto.Message(&htlctypes.HTLC{})
	bz, err = htlctestutil.QueryHTLCExec(val.ClientCtx, hashLock)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	htlcItem := respType.(*htlctypes.HTLC)
	s.Require().Equal(amount, htlcItem.Amount.String())
	s.Require().Equal(from.String(), htlcItem.Sender)
	s.Require().Equal(to.String(), htlcItem.To)
	s.Require().Equal(receiverOnOtherChain, htlcItem.ReceiverOnOtherChain)
	s.Require().Equal(timestamp, htlcItem.Timestamp)
	s.Require().Equal(stateOpen, htlcItem.State.String())

	//------test GetCmdClaimHTLC()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = htlctestutil.ClaimHTLCExec(val.ClientCtx, from.String(), hashLock, secretHex, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&htlctypes.HTLC{})
	bz, err = htlctestutil.QueryHTLCExec(val.ClientCtx, hashLock)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	htlcItem = respType.(*htlctypes.HTLC)
	s.Require().Equal(stateCompleted, htlcItem.State.String())

	coinType := proto.Message(&sdk.Coin{})
	out, err := simapp.QueryBalanceExec(val.ClientCtx, to.String(), sdk.DefaultBondDenom)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), coinType))
	balance := coinType.(*sdk.Coin)
	s.Require().Equal(amount, balance.Amount.String()+sdk.DefaultBondDenom)

	//------test GetCmdClaimHTLC()-------------
	// testdata
	hashLock = "f054e34abd9ccc3cab12a5b797b8e9c053507f279e7e53fb3f9f44d178c94b20"
	timestamp = uint64(0)
	timeLock = uint64(50)

	// create an htlc
	args = []string{
		fmt.Sprintf("--%s=%s", htlccli.FlagTo, to),
		fmt.Sprintf("--%s=%s", htlccli.FlagAmount, amount),
		fmt.Sprintf("--%s=%s", htlccli.FlagReceiverOnOtherChain, receiverOnOtherChain),
		fmt.Sprintf("--%s=%s", htlccli.FlagHashLock, hashLock),
		fmt.Sprintf("--%s=%d", htlccli.FlagTimeLock, timeLock),
		fmt.Sprintf("--%s=%d", htlccli.FlagTimestamp, timestamp),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = htlctestutil.CreateHTLCExec(val.ClientCtx, from.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&htlctypes.HTLC{})
	bz, err = htlctestutil.QueryHTLCExec(val.ClientCtx, hashLock)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	htlcItem = respType.(*htlctypes.HTLC)
	s.Require().Equal(amount, htlcItem.Amount.String())
	s.Require().Equal(from.String(), htlcItem.Sender)
	s.Require().Equal(to.String(), htlcItem.To)
	s.Require().Equal(receiverOnOtherChain, htlcItem.ReceiverOnOtherChain)
	s.Require().Equal(timestamp, htlcItem.Timestamp)
	s.Require().Equal(stateOpen, htlcItem.State.String())

	// refund an htlc and expect failure
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})
	bz, err = htlctestutil.RefundHTLCExec(val.ClientCtx, from.String(), hashLock, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal("failed to execute message; message index: 0: F054E34ABD9CCC3CAB12A5B797B8E9C053507F279E7E53FB3F9F44D178C94B20: htlc not expired", txResp.RawLog)

	respType = proto.Message(&htlctypes.HTLC{})
	bz, err = htlctestutil.QueryHTLCExec(val.ClientCtx, hashLock)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	htlcItem = respType.(*htlctypes.HTLC)
	s.Require().Equal(stateOpen, htlcItem.State.String())

	// refund an htlc and expect success
	lastHeigth, err := s.network.LatestHeight()
	s.Require().NoError(err)
	_, _ = s.network.WaitForHeightWithTimeout(lastHeigth+int64(timeLock), 3600*time.Second)
	respType = proto.Message(&htlctypes.HTLC{})
	bz, err = htlctestutil.QueryHTLCExec(val.ClientCtx, hashLock)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	htlcItem = respType.(*htlctypes.HTLC)
	s.Equal(stateExpired, htlcItem.State.String())

	respType = proto.Message(&sdk.TxResponse{})
	bz, err = htlctestutil.RefundHTLCExec(val.ClientCtx, from.String(), hashLock, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&htlctypes.HTLC{})
	bz, err = htlctestutil.QueryHTLCExec(val.ClientCtx, hashLock)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	htlcItem = respType.(*htlctypes.HTLC)
	s.Equal(stateRefunded, htlcItem.State.String())
	// ---------------------------------------------------------------------------
}
