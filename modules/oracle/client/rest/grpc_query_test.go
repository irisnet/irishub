package rest

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"

	oraclecli "github.com/irisnet/irismod/modules/oracle/client/cli"
	oracletestutil "github.com/irisnet/irismod/modules/oracle/client/testutil"
	oracletypes "github.com/irisnet/irismod/modules/oracle/types"

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
	feedName := "test-feed"
	aggregateFunc := "avg"
	valueJsonPath := "high"
	latestHistory := 10
	description := "description"
	serviceName := "test-service"
	input := `{"header":{},"body":{}}`
	providers := []string{
		from.String(),
	}
	timeout := 2
	serviceFeeCap := "1iris"
	threshold := 1
	frequency := 12
	creator := from.String()

	args := []string{
		fmt.Sprintf("--%s=%s", oraclecli.FlagFeedName, feedName),
		fmt.Sprintf("--%s=%s", oraclecli.FlagAggregateFunc, aggregateFunc),
		fmt.Sprintf("--%s=%s", oraclecli.FlagValueJsonPath, valueJsonPath),
		fmt.Sprintf("--%s=%d", oraclecli.FlagLatestHistory, latestHistory),
		fmt.Sprintf("--%s=%s", oraclecli.FlagDescription, description),
		fmt.Sprintf("--%s=%s", oraclecli.FlagServiceFeeCap, serviceFeeCap),
		fmt.Sprintf("--%s=%s", oraclecli.FlagServiceName, serviceName),
		fmt.Sprintf("--%s=%s", oraclecli.FlagInput, input),
		fmt.Sprintf("--%s=%t", oraclecli.FlagProviders, providers),
		fmt.Sprintf("--%s=%s", oraclecli.FlagTimeout, timeout),
		fmt.Sprintf("--%s=%s", oraclecli.FlagThreshold, threshold),
		fmt.Sprintf("--%s=%s", oraclecli.FlagFrequency, frequency),
		fmt.Sprintf("--%s=%s", oraclecli.FlagCreator, creator),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := oracletestutil.CreateFeedExec(clientCtx, from.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	// ---------------------------------------------------------------------------

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = oracletestutil.StartFeedExec(clientCtx, from.String(), feedName, args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	// ---------------------------------------------------------------------------

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = oracletestutil.PauseFeedExec(clientCtx, from.String(), feedName, args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	// ---------------------------------------------------------------------------
	args = []string{
		fmt.Sprintf("--%s=%d", oraclecli.FlagLatestHistory, latestHistory),
		fmt.Sprintf("--%s=%t", oraclecli.FlagProviders, providers),
		fmt.Sprintf("--%s=%s", oraclecli.FlagServiceFeeCap, serviceFeeCap),
		fmt.Sprintf("--%s=%s", oraclecli.FlagTimeout, timeout),
		fmt.Sprintf("--%s=%s", oraclecli.FlagFrequency, frequency),
		fmt.Sprintf("--%s=%s", oraclecli.FlagThreshold, threshold),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})

	bz, err = oracletestutil.EditFeedExec(clientCtx, from.String(), feedName, args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&oracletypes.FeedContext{})
	bz, err = oracletestutil.QueryFeedExec(clientCtx, feedName)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	//feedResp := respType.(*oracletypes.FeedContext)
	s.Require().NoError(err)

	respType = proto.Message(&oracletypes.QueryFeedsResponse{})
	bz, err = oracletestutil.QueryFeedsExec(clientCtx)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	//feedResp := respType.(*oracletypes.QueryFeedsResponse)
	s.Require().NoError(err)

	respType = proto.Message(&oracletypes.QueryFeedValueResponse{})
	bz, err = oracletestutil.QueryFeedValueExec(clientCtx, feedName)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	//feedResp := respType.(*oracletypes.QueryFeedValueResponse)
	s.Require().NoError(err)
}
