package cli_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)
	bz, err := servicetestutil.DefineServiceExec(clientCtx, author.String(), args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
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
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
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
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	// ------test GetCmdQueryFeed()-------------
	respType = proto.Message(&oracletypes.FeedContext{})
	bz, err = oracletestutil.QueryFeedExec(clientCtx, feedName)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	feedResp := respType.(*oracletypes.FeedContext)
	s.Require().NoError(err)
	s.Require().Equal(feedName, feedResp.Feed.FeedName)
	s.Require().Equal(servicetypes.PAUSED, feedResp.State)

	// ------test GetCmdQueryFeeds()-------------
	respType = proto.Message(&oracletypes.QueryFeedsResponse{})
	bz, err = oracletestutil.QueryFeedsExec(clientCtx)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	feedsResp := respType.(*oracletypes.QueryFeedsResponse)
	s.Require().NoError(err)
	s.Require().Len(feedsResp.Feeds, 1)
	s.Require().Equal(*feedResp, feedsResp.Feeds[0])

	// ------test GetCmdStartFeed()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})

	bz, err = oracletestutil.StartFeedExec(clientCtx, creator.String(), feedName, args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&oracletypes.FeedContext{})
	bz, err = oracletestutil.QueryFeedExec(clientCtx, feedName)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	feedResp = respType.(*oracletypes.FeedContext)
	s.Require().NoError(err)
	s.Require().Equal(servicetypes.RUNNING, feedResp.State)

	// ------test GetCmdPauseFeed()-------------
	respType = proto.Message(&sdk.TxResponse{})

	bz, err = oracletestutil.PauseFeedExec(clientCtx, creator.String(), feedName, args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&oracletypes.FeedContext{})
	bz, err = oracletestutil.QueryFeedExec(clientCtx, feedName)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	feedResp = respType.(*oracletypes.FeedContext)
	s.Require().NoError(err)
	s.Require().Equal(servicetypes.PAUSED, feedResp.State)

	// ------test GetCmdEditFeed()-------------
	args = []string{
		fmt.Sprintf("--%s=%d", oraclecli.FlagTimeout, newTimeout),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})

	bz, err = oracletestutil.EditFeedExec(clientCtx, creator.String(), feedName, args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&oracletypes.FeedContext{})
	bz, err = oracletestutil.QueryFeedExec(clientCtx, feedName)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	feedResp = respType.(*oracletypes.FeedContext)
	s.Require().NoError(err)
	s.Require().Equal(newTimeout, feedResp.Timeout)
	s.Require().Equal(servicetypes.PAUSED, feedResp.State)

	// ------test GetCmdQueryFeedValue()-------------
	respType = proto.Message(&oracletypes.QueryFeedValueResponse{})
	bz, err = oracletestutil.QueryFeedValueExec(clientCtx, feedName)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	feedValueResp := respType.(*oracletypes.QueryFeedValueResponse)
	s.Require().NoError(err)
	s.Require().Len(feedValueResp.FeedValues, 0)

	// ------restart Feed-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})

	bz, err = oracletestutil.StartFeedExec(clientCtx, creator.String(), feedName, args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&oracletypes.FeedContext{})
	bz, err = oracletestutil.QueryFeedExec(clientCtx, feedName)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	feedResp = respType.(*oracletypes.FeedContext)
	s.Require().NoError(err)
	s.Require().Equal(servicetypes.RUNNING, feedResp.State)

	// ------get request-------------
	requestHeight := txResp.Height

	blockResult, err := clientCtx.Client.BlockResults(context.Background(), &requestHeight)
	s.Require().NoError(err)
	var requestId string
	for _, event := range blockResult.EndBlockEvents {
		if event.Type == servicetypes.EventTypeNewBatchRequestProvider {
			var found bool
			var requestIds []string
			var requestsBz []byte
			for _, attribute := range event.Attributes {
				if string(attribute.Key) == servicetypes.AttributeKeyRequests {
					requestsBz = attribute.GetValue()
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
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	expectedCode = uint32(0)
	bz, err = servicetestutil.RespondServiceExec(clientCtx, provider.String(), args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	// ------get feedValue-------------
	respType = proto.Message(&oracletypes.QueryFeedValueResponse{})
	bz, err = oracletestutil.QueryFeedValueExec(clientCtx, feedName)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	feedValueResp = respType.(*oracletypes.QueryFeedValueResponse)
	s.Require().NoError(err)
	s.Require().Len(feedValueResp.FeedValues, 1)
	s.Require().Equal((strconv.FormatFloat(2, 'f', 8, 64)), feedValueResp.FeedValues[0].Data)
}
