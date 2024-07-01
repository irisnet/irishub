package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/gogoproto/proto"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/e2e"
	servicecli "mods.irisnet.org/modules/service/client/cli"
	"mods.irisnet.org/modules/service/types"
	servicetypes "mods.irisnet.org/modules/service/types"
	"mods.irisnet.org/simapp"
)

// QueryTestSuite is a suite of end-to-end tests for the service module
type QueryTestSuite struct {
	e2e.TestSuite
}

// SetupSuite sets up test suite
func (s *QueryTestSuite) SetupSuite() {
	s.SetModifyConfigFn(func(cfg *network.Config) {
		var serviceGenesisState servicetypes.GenesisState
		cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[servicetypes.ModuleName], &serviceGenesisState)

		serviceGenesisState.Params.ArbitrationTimeLimit = time.Duration(time.Second)
		serviceGenesisState.Params.ComplaintRetrospect = time.Duration(time.Second)
		cfg.GenesisState[servicetypes.ModuleName] = cfg.Codec.MustMarshalJSON(&serviceGenesisState)
		cfg.NumValidators = 1
	})
	s.TestSuite.SetupSuite()
}

// TestQueryCmd tests all query command in the service module
func (s *QueryTestSuite) TestQueryCmd() {
	val := s.Network.Validators[0]
	clientCtx := val.ClientCtx
	expectedCode := uint32(0)
	// ---------------------------------------------------------------------------

	serviceName := "test-service"
	serviceDesc := "test-description"
	serviceAuthorDesc := "test-author-description"
	serviceTags := "tags3,tags4"
	serviceSchemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	serviceDenom := sdk.DefaultBondDenom
	baseURL := val.APIAddress

	serviceDeposit := fmt.Sprintf("50000%s", serviceDenom)
	servicePrices := fmt.Sprintf(`{"price": "50%s"}`, serviceDenom)
	qos := uint64(3)
	options := "{}"

	author := val.Address
	provider := author

	consumerInfo, _, _ := val.ClientCtx.Keyring.NewMnemonic(
		"NewValidator",
		keyring.English,
		sdk.FullFundraiserPath,
		keyring.DefaultBIP39Passphrase,
		hd.Secp256k1,
	)

	consumer, err := consumerInfo.GetAddress()
	s.Require().NoError(err)

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
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult := DefineServiceExec(
		s.T(),
		s.Network,
		clientCtx,
		author.String(),
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	//------test GetCmdQueryServiceDefinition()-------------
	url := fmt.Sprintf("%s/irismod/service/definitions/%s", baseURL, serviceName)
	resp, err := testutil.GetRequest(url)
	respType := proto.Message(&servicetypes.QueryDefinitionResponse{})
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
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = BindServiceExec(
		s.T(),
		s.Network,
		clientCtx,
		provider.String(),
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	//------test GetCmdQueryServiceBinding()-------------
	url = fmt.Sprintf("%s/irismod/service/bindings/%s/%s", baseURL, serviceName, provider.String())
	resp, err = testutil.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryBindingResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	serviceBindingResp := respType.(*servicetypes.QueryBindingResponse)
	s.Require().Equal(serviceName, serviceBindingResp.ServiceBinding.ServiceName)
	s.Require().Equal(provider.String(), serviceBindingResp.ServiceBinding.Provider)

	//------test GetCmdQueryServiceBindings()-------------
	url = fmt.Sprintf("%s/irismod/service/bindings/%s", baseURL, serviceName)
	resp, err = testutil.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryBindingsResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	serviceBindings := respType.(*servicetypes.QueryBindingsResponse)
	s.Require().Len(serviceBindings.ServiceBindings, 1)

	//------test GetCmdDisableServiceBinding()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = DisableServiceExec(
		s.T(),
		s.Network,
		clientCtx,
		serviceName,
		provider.String(),
		provider.String(),
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	url = fmt.Sprintf("%s/irismod/service/bindings/%s/%s", baseURL, serviceName, provider.String())
	resp, err = testutil.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryBindingResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	serviceBindingResp = respType.(*servicetypes.QueryBindingResponse)
	s.Require().False(serviceBindingResp.ServiceBinding.Available)

	//------test GetCmdRefundServiceDeposit()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = RefundDepositExec(
		s.T(),
		s.Network,
		clientCtx,
		serviceName,
		provider.String(),
		provider.String(),
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	url = fmt.Sprintf("%s/irismod/service/bindings/%s/%s", baseURL, serviceName, provider.String())
	resp, err = testutil.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryBindingResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	serviceBindingResp = respType.(*servicetypes.QueryBindingResponse)
	s.Require().True(serviceBindingResp.ServiceBinding.Deposit.IsZero())

	//------test GetCmdEnableServiceBinding()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", servicecli.FlagDeposit, serviceDeposit),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = EnableServiceExec(
		s.T(),
		s.Network,
		clientCtx,
		serviceName,
		provider.String(),
		provider.String(),
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	url = fmt.Sprintf("%s/irismod/service/bindings/%s/%s", baseURL, serviceName, provider.String())
	resp, err = testutil.GetRequest(url)
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
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = simapp.MsgSendExec(s.T(), s.Network, clientCtx, provider, consumer, amount, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	//------test GetCmdCallService()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", servicecli.FlagServiceName, serviceName),
		fmt.Sprintf("--%s=%s", servicecli.FlagProviders, provider),
		fmt.Sprintf("--%s=%s", servicecli.FlagServiceFeeCap, reqServiceFee),
		fmt.Sprintf("--%s=%s", servicecli.FlagData, reqInput),
		fmt.Sprintf("--%s=%d", servicecli.FlagTimeout, timeout),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = CallServiceExec(
		s.T(),
		s.Network,
		clientCtx,
		consumer.String(),
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	requestContextId := s.Network.GetAttribute(
		servicetypes.EventTypeCreateContext,
		servicetypes.AttributeKeyRequestContextID,
		txResult.Events,
	)
	requestHeight := txResult.Height

	blockResult, err := val.RPCClient.BlockResults(context.Background(), &requestHeight)
	s.Require().NoError(err)
	var compactRequest servicetypes.CompactRequest
	for _, event := range blockResult.EndBlockEvents {
		if event.Type == servicetypes.EventTypeNewBatchRequest {
			var found bool
			var requests []servicetypes.CompactRequest
			var requestsBz []byte
			for _, attribute := range event.Attributes {
				if string(attribute.Key) == types.AttributeKeyRequests {
					requestsBz = []byte(attribute.Value)
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
	resp, err = testutil.GetRequest(url)
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
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = RespondServiceExec(
		s.T(),
		s.Network,
		clientCtx,
		provider.String(),
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	//------test GetCmdQueryEarnedFees()-------------
	url = fmt.Sprintf("%s/irismod/service/fees/%s", baseURL, provider.String())
	resp, err = testutil.GetRequest(url)
	respType = proto.Message(&servicetypes.QueryEarnedFeesResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	earnedFees := respType.(*servicetypes.QueryEarnedFeesResponse).Fees
	s.Require().Equal(expectedEarnedFees, earnedFees.String())
}
