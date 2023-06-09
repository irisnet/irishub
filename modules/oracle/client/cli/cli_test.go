package cli_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	oraclecli "github.com/irisnet/irismod/modules/oracle/client/cli"
	oracletestutil "github.com/irisnet/irismod/modules/oracle/client/testutil"
	servicecli "github.com/irisnet/irismod/modules/service/client/cli"
	servicetestutil "github.com/irisnet/irismod/modules/service/client/testutil"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	"github.com/irisnet/irismod/simapp"
)

type IntegrationTestSuite struct {
	suite.Suite

	network simapp.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = simapp.SetupNetwork(s.T())
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
	expectedCode := uint32(0)

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
	respResult := `{"code":200,"message":""}`
	respOutput := `{"header":{},"body":{"price":"2"}}`
	providers := provider
	timeout := 2
	newTimeout := qos
	serviceFeeCap := fmt.Sprintf("50%s", serviceDenom)
	threshold := 1
	frequency := 12

	//------Define && Bind Service-------------
	args := []string{
		fmt.Sprintf("--%s=%s", servicecli.FlagName, serviceName),
		fmt.Sprintf("--%s=%s", servicecli.FlagDescription, serviceDesc),
		fmt.Sprintf("--%s=%s", servicecli.FlagTags, serviceTags),
		fmt.Sprintf("--%s=%s", servicecli.FlagAuthorDescription, serviceAuthorDesc),
		fmt.Sprintf("--%s=%s", servicecli.FlagSchemas, serviceSchemas),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}

	txResult := servicetestutil.DefineServiceExec(s.T(), s.network, clientCtx, author.String(), args...)
	s.Require().Equal(expectedCode, txResult.Code)

	//------test GetCmdBindService()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", servicecli.FlagServiceName, serviceName),
		fmt.Sprintf("--%s=%s", servicecli.FlagDeposit, serviceDeposit),
		fmt.Sprintf("--%s=%s", servicecli.FlagPricing, servicePrices),
		fmt.Sprintf("--%s=%d", servicecli.FlagQoS, qos),
		fmt.Sprintf("--%s=%s", servicecli.FlagOptions, options),
		fmt.Sprintf("--%s=%s", servicecli.FlagProvider, provider),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}

	txResult = servicetestutil.BindServiceExec(s.T(), s.network, clientCtx, provider.String(), args...)
	s.Require().Equal(expectedCode, txResult.Code)

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
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}

	txResult = oracletestutil.CreateFeedExec(s.T(), s.network, clientCtx, creator.String(), args...)
	s.Require().Equal(expectedCode, txResult.Code)

	// ------test GetCmdQueryFeed()-------------

	feedContext := oracletestutil.QueryFeedExec(s.T(), s.network, clientCtx, feedName)
	s.Require().Equal(feedName, feedContext.Feed.FeedName)
	s.Require().Equal(servicetypes.PAUSED, feedContext.State)

	// ------test GetCmdQueryFeeds()-------------
	feedsResp := oracletestutil.QueryFeedsExec(s.T(), s.network, clientCtx)
	s.Require().Len(feedsResp.Feeds, 1)
	s.Require().Equal(*feedContext, feedsResp.Feeds[0])

	// ------test GetCmdStartFeed()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}

	txResult = oracletestutil.StartFeedExec(s.T(), s.network, clientCtx, creator.String(), feedName, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	feedContext = oracletestutil.QueryFeedExec(s.T(), s.network, clientCtx, feedName)
	s.Require().Equal(servicetypes.RUNNING, feedContext.State)

	// ------test GetCmdPauseFeed()-------------
	txResult = oracletestutil.PauseFeedExec(s.T(), s.network, clientCtx, creator.String(), feedName, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	feedContext = oracletestutil.QueryFeedExec(s.T(), s.network, clientCtx, feedName)
	s.Require().Equal(servicetypes.PAUSED, feedContext.State)

	// ------test GetCmdEditFeed()-------------
	args = []string{
		fmt.Sprintf("--%s=%d", oraclecli.FlagTimeout, newTimeout),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}

	txResult = oracletestutil.EditFeedExec(s.T(), s.network, clientCtx, creator.String(), feedName, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	feedContext = oracletestutil.QueryFeedExec(s.T(), s.network, clientCtx, feedName)
	s.Require().Equal(newTimeout, feedContext.Timeout)
	s.Require().Equal(servicetypes.PAUSED, feedContext.State)

	// ------test GetCmdQueryFeedValue()-------------
	feedValueResp := oracletestutil.QueryFeedValueExec(s.T(), s.network, clientCtx, feedName)
	s.Require().Len(feedValueResp.FeedValues, 0)

	// ------restart Feed-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}

	txResult = oracletestutil.StartFeedExec(s.T(), s.network, clientCtx, creator.String(), feedName, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	feedContext = oracletestutil.QueryFeedExec(s.T(), s.network, clientCtx, feedName)
	s.Require().Equal(servicetypes.RUNNING, feedContext.State)

	// ------get request-------------
	requestHeight := txResult.Height

	blockResult, err := val.RPCClient.BlockResults(context.Background(), &requestHeight)
	s.Require().NoError(err)
	var requestId string
	for _, event := range blockResult.EndBlockEvents {
		if event.Type == servicetypes.EventTypeNewBatchRequestProvider {
			var found bool
			var requestIds []string
			var requestsBz []byte
			for _, attribute := range event.Attributes {
				if string(attribute.Key) == servicetypes.AttributeKeyRequests {
					requestsBz = []byte(attribute.GetValue())
					found = true
				}
			}
			s.Require().True(found)
			if found {
				err := json.Unmarshal(requestsBz, &requestIds)
				s.Require().NoError(err)
			}
			s.Require().Len(requestIds, 1)
			requestId = requestIds[0]
		}
	}
	s.Require().NotNil(requestId)

	//------respond service-------------
	args = []string{
		fmt.Sprintf("--%s=%s", servicecli.FlagRequestID, requestId),
		fmt.Sprintf("--%s=%s", servicecli.FlagResult, respResult),
		fmt.Sprintf("--%s=%s", servicecli.FlagData, respOutput),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}

	txResult = servicetestutil.RespondServiceExec(s.T(), s.network, clientCtx, provider.String(), args...)
	s.Require().Equal(expectedCode, txResult.Code)

	// ------get feedValue-------------
	feedValueResp = oracletestutil.QueryFeedValueExec(s.T(), s.network, clientCtx, feedName)
	s.Require().Len(feedValueResp.FeedValues, 1)
	s.Require().Equal((strconv.FormatFloat(2, 'f', 8, 64)), feedValueResp.FeedValues[0].Data)
}
