package rest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"

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
	baseURL := val.APIAddress

	serviceDeposit := fmt.Sprintf("50000%s", serviceDenom)
	servicePrices := fmt.Sprintf(`{"price": "50%s"}`, serviceDenom)
	qos := uint64(3)
	options := "{}"

	author := val.Address
	provider := author

	consumerInfo, _, _ := val.ClientCtx.Keyring.NewMnemonic("NewValidator", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	consumer := sdk.AccAddress(consumerInfo.GetPubKey().Address())

	reqServiceFee := fmt.Sprintf("50%s", serviceDenom)
	reqInput := `{"header":{},"body":{}}`
	respResult := `{"code":200,"message":""}`
	respOutput := `{"header":{},"body":{}}`
	timeout := qos

	expectedEarnedFees := fmt.Sprintf("48%s", serviceDenom)

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
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	//------test GetCmdQueryServiceDefinition()-------------
	url := fmt.Sprintf("%s/irismod/service/definitions/%s", baseURL, serviceName)
	resp, err := rest.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryDefinitionResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	serviceDefinitionResp := respType.(*servicetypes.QueryDefinitionResponse)
	s.Require().Equal(serviceName, serviceDefinitionResp.ServiceDefinition.Name)

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

	//------test GetCmdQueryServiceBinding()-------------
	url = fmt.Sprintf("%s/irismod/service/bindings/%s/%s", baseURL, serviceName, provider.String())
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryBindingResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	serviceBindingResp := respType.(*servicetypes.QueryBindingResponse)
	s.Require().Equal(serviceName, serviceBindingResp.ServiceBinding.ServiceName)
	s.Require().Equal(provider.String(), serviceBindingResp.ServiceBinding.Provider)

	//------test GetCmdQueryServiceBindings()-------------
	url = fmt.Sprintf("%s/irismod/service/bindings/%s", baseURL, serviceName)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryBindingsResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
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
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	url = fmt.Sprintf("%s/irismod/service/bindings/%s/%s", baseURL, serviceName, provider.String())
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryBindingResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	serviceBindingResp = respType.(*servicetypes.QueryBindingResponse)
	s.Require().False(serviceBindingResp.ServiceBinding.Available)

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
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	url = fmt.Sprintf("%s/irismod/service/bindings/%s/%s", baseURL, serviceName, provider.String())
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryBindingResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	serviceBindingResp = respType.(*servicetypes.QueryBindingResponse)
	s.Require().True(serviceBindingResp.ServiceBinding.Deposit.IsZero())

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
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	url = fmt.Sprintf("%s/irismod/service/bindings/%s/%s", baseURL, serviceName, provider.String())
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryBindingResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	serviceBindingResp = respType.(*servicetypes.QueryBindingResponse)
	s.Require().Equal(serviceDeposit, serviceBindingResp.ServiceBinding.Deposit.String())

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
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
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
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
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
	url = fmt.Sprintf("%s/irismod/service/requests/%s/%s", baseURL, serviceName, provider.String())
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryRequestsResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	requests := respType.(*servicetypes.QueryRequestsResponse).Requests
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
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	//------test GetCmdQueryEarnedFees()-------------
	url = fmt.Sprintf("%s/irismod/service/fees/%s", baseURL, provider.String())
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryEarnedFeesResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	earnedFees := respType.(*servicetypes.QueryEarnedFeesResponse).Fees
	s.Require().Equal(expectedEarnedFees, earnedFees.String())
}
