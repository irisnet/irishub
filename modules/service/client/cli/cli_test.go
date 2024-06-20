package cli_test

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/suite"

// 	"github.com/cometbft/cometbft/crypto"

// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	"github.com/cosmos/cosmos-sdk/crypto/hd"
// 	"github.com/cosmos/cosmos-sdk/crypto/keyring"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

// 	"github.com/irisnet/irismod/simapp"
// 	servicecli "irismod.io/service/client/cli"
// 	servicetestutil "irismod.io/service/client/testutil"
// 	"irismod.io/service/types"
// 	servicetypes "irismod.io/service/types"
// )

// type IntegrationTestSuite struct {
// 	suite.Suite

// 	network simapp.Network
// }

// func (s *IntegrationTestSuite) SetupSuite() {
// 	s.T().Log("setting up integration test suite")

// 	cfg := simapp.NewConfig()
// 	cfg.NumValidators = 1

// 	var serviceGenesisState servicetypes.GenesisState
// 	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[servicetypes.ModuleName], &serviceGenesisState)

// 	serviceGenesisState.Params.ArbitrationTimeLimit = time.Duration(time.Second)
// 	serviceGenesisState.Params.ComplaintRetrospect = time.Duration(time.Second)
// 	cfg.GenesisState[servicetypes.ModuleName] = cfg.Codec.MustMarshalJSON(&serviceGenesisState)

// 	s.network = simapp.SetupNetworkWithConfig(s.T(), cfg)
// }

// func (s *IntegrationTestSuite) TearDownSuite() {
// 	s.T().Log("tearing down integration test suite")
// 	s.network.Cleanup()
// }

// func TestIntegrationTestSuite(t *testing.T) {
// 	suite.Run(t, new(IntegrationTestSuite))
// }

// func (s *IntegrationTestSuite) TestService() {
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
// 	qos := uint64(3)
// 	options := "{}"

// 	author := val.Address
// 	provider := author

// 	consumerInfo, _, _ := val.ClientCtx.Keyring.NewMnemonic(
// 		"NewValidator",
// 		keyring.English,
// 		sdk.FullFundraiserPath,
// 		keyring.DefaultBIP39Passphrase,
// 		hd.Secp256k1,
// 	)
// 	pubKey, err := consumerInfo.GetPubKey()
// 	s.Require().NoError(err)
// 	consumer := sdk.AccAddress(pubKey.Address())

// 	reqServiceFee := fmt.Sprintf("50%s", serviceDenom)
// 	reqInput := `{"header":{},"body":{}}`
// 	respResult := `{"code":200,"message":""}`
// 	respOutput := `{"header":{},"body":{}}`
// 	timeout := qos

// 	expectedEarnedFees := fmt.Sprintf("48%s", serviceDenom)
// 	expectedTaxFees := fmt.Sprintf("2%s", serviceDenom)

// 	withdrawalAddress := sdk.AccAddress(crypto.AddressHash([]byte("withdrawalAddress")))

// 	//------test GetCmdDefineService()-------------
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

// 	//------test GetCmdQueryServiceDefinition()-------------
// 	serviceDefinition := servicetestutil.QueryServiceDefinitionExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		serviceName,
// 	)
// 	s.Require().Equal(serviceName, serviceDefinition.Name)

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

// 	//------test GetCmdQueryServiceBinding()-------------
// 	serviceBinding := servicetestutil.QueryServiceBindingExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		serviceName,
// 		provider.String(),
// 	)
// 	s.Require().Equal(serviceName, serviceBinding.ServiceName)
// 	s.Require().Equal(provider.String(), serviceBinding.Provider)

// 	//------test GetCmdQueryServiceBindings()-------------
// 	serviceBindings := servicetestutil.QueryServiceBindingsExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		serviceName,
// 	)
// 	s.Require().Len(serviceBindings.ServiceBindings, 1)

// 	//------test GetCmdDisableServiceBinding()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = servicetestutil.DisableServiceExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		serviceName,
// 		provider.String(),
// 		provider.String(),
// 		args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	serviceBinding = servicetestutil.QueryServiceBindingExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		serviceName,
// 		provider.String(),
// 	)
// 	s.Require().False(serviceBinding.Available)

// 	//------test GetCmdRefundServiceDeposit()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = servicetestutil.RefundDepositExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		serviceName,
// 		provider.String(),
// 		provider.String(),
// 		args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	serviceBinding = servicetestutil.QueryServiceBindingExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		serviceName,
// 		provider.String(),
// 	)
// 	s.Require().True(serviceBinding.Deposit.IsZero())

// 	//------test GetCmdEnableServiceBinding()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=%s", servicecli.FlagDeposit, serviceDeposit),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = servicetestutil.EnableServiceExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		serviceName,
// 		provider.String(),
// 		provider.String(),
// 		args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	serviceBinding = servicetestutil.QueryServiceBindingExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		serviceName,
// 		provider.String(),
// 	)
// 	s.Require().Equal(serviceDeposit, serviceBinding.Deposit.String())

// 	//------send token to consumer------------------------
// 	amount := sdk.NewCoins(
// 		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(50000000)),
// 	)
// 	args = []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = simapp.MsgSendExec(s.T(), s.network, clientCtx, provider, consumer, amount, args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	//------test GetCmdCallService()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=%s", servicecli.FlagServiceName, serviceName),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagProviders, provider),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagServiceFeeCap, reqServiceFee),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagData, reqInput),
// 		fmt.Sprintf("--%s=%d", servicecli.FlagTimeout, timeout),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = servicetestutil.CallServiceExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		consumer.String(),
// 		args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	requestContextId := s.network.GetAttribute(
// 		servicetypes.EventTypeCreateContext,
// 		servicetypes.AttributeKeyRequestContextID,
// 		txResult.Events,
// 	)
// 	requestHeight := txResult.Height

// 	blockResult, err := val.RPCClient.BlockResults(context.Background(), &requestHeight)
// 	s.Require().NoError(err)
// 	var compactRequest servicetypes.CompactRequest
// 	for _, event := range blockResult.EndBlockEvents {
// 		if event.Type == servicetypes.EventTypeNewBatchRequest {
// 			var found bool
// 			var requests []servicetypes.CompactRequest
// 			var requestsBz []byte
// 			for _, attribute := range event.Attributes {
// 				if string(attribute.Key) == types.AttributeKeyRequests {
// 					requestsBz = []byte(attribute.GetValue())
// 				}
// 				if string(attribute.Key) == types.AttributeKeyRequestContextID &&
// 					string(attribute.GetValue()) == requestContextId {
// 					found = true
// 				}
// 			}
// 			s.Require().True(found)
// 			if found {
// 				err := json.Unmarshal(requestsBz, &requests)
// 				s.Require().NoError(err)
// 			}
// 			s.Require().Len(requests, 1)
// 			compactRequest = requests[0]
// 		}
// 	}
// 	s.Require().Equal(requestContextId, compactRequest.RequestContextId)

// 	//------test GetCmdQueryServiceRequests()-------------
// 	queryRequestsResponse := servicetestutil.QueryServiceRequestsExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		serviceName,
// 		provider.String(),
// 	)
// 	s.Require().Len(queryRequestsResponse.Requests, 1)
// 	s.Require().Equal(requestContextId, queryRequestsResponse.Requests[0].RequestContextId)

// 	//------test GetCmdQueryServiceRequests()-------------
// 	queryRequestsResponse = servicetestutil.QueryServiceRequestsByReqCtx(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		queryRequestsResponse.Requests[0].RequestContextId,
// 		fmt.Sprint(queryRequestsResponse.Requests[0].RequestContextBatchCounter),
// 	)
// 	s.Require().Len(queryRequestsResponse.Requests, 1)
// 	s.Require().Equal(requestContextId, queryRequestsResponse.Requests[0].RequestContextId)

// 	//------test GetCmdRespondService()-------------
// 	request := queryRequestsResponse.Requests[0]
// 	args = []string{
// 		fmt.Sprintf("--%s=%s", servicecli.FlagRequestID, request.Id),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagResult, respResult),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagData, respOutput),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = servicetestutil.RespondServiceExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		provider.String(),
// 		args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	//------test GetCmdQueryEarnedFees()-------------
// 	queryEarnedFeesResponse := servicetestutil.QueryEarnedFeesExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		provider.String(),
// 	)
// 	s.Require().Equal(expectedEarnedFees, queryEarnedFeesResponse.Fees.String())

// 	//------GetCmdSetWithdrawAddr()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = servicetestutil.SetWithdrawAddrExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		withdrawalAddress.String(),
// 		provider.String(),
// 		args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	//------GetCmdWithdrawEarnedFees()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = servicetestutil.WithdrawEarnedFeesExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		provider.String(),
// 		provider.String(),
// 		args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	withdrawalFees := simapp.QueryBalancesExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		withdrawalAddress.String(),
// 	)
// 	s.Require().Equal(expectedEarnedFees, withdrawalFees.String())

// 	//------check service tax-------------
// 	taxFees := simapp.QueryBalancesExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		authtypes.NewModuleAddress(servicetypes.FeeCollectorName).String(),
// 	)
// 	s.Require().Equal(expectedTaxFees, taxFees.String())

// 	//------GetCmdQueryRequestContext()-------------
// 	contextId := request.RequestContextId
// 	contextResp := servicetestutil.QueryRequestContextExec(s.T(), s.network, clientCtx, contextId)
// 	s.Require().False(contextResp.Empty())

// 	//------GetCmdQueryServiceRequest()-------------
// 	requestId := request.Id
// 	requestResp := servicetestutil.QueryServiceRequestExec(s.T(), s.network, clientCtx, requestId)
// 	s.Require().False(requestResp.Empty())
// 	s.Require().Equal(requestId, requestResp.Id)

// 	//------GetCmdQueryServiceResponse()-------------
// 	responseResp := servicetestutil.QueryServiceResponseExec(s.T(), s.network, clientCtx, requestId)
// 	s.Require().False(responseResp.Empty())
// }
