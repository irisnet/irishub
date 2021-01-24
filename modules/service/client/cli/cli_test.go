package cli_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	servicecli "github.com/irisnet/irismod/modules/service/client/cli"
	servicetestutil "github.com/irisnet/irismod/modules/service/client/testutil"
	"github.com/irisnet/irismod/modules/service/types"
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

	var serviceGenesisState servicetypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[servicetypes.ModuleName], &serviceGenesisState)

	serviceGenesisState.Params.ArbitrationTimeLimit = time.Duration(time.Second)
	serviceGenesisState.Params.ComplaintRetrospect = time.Duration(time.Second)
	cfg.GenesisState[servicetypes.ModuleName] = cfg.Codec.MustMarshalJSON(&serviceGenesisState)

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

func (s *IntegrationTestSuite) TestService() {
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
	qos := uint64(3)
	options := "{}"

	author := val.Address
	provider := author

	consumerInfo, _, _ := val.ClientCtx.Keyring.NewMnemonic("NewValidator", keyring.English, sdk.FullFundraiserPath, hd.Secp256k1)
	consumer := sdk.AccAddress(consumerInfo.GetPubKey().Address())

	reqServiceFee := fmt.Sprintf("50%s", serviceDenom)
	reqInput := `{"header":{},"body":{}}`
	respResult := `{"code":200,"message":""}`
	respOutput := `{"header":{},"body":{}}`
	timeout := qos

	expectedEarnedFees := fmt.Sprintf("48%s", serviceDenom)

	withdrawalAddress := sdk.AccAddress(crypto.AddressHash([]byte("withdrawalAddress")))

	//------test GetCmdDefineService()-------------
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

	//------test GetCmdQueryServiceDefinition()-------------
	respType = proto.Message(&servicetypes.ServiceDefinition{})
	bz, err = servicetestutil.QueryServiceDefinitionExec(val.ClientCtx, serviceName)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	serviceDefinition := respType.(*servicetypes.ServiceDefinition)
	s.Require().Equal(serviceName, serviceDefinition.Name)

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

	//------test GetCmdQueryServiceBinding()-------------
	respType = proto.Message(&servicetypes.ServiceBinding{})
	bz, err = servicetestutil.QueryServiceBindingExec(val.ClientCtx, serviceName, provider.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	serviceBinding := respType.(*servicetypes.ServiceBinding)
	s.Require().Equal(serviceName, serviceBinding.ServiceName)
	s.Require().Equal(provider.String(), serviceBinding.Provider)

	//------test GetCmdQueryServiceBindings()-------------
	respType = proto.Message(&servicetypes.QueryBindingsResponse{})
	bz, err = servicetestutil.QueryServiceBindingsExec(val.ClientCtx, serviceName)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	serviceBindings := respType.(*servicetypes.QueryBindingsResponse)
	s.Require().Len(serviceBindings.ServiceBindings, 1)

	//------test GetCmdDisableServiceBinding()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	expectedCode = uint32(0)
	bz, err = servicetestutil.DisableServiceExec(clientCtx, serviceName, provider.String(), provider.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&servicetypes.ServiceBinding{})
	bz, err = servicetestutil.QueryServiceBindingExec(val.ClientCtx, serviceName, provider.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	serviceBinding = respType.(*servicetypes.ServiceBinding)
	s.Require().False(serviceBinding.Available)

	//------test GetCmdRefundServiceDeposit()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	expectedCode = uint32(0)
	bz, err = servicetestutil.RefundDepositExec(clientCtx, serviceName, provider.String(), provider.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&servicetypes.ServiceBinding{})
	bz, err = servicetestutil.QueryServiceBindingExec(val.ClientCtx, serviceName, provider.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	serviceBinding = respType.(*servicetypes.ServiceBinding)
	s.Require().True(serviceBinding.Deposit.IsZero())

	//------test GetCmdEnableServiceBinding()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", servicecli.FlagDeposit, serviceDeposit),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	expectedCode = uint32(0)
	bz, err = servicetestutil.EnableServiceExec(clientCtx, serviceName, provider.String(), provider.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&servicetypes.ServiceBinding{})
	bz, err = servicetestutil.QueryServiceBindingExec(val.ClientCtx, serviceName, provider.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	serviceBinding = respType.(*servicetypes.ServiceBinding)
	s.Require().Equal(serviceDeposit, serviceBinding.Deposit.String())

	//------send token to consumer------------------------
	amount := sdk.NewCoins(
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(50000000)),
	)
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})
	expectedCode = uint32(0)
	bz, err = banktestutil.MsgSendExec(clientCtx, provider, consumer, amount, args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	//------test GetCmdCallService()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", servicecli.FlagServiceName, serviceName),
		fmt.Sprintf("--%s=%s", servicecli.FlagProviders, provider),
		fmt.Sprintf("--%s=%s", servicecli.FlagServiceFeeCap, reqServiceFee),
		fmt.Sprintf("--%s=%s", servicecli.FlagData, reqInput),
		fmt.Sprintf("--%s=%d", servicecli.FlagTimeout, timeout),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	expectedCode = uint32(0)
	bz, err = servicetestutil.CallServiceExec(clientCtx, consumer.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
	requestContextId := gjson.Get(txResp.RawLog, "0.events.0.attributes.0.value").String()
	requestHeight := txResp.Height

	blockResult, err := clientCtx.Client.BlockResults(context.Background(), &requestHeight)
	s.Require().NoError(err)
	var compactRequest servicetypes.CompactRequest
	for _, event := range blockResult.EndBlockEvents {
		if event.Type == servicetypes.EventTypeNewBatchRequest {
			var found bool
			var requests []servicetypes.CompactRequest
			var requestsBz []byte
			for _, attribute := range event.Attributes {
				if string(attribute.Key) == types.AttributeKeyRequests {
					requestsBz = attribute.GetValue()
				}
				if string(attribute.Key) == types.AttributeKeyRequestContextID &&
					string(attribute.GetValue()) == requestContextId {
					found = true
				}
			}
			s.Require().True(found)
			if found {
				err := json.Unmarshal(requestsBz, &requests)
				s.Require().NoError(err)
			}
			s.Require().Len(requests, 1)
			compactRequest = requests[0]
		}
	}
	s.Require().Equal(requestContextId, compactRequest.RequestContextId)

	//------test GetCmdQueryServiceRequests()-------------
	respType = proto.Message(&servicetypes.QueryRequestsResponse{})
	bz, err = servicetestutil.QueryServiceRequestsExec(val.ClientCtx, serviceName, provider.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	requests := respType.(*servicetypes.QueryRequestsResponse).Requests
	s.Require().Len(requests, 1)
	s.Require().Equal(requestContextId, requests[0].RequestContextId)

	//------test GetCmdQueryServiceRequests()-------------
	respType = proto.Message(&servicetypes.QueryRequestsResponse{})
	bz, err = servicetestutil.QueryServiceRequestsByReqCtx(val.ClientCtx, requests[0].RequestContextId, fmt.Sprint(requests[0].RequestContextBatchCounter))
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	requests = respType.(*servicetypes.QueryRequestsResponse).Requests
	s.Require().Len(requests, 1)
	s.Require().Equal(requestContextId, requests[0].RequestContextId)

	//------test GetCmdRespondService()-------------
	request := requests[0]
	args = []string{
		fmt.Sprintf("--%s=%s", servicecli.FlagRequestID, request.Id),
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

	//------test GetCmdQueryEarnedFees()-------------
	respType = proto.Message(&servicetypes.QueryEarnedFeesResponse{})
	bz, err = servicetestutil.QueryEarnedFeesExec(val.ClientCtx, provider.String())
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	earnedFees := respType.(*servicetypes.QueryEarnedFeesResponse).Fees
	s.Require().Equal(expectedEarnedFees, earnedFees.String())

	//------GetCmdSetWithdrawAddr()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	expectedCode = uint32(0)
	bz, err = servicetestutil.SetWithdrawAddrExec(clientCtx, withdrawalAddress.String(), provider.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	//------GetCmdWithdrawEarnedFees()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType = proto.Message(&sdk.TxResponse{})
	expectedCode = uint32(0)
	bz, err = servicetestutil.WithdrawEarnedFeesExec(clientCtx, provider.String(), provider.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&banktypes.QueryAllBalancesResponse{})
	bz, err = banktestutil.QueryBalancesExec(val.ClientCtx, withdrawalAddress)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	withdrawalFees := respType.(*banktypes.QueryAllBalancesResponse).Balances
	s.Require().Equal(expectedEarnedFees, withdrawalFees.String())

	//------GetCmdQueryRequestContext()-------------
	contextId := request.RequestContextId
	respType = proto.Message(&servicetypes.RequestContext{})
	bz, err = servicetestutil.QueryRequestContextExec(val.ClientCtx, contextId)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	contextResp := respType.(*servicetypes.RequestContext)
	s.Require().False(contextResp.Empty())

	//------GetCmdQueryServiceRequest()-------------
	requestId := request.Id
	respType = proto.Message(&servicetypes.Request{})
	bz, err = servicetestutil.QueryServiceRequestExec(val.ClientCtx, requestId)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	requestResp := respType.(*servicetypes.Request)
	s.Require().False(requestResp.Empty())
	s.Require().Equal(requestId, requestResp.Id)

	//------GetCmdQueryServiceResponse()-------------
	respType = proto.Message(&servicetypes.Response{})
	bz, err = servicetestutil.QueryServiceResponseExec(val.ClientCtx, requestId)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	responseResp := respType.(*servicetypes.Response)
	s.Require().False(responseResp.Empty())

}
