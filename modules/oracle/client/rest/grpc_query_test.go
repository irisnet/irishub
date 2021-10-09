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

	oraclecli "github.com/irisnet/irismod/modules/oracle/client/cli"
	oracletestutil "github.com/irisnet/irismod/modules/oracle/client/testutil"
	oracletypes "github.com/irisnet/irismod/modules/oracle/types"
	servicecli "github.com/irisnet/irismod/modules/service/client/cli"
	servicetestutil "github.com/irisnet/irismod/modules/service/client/testutil"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
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
	serviceName := "test-service"
	serviceDesc := "test-description"
	serviceAuthorDesc := "test-author-description"
	serviceTags := "tags1,tags2"
	serviceSchemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	serviceDenom := sdk.DefaultBondDenom

	serviceDeposit := fmt.Sprintf("50000%s", serviceDenom)
	servicePrices := fmt.Sprintf(`{"price": "50%s"}`, serviceDenom)
	qos := int64(3)
	options := "{}"

	author := val.Address
	provider := author
	creator := author

	feedName := "test-feed"
	aggregateFunc := "avg"
	valueJsonPath := "price"
	latestHistory := 10
	description := "description"
	input := `{"header":{},"body":{}}`
	providers := provider
	timeout := 2
	serviceFeeCap := fmt.Sprintf("50%s", serviceDenom)
	threshold := 1
	frequency := 12
	baseURL := val.APIAddress

	//------Define && Bind Service-------------
	args := []string{
		fmt.Sprintf("--%s=%s", servicecli.FlagName, serviceName),
		fmt.Sprintf("--%s=%s", servicecli.FlagDescription, serviceDesc),
		fmt.Sprintf("--%s=%s", servicecli.FlagTags, serviceTags),
		fmt.Sprintf("--%s=%s", servicecli.FlagAuthorDescription, serviceAuthorDesc),
		fmt.Sprintf("--%s=%s", servicecli.FlagSchemas, serviceSchemas),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)
	bz, err := servicetestutil.DefineServiceExec(clientCtx, author.String(), args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	//------test GetCmdBindService()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", servicecli.FlagServiceName, serviceName),
		fmt.Sprintf("--%s=%s", servicecli.FlagDeposit, serviceDeposit),
		fmt.Sprintf("--%s=%s", servicecli.FlagPricing, servicePrices),
		fmt.Sprintf("--%s=%d", servicecli.FlagQoS, qos),
		fmt.Sprintf("--%s=%s", servicecli.FlagOptions, options),
		fmt.Sprintf("--%s=%s", servicecli.FlagProvider, provider),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	expectedCode = uint32(0)
	bz, err = servicetestutil.BindServiceExec(clientCtx, provider.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	//------test GetCmdCreateFeed()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", oraclecli.FlagFeedName, feedName),
		fmt.Sprintf("--%s=%s", oraclecli.FlagAggregateFunc, aggregateFunc),
		fmt.Sprintf("--%s=%s", oraclecli.FlagValueJsonPath, valueJsonPath),
		fmt.Sprintf("--%s=%d", oraclecli.FlagLatestHistory, latestHistory),
		fmt.Sprintf("--%s=%s", oraclecli.FlagDescription, description),
		fmt.Sprintf("--%s=%s", oraclecli.FlagServiceFeeCap, serviceFeeCap),
		fmt.Sprintf("--%s=%s", oraclecli.FlagServiceName, serviceName),
		fmt.Sprintf("--%s=%s", oraclecli.FlagInput, input),
		fmt.Sprintf("--%s=%s", oraclecli.FlagProviders, providers),
		fmt.Sprintf("--%s=%d", oraclecli.FlagTimeout, timeout),
		fmt.Sprintf("--%s=%d", oraclecli.FlagThreshold, threshold),
		fmt.Sprintf("--%s=%d", oraclecli.FlagFrequency, frequency),
		fmt.Sprintf("--%s=%s", oraclecli.FlagCreator, creator),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})
	expectedCode = uint32(0)

	bz, err = oracletestutil.CreateFeedExec(clientCtx, creator.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	// ------test GetCmdQueryFeed()-------------
	url := fmt.Sprintf("%s/irismod/oracle/feeds/%s", baseURL, feedName)
	resp, err := rest.GetRequest(url)
	s.Require().NoError(err)
	respType = proto.Message(&oracletypes.QueryFeedResponse{})
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	feedResp := respType.(*oracletypes.QueryFeedResponse)
	s.Require().NoError(err)
	s.Require().Equal(feedName, feedResp.Feed.Feed.FeedName)
	s.Require().Equal(servicetypes.PAUSED, feedResp.Feed.State)

	// ------test GetCmdQueryFeeds()-------------
	url = fmt.Sprintf("%s/irismod/oracle/feeds", baseURL)
	resp, err = rest.GetRequest(url)
	s.Require().NoError(err)
	respType = proto.Message(&oracletypes.QueryFeedsResponse{})
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	feedsResp := respType.(*oracletypes.QueryFeedsResponse)
	s.Require().NoError(err)
	s.Require().Len(feedsResp.Feeds, 1)
	s.Require().Equal(feedResp.Feed, feedsResp.Feeds[0])

	// ------test GetCmdQueryFeedValue()-------------
	url = fmt.Sprintf("%s/irismod/oracle/feeds/%s/values", baseURL, feedName)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&oracletypes.QueryFeedValueResponse{})
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	feedValueResp := respType.(*oracletypes.QueryFeedValueResponse)
	s.Require().NoError(err)
	s.Require().Len(feedValueResp.FeedValues, 0)
}
