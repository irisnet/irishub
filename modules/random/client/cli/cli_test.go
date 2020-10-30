package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/tidwall/gjson"

	randomcli "github.com/irisnet/irismod/modules/random/client/cli"
	randomtestutil "github.com/irisnet/irismod/modules/random/client/testutil"
	randomtypes "github.com/irisnet/irismod/modules/random/types"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/stretchr/testify/suite"

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

func (s *IntegrationTestSuite) TestOracle() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	// ---------------------------------------------------------------------------

	from := val.Address
	blockInterval := 10
	oracle := true
	serviceFeeCap := "1iris"

	args := []string{
		fmt.Sprintf("--%s=%s", randomcli.FlagServiceFeeCap, serviceFeeCap),
		fmt.Sprintf("--%s=%t", randomcli.FlagOracle, oracle),
		fmt.Sprintf("--%s=%d", randomcli.FlagBlockInterval, blockInterval),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := randomtestutil.RequestRandomExec(clientCtx, from.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
	requestID := gjson.Get(txResp.RawLog, "0.events.1.attributes.0.value").String()
	generateHeight := gjson.Get(txResp.RawLog, "0.events.1.attributes.1.value").String()
	println(txResp.RawLog)

	// ---------------------------------------------------------------------------

	respType = proto.Message(&randomtypes.Random{})
	bz, err = randomtestutil.QueryRandomExec(clientCtx, requestID)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	//randomResp := respType.(*randomtypes.Random)
	s.Require().NoError(err)

	respType = proto.Message(&randomtypes.QueryRandomRequestQueueResponse{})
	bz, err = randomtestutil.QueryRandomRequestQueueExec(clientCtx, generateHeight)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	//qrrResp := respType.(*randomtypes.QueryRandomRequestQueueResponse)
	s.Require().NoError(err)

}
