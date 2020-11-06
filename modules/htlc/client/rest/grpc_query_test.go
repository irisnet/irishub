package rest_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto"
	"testing"

	htlctypes "github.com/irisnet/irismod/modules/htlc/types"

	htlccli "github.com/irisnet/irismod/modules/htlc/client/cli"
	htlctestutil "github.com/irisnet/irismod/modules/htlc/client/testutil"
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

func (s *IntegrationTestSuite) TestHtlc() {
	val := s.network.Validators[0]

	//------test GetCmdCreateHTLC()-------------
	baseURL:=val.APIAddress
	from := val.Address
	to := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))
	amount := "1000" + sdk.DefaultBondDenom
	receiverOnOtherChain := "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826"
	hashLock := "e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561"
	timeLock := uint64(50)
	timestamp := uint64(1580000000)
	stateOpen := "HTLC_STATE_OPEN"

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
	url :=fmt.Sprintf("%s/irismod/htlc/htlcs/%s", baseURL,hashLock)
	resp, err := rest.GetRequest(url)
	respType = proto.Message(&htlctypes.QueryHTLCResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(resp, respType))
	htlcResp := respType.(*htlctypes.QueryHTLCResponse)
	s.Require().Equal(amount, htlcResp.Htlc.Amount.String())
	s.Require().Equal(from.String(), htlcResp.Htlc.Sender)
	s.Require().Equal(to.String(), htlcResp.Htlc.To)
	s.Require().Equal(receiverOnOtherChain, htlcResp.Htlc.ReceiverOnOtherChain)
	s.Require().Equal(timestamp, htlcResp.Htlc.Timestamp)
	s.Require().Equal(stateOpen, htlcResp.Htlc.State.String())

}
