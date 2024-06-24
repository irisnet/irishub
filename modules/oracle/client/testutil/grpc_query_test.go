package testutil_test

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/cosmos/gogoproto/proto"
// 	"github.com/stretchr/testify/suite"

// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	"github.com/cosmos/cosmos-sdk/testutil"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	servicecli "mods.irisnet.org/modules/service/client/cli"
// 	servicetestutil "mods.irisnet.org/modules/service/client/testutil"
// 	servicetypes "mods.irisnet.org/modules/service/types"
// 	"mods.irisnet.org/simapp"
// 	oraclecli "mods.irisnet.org/modules/oracle/client/cli"
// 	oracletestutil "mods.irisnet.org/modules/oracle/client/testutil"
// 	oracletypes "mods.irisnet.org/modules/oracle/types"
// )

// type IntegrationTestSuite struct {
// 	suite.Suite

// 	network simapp.Network
// }

// func (s *IntegrationTestSuite) SetupSuite() {
// 	s.T().Log("setting up integration test suite")

// 	s.network = simapp.SetupNetwork(s.T())
// }

// func (s *IntegrationTestSuite) TearDownSuite() {
// 	s.T().Log("tearing down integration test suite")
// 	s.network.Cleanup()
// }

// func TestIntegrationTestSuite(t *testing.T) {
// 	suite.Run(t, new(IntegrationTestSuite))
// }

// func (s *IntegrationTestSuite) TestOracle() {
// 	val := s.network.Validators[0]
// 	clientCtx := val.ClientCtx
// 	expectedCode := uint32(0)

// 	// ---------------------------------------------------------------------------
// 	serviceName := "test-service"
// 	serviceDesc := "test-description"
// 	serviceAuthorDesc := "test-author-description"
// 	serviceTags := "tags1,tags2"
// 	serviceSchemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
// 	serviceDenom := sdk.DefaultBondDenom

// 	serviceDeposit := fmt.Sprintf("50000%s", serviceDenom)
// 	servicePrices := fmt.Sprintf(`{"price": "50%s"}`, serviceDenom)
// 	qos := int64(3)
// 	options := "{}"

// 	author := val.Address
// 	provider := author
// 	creator := author

// 	feedName := "test-feed"
// 	aggregateFunc := "avg"
// 	valueJsonPath := "price"
// 	latestHistory := 10
// 	description := "description"
// 	input := `{"header":{},"body":{}}`
// 	providers := provider
// 	timeout := 2
// 	serviceFeeCap := fmt.Sprintf("50%s", serviceDenom)
// 	threshold := 1
// 	frequency := 12
// 	baseURL := val.APIAddress

// 	//------Define && Bind Service-------------
// 	args := []string{
// 		fmt.Sprintf("--%s=%s", servicecli.FlagName, serviceName),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagDescription, serviceDesc),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagTags, serviceTags),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagAuthorDescription, serviceAuthorDesc),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagSchemas, serviceSchemas),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult := servicetestutil.DefineServiceExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		author.String(),
// 		args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	//------test GetCmdBindService()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=%s", servicecli.FlagServiceName, serviceName),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagDeposit, serviceDeposit),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagPricing, servicePrices),
// 		fmt.Sprintf("--%s=%d", servicecli.FlagQoS, qos),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagOptions, options),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagProvider, provider),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = servicetestutil.BindServiceExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		provider.String(),
// 		args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	//------test GetCmdCreateFeed()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=%s", oraclecli.FlagFeedName, feedName),
// 		fmt.Sprintf("--%s=%s", oraclecli.FlagAggregateFunc, aggregateFunc),
// 		fmt.Sprintf("--%s=%s", oraclecli.FlagValueJsonPath, valueJsonPath),
// 		fmt.Sprintf("--%s=%d", oraclecli.FlagLatestHistory, latestHistory),
// 		fmt.Sprintf("--%s=%s", oraclecli.FlagDescription, description),
// 		fmt.Sprintf("--%s=%s", oraclecli.FlagServiceFeeCap, serviceFeeCap),
// 		fmt.Sprintf("--%s=%s", oraclecli.FlagServiceName, serviceName),
// 		fmt.Sprintf("--%s=%s", oraclecli.FlagInput, input),
// 		fmt.Sprintf("--%s=%s", oraclecli.FlagProviders, providers),
// 		fmt.Sprintf("--%s=%d", oraclecli.FlagTimeout, timeout),
// 		fmt.Sprintf("--%s=%d", oraclecli.FlagThreshold, threshold),
// 		fmt.Sprintf("--%s=%d", oraclecli.FlagFrequency, frequency),
// 		fmt.Sprintf("--%s=%s", oraclecli.FlagCreator, creator),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = oracletestutil.CreateFeedExec(s.T(), s.network, clientCtx, creator.String(), args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	// ------test GetCmdQueryFeed()-------------
// 	url := fmt.Sprintf("%s/irismod/oracle/feeds/%s", baseURL, feedName)
// 	resp, err := testutil.GetRequest(url)
// 	s.Require().NoError(err)
// 	respType := proto.Message(&oracletypes.QueryFeedResponse{})
// 	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
// 	feedResp := respType.(*oracletypes.QueryFeedResponse)
// 	s.Require().NoError(err)
// 	s.Require().Equal(feedName, feedResp.Feed.Feed.FeedName)
// 	s.Require().Equal(servicetypes.PAUSED, feedResp.Feed.State)

// 	// ------test GetCmdQueryFeeds()-------------
// 	url = fmt.Sprintf("%s/irismod/oracle/feeds", baseURL)
// 	resp, err = testutil.GetRequest(url)
// 	s.Require().NoError(err)
// 	respType = proto.Message(&oracletypes.QueryFeedsResponse{})
// 	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
// 	feedsResp := respType.(*oracletypes.QueryFeedsResponse)
// 	s.Require().NoError(err)
// 	s.Require().Len(feedsResp.Feeds, 1)
// 	s.Require().Equal(feedResp.Feed, feedsResp.Feeds[0])

// 	// ------test GetCmdQueryFeedValue()-------------
// 	url = fmt.Sprintf("%s/irismod/oracle/feeds/%s/values", baseURL, feedName)
// 	resp, err = testutil.GetRequest(url)
// 	respType = proto.Message(&oracletypes.QueryFeedValueResponse{})
// 	s.Require().NoError(err)
// 	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
// 	feedValueResp := respType.(*oracletypes.QueryFeedValueResponse)
// 	s.Require().NoError(err)
// 	s.Require().Len(feedValueResp.FeedValues, 0)
// }
