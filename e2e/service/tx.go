package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"mods.irisnet.org/e2e"
	servicecli "mods.irisnet.org/modules/service/client/cli"
	servicetypes "mods.irisnet.org/modules/service/types"
	"mods.irisnet.org/simapp"
)

// TxTestSuite is a suite of end-to-end tests for the service module
type TxTestSuite struct {
	e2e.TestSuite
}

// SetupSuite sets up test suite
func (s *TxTestSuite) SetupSuite() {
	s.SetModifyConfigFn(func(cfg *network.Config) {
		var serviceGenesisState servicetypes.GenesisState
		cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[servicetypes.ModuleName], &serviceGenesisState)

		serviceGenesisState.Params.ArbitrationTimeLimit = time.Second
		serviceGenesisState.Params.ComplaintRetrospect = time.Second
		cfg.GenesisState[servicetypes.ModuleName] = cfg.Codec.MustMarshalJSON(&serviceGenesisState)
		cfg.NumValidators = 1
	})
	s.TestSuite.SetupSuite()
}

// TestQueryCmd tests all query command in the service module
func (s *TxTestSuite) TestQueryCmd() {
	val := s.Network.Validators[0]
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
	pubKey, err := consumerInfo.GetPubKey()
	s.Require().NoError(err)
	consumer := sdk.AccAddress(pubKey.Address())

	reqServiceFee := fmt.Sprintf("50%s", serviceDenom)
	reqInput := `{"header":{},"body":{}}`
	respResult := `{"code":200,"message":""}`
	respOutput := `{"header":{},"body":{}}`
	timeout := qos

	expectedEarnedFees := fmt.Sprintf("48%s", serviceDenom)
	expectedTaxFees := fmt.Sprintf("2%s", serviceDenom)

	withdrawalAddress := sdk.AccAddress(crypto.AddressHash([]byte("withdrawalAddress")))

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
	serviceDefinition := QueryServiceDefinitionExec(
		s.T(),
		s.Network,
		clientCtx,
		serviceName,
	)
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
	serviceBinding := QueryServiceBindingExec(
		s.T(),
		s.Network,
		clientCtx,
		serviceName,
		provider.String(),
	)
	s.Require().Equal(serviceName, serviceBinding.ServiceName)
	s.Require().Equal(provider.String(), serviceBinding.Provider)

	//------test GetCmdQueryServiceBindings()-------------
	serviceBindings := QueryServiceBindingsExec(
		s.T(),
		s.Network,
		clientCtx,
		serviceName,
	)
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

	serviceBinding = QueryServiceBindingExec(
		s.T(),
		s.Network,
		clientCtx,
		serviceName,
		provider.String(),
	)
	s.Require().False(serviceBinding.Available)

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

	serviceBinding = QueryServiceBindingExec(
		s.T(),
		s.Network,
		clientCtx,
		serviceName,
		provider.String(),
	)
	s.Require().True(serviceBinding.Deposit.IsZero())

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

	serviceBinding = QueryServiceBindingExec(
		s.T(),
		s.Network,
		clientCtx,
		serviceName,
		provider.String(),
	)
	s.Require().Equal(serviceDeposit, serviceBinding.Deposit.String())

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
				if attribute.Key == servicetypes.AttributeKeyRequests {
					requestsBz = []byte(attribute.GetValue())
				}
				if attribute.Key == servicetypes.AttributeKeyRequestContextID &&
					attribute.GetValue() == requestContextId {
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
	queryRequestsResponse := QueryServiceRequestsExec(
		s.T(),
		s.Network,
		clientCtx,
		serviceName,
		provider.String(),
	)
	s.Require().Len(queryRequestsResponse.Requests, 1)
	s.Require().Equal(requestContextId, queryRequestsResponse.Requests[0].RequestContextId)

	//------test GetCmdQueryServiceRequests()-------------
	queryRequestsResponse = QueryServiceRequestsByReqCtx(
		s.T(),
		s.Network,
		clientCtx,
		queryRequestsResponse.Requests[0].RequestContextId,
		fmt.Sprint(queryRequestsResponse.Requests[0].RequestContextBatchCounter),
	)
	s.Require().Len(queryRequestsResponse.Requests, 1)
	s.Require().Equal(requestContextId, queryRequestsResponse.Requests[0].RequestContextId)

	//------test GetCmdRespondService()-------------
	request := queryRequestsResponse.Requests[0]
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
	queryEarnedFeesResponse := QueryEarnedFeesExec(
		s.T(),
		s.Network,
		clientCtx,
		provider.String(),
	)
	s.Require().Equal(expectedEarnedFees, queryEarnedFeesResponse.Fees.String())

	//------GetCmdSetWithdrawAddr()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = SetWithdrawAddrExec(
		s.T(),
		s.Network,
		clientCtx,
		withdrawalAddress.String(),
		provider.String(),
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	//------GetCmdWithdrawEarnedFees()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = WithdrawEarnedFeesExec(
		s.T(),
		s.Network,
		clientCtx,
		provider.String(),
		provider.String(),
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	withdrawalFees := simapp.QueryBalancesExec(
		s.T(),
		s.Network,
		clientCtx,
		withdrawalAddress.String(),
	)
	s.Require().Equal(expectedEarnedFees, withdrawalFees.String())

	//------check service tax-------------
	taxFees := simapp.QueryBalancesExec(
		s.T(),
		s.Network,
		clientCtx,
		authtypes.NewModuleAddress(servicetypes.FeeCollectorName).String(),
	)
	s.Require().Equal(expectedTaxFees, taxFees.String())

	//------GetCmdQueryRequestContext()-------------
	contextId := request.RequestContextId
	contextResp := QueryRequestContextExec(s.T(), s.Network, clientCtx, contextId)
	s.Require().False(contextResp.Empty())

	//------GetCmdQueryServiceRequest()-------------
	requestId := request.Id
	requestResp := QueryServiceRequestExec(s.T(), s.Network, clientCtx, requestId)
	s.Require().False(requestResp.Empty())
	s.Require().Equal(requestId, requestResp.Id)

	//------GetCmdQueryServiceResponse()-------------
	responseResp := QueryServiceResponseExec(s.T(), s.Network, clientCtx, requestId)
	s.Require().False(responseResp.Empty())
}
