package clitest

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/modules/service"
)

func TestIrisCLIService(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	cdc := app.MakeCodec()

	serviceFeeTax := sdk.NewDecWithPrec(1, 1) // 10%
	maxRequestTimeout := int64(2)

	// Update service params for test
	genesisState := f.GenesisState()
	var serviceData service.GenesisState
	err := cdc.UnmarshalJSON(genesisState[service.ModuleName], &serviceData)
	require.NoError(t, err)
	serviceData.Params.ServiceFeeTax = serviceFeeTax
	serviceData.Params.MaxRequestTimeout = maxRequestTimeout
	serviceDataBz, err := cdc.MarshalJSON(serviceData)
	require.NoError(t, err)
	genesisState[service.ModuleName] = serviceDataBz

	genFile := filepath.Join(f.IrisHome, "config", "genesis.json")
	genDoc, err := tmtypes.GenesisDocFromFile(genFile)
	require.NoError(t, err)
	genDoc.AppState, err = cdc.MarshalJSON(genesisState)
	require.NoError(t, genDoc.SaveAs(genFile))

	proc := f.GDStart()
	defer proc.Stop(false)

	tests.WaitForTMStart(f.Port)
	tests.WaitForNextNBlocksTM(2, f.Port)

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	sendTokens := sdk.TokensFromConsensusPower(10)
	f.TxSend(keyFoo, barAddr, sdk.NewCoin(denom, sendTokens), "-y")
	tests.WaitForNextNBlocksTM(1, f.Port)

	fooAcc := f.QueryAccount(fooAddr)
	fooCoinAmt := fooAcc.Coins.AmountOf(sdk.DefaultBondDenom).String()
	require.Equal(t, "40000000", fooCoinAmt)

	barAcc := f.QueryAccount(barAddr)
	barCoinAmt := barAcc.Coins.AmountOf(sdk.DefaultBondDenom).String()
	require.Equal(t, "10000000", barCoinAmt)

	// testing variables
	chainID := f.ChainID
	serviceName := "test-service"
	serviceDesc := "test-description"
	serviceAuthorDesc := "test-author-description"
	serviceTags := "tags1,tags2"
	serviceSchemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	serviceDenom := sdk.DefaultBondDenom

	serviceDeposit := fmt.Sprintf("50000%s", serviceDenom)
	servicePrices := fmt.Sprintf("50%s", serviceDenom)
	bindingType := "Local"
	avgRspTime := int64(10000)
	usableTime := int64(9999)

	reqMethodID := int16(1)
	reqServiceFees := fmt.Sprintf("50%s", serviceDenom)
	reqInput := "AB"
	respOutput := "CD"

	author := fooAddr.String()
	provider := fooAddr.String()
	consumer := barAddr.String()

	guardianAddr := fooAddr
	taxWithdrawAddr := barAddr
	taxWithdrawAmt := fmt.Sprintf("5%s", serviceDenom)

	// define service
	success, _, _ := f.TxServiceDefine(serviceName, serviceDesc, serviceTags, serviceAuthorDesc, serviceSchemas, author, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(1, f.Port)

	svcDef := f.QueryServiceDefinition(serviceName)
	require.Equal(t, serviceName, svcDef.Name)

	// bind service
	success, _, _ = f.TxServiceBind(chainID, serviceName, bindingType, serviceDeposit, servicePrices, avgRspTime, usableTime, provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(1, f.Port)

	binding := f.QueryServiceBinding(chainID, serviceName, chainID, provider)
	require.Equal(t, serviceName, binding.DefName)
	require.Equal(t, chainID, binding.DefChainID)
	require.Equal(t, chainID, binding.BindChainID)
	require.Equal(t, provider, binding.Provider.String())
	require.Equal(t, serviceDeposit, binding.Deposit.String())

	bindings := f.QueryServiceBindings(chainID, serviceName)
	require.Equal(t, 1, len(bindings))

	// disable service binding
	success, _, _ = f.TxServiceDisable(chainID, serviceName, provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(1, f.Port)

	binding = f.QueryServiceBinding(chainID, serviceName, chainID, provider)
	require.Equal(t, false, binding.Available)

	// refund deposit
	success, _, _ = f.TxServiceRefundDeposit(chainID, serviceName, provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(1, f.Port)

	binding = f.QueryServiceBinding(chainID, serviceName, chainID, provider)
	require.Equal(t, serviceDeposit, binding.Deposit.String())

	// enable service binding
	success, _, _ = f.TxServiceEnable(chainID, serviceName, serviceDeposit, provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(1, f.Port)

	binding = f.QueryServiceBinding(chainID, serviceName, chainID, provider)
	require.Equal(t, true, binding.Available)

	// service call
	success, _, _ = f.TxServiceCall(chainID, serviceName, chainID, provider, reqMethodID, reqInput, reqServiceFees, consumer, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(1, f.Port)

	// QueryTxs
	searchResult := f.QueryTxs(1, 50, "message.action=call_service", fmt.Sprintf("message.sender=%s", consumer))
	require.Len(t, searchResult.Txs, 1)

	events := searchResult.Txs[0].Logs[0].Events
	var reqID string
	for _, e := range events {
		for _, attribute := range e.Attributes {
			if attribute.Key == service.AttributeKeyRequestID {
				reqID = attribute.Value
				break
			}
		}
	}
	require.NotEmpty(t, reqID)
	requests := f.QueryServiceRequests(chainID, serviceName, chainID, provider)
	require.Equal(t, 1, len(requests))
	require.Equal(t, reqID, requests[0].RequestID())
	require.Equal(t, consumer, requests[0].Consumer.String())
	require.Equal(t, provider, requests[0].Provider.String())

	// respond service request
	success, _, _ = f.TxServiceRespond(chainID, reqID, respOutput, provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(1, f.Port)

	// query fees
	fees := f.QueryServiceFees(provider)
	require.Nil(t, fees.ReturnedFee)
	require.Equal(t, "45stake", fees.IncomingFee.String()) // servicePrices * (1-serviceFeeTax)

	fees = f.QueryServiceFees(consumer)
	require.Nil(t, fees.IncomingFee)
	require.Nil(t, fees.ReturnedFee)

	// withdraw fees
	success, _, _ = f.TxServiceWithdrawFees(provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(1, f.Port)

	fees = f.QueryServiceFees(provider)
	require.Nil(t, fees.ReturnedFee)
	require.Nil(t, fees.IncomingFee)

	// service call but does not respond
	success, _, _ = f.TxServiceCall(chainID, serviceName, chainID, provider, reqMethodID, reqInput, reqServiceFees, consumer, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(maxRequestTimeout+1, f.Port)

	// query fees
	fees = f.QueryServiceFees(provider)
	require.Nil(t, fees.IncomingFee)
	require.Nil(t, fees.ReturnedFee)

	fees = f.QueryServiceFees(consumer)
	require.Equal(t, servicePrices, fees.ReturnedFee.String())
	require.Nil(t, fees.IncomingFee)

	// refund fees
	success, _, _ = f.TxServiceRefundFees(consumer, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(1, f.Port)

	fees = f.QueryServiceFees(consumer)
	require.Nil(t, fees.IncomingFee)
	require.Nil(t, fees.ReturnedFee)

	// withdraw tax
	success, _, _ = f.TxServiceWithdrawTax(taxWithdrawAmt, taxWithdrawAddr, guardianAddr, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(1, f.Port)

	// Cleanup testing directories
	f.Cleanup()
}

// TxServiceDefine is iriscli tx service define
func (f *Fixtures) TxServiceDefine(serviceName, serviceDesc, tags, serviceAuthorDesc, serviceSchemas, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service define --keyring-backend=test --name %s --description %s --tags %s --author-description %s --schemas %s --from=%s %v", f.IriscliBinary, serviceName, serviceDesc, tags, serviceAuthorDesc, serviceSchemas, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxServiceBind is iriscli tx service bind
func (f *Fixtures) TxServiceBind(defChainID, serviceName, bindType, deposit, prices string, avgRspTime int64, usableTime int64, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service bind --keyring-backend=test --def-chain-id %s --service-name %s --bind-type %s --deposit %s --prices %s --avg-rsp-time %d --usable-time %d --from=%s %v", f.IriscliBinary, defChainID, serviceName, bindType, deposit, prices, avgRspTime, usableTime, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxServiceRefundDeposit is iriscli tx service refund-deposit
func (f *Fixtures) TxServiceRefundDeposit(defChainID, serviceName, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service refund-deposit --keyring-backend=test --def-chain-id %s --service-name %s --from=%s %v", f.IriscliBinary, defChainID, serviceName, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxServiceUpdateBinding is iriscli tx service update-binding
func (f *Fixtures) TxServiceUpdateBinding(defChainID, serviceName, bindType, deposit, prices string, avgRspTime int64, usableTime int64, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service update-binding --keyring-backend=test --def-chain-id %s --service-name %s --bind-type %s --deposit %s --prices %s --avg-rsp-time %d --usable-time %d --from=%s %v", f.IriscliBinary, defChainID, serviceName, bindType, deposit, prices, avgRspTime, usableTime, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxServiceDisable is iriscli tx service disable
func (f *Fixtures) TxServiceDisable(defChainID, serviceName, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service disable --keyring-backend=test --def-chain-id %s --service-name %s --from=%s %v", f.IriscliBinary, defChainID, serviceName, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxServiceEnable is iriscli tx service enable
func (f *Fixtures) TxServiceEnable(defChainID, serviceName, deposit, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service enable --keyring-backend=test --def-chain-id %s --service-name %s --deposit %s --from=%s %v", f.IriscliBinary, defChainID, serviceName, deposit, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxServiceCall is iriscli tx service call
func (f *Fixtures) TxServiceCall(defChainID, serviceName, bindChainID, provider string, methodID int16, requestData, serviceFees, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service call --keyring-backend=test --def-chain-id %s --service-name %s --bind-chain-id %s --provider %s --method-id %d --request-data %s --service-fee %s --from=%s %v", f.IriscliBinary, defChainID, serviceName, bindChainID, provider, methodID, requestData, serviceFees, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxServiceRespond is iriscli tx service respond
func (f *Fixtures) TxServiceRespond(reqChainID, requestID, responseData, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service respond --keyring-backend=test --request-chain-id %s --request-id %s --response-data %s --from=%s %v", f.IriscliBinary, reqChainID, requestID, responseData, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxServiceRefundFees is iriscli tx service refund-fees
func (f *Fixtures) TxServiceRefundFees(from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service refund-fees --keyring-backend=test --from=%s %v", f.IriscliBinary, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxServiceWithdrawFees is iriscli tx service withdraw-fees
func (f *Fixtures) TxServiceWithdrawFees(from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service withdraw-fees --keyring-backend=test --from=%s %v", f.IriscliBinary, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxServiceWithdrawTax is iriscli tx service withdraw-tax
func (f *Fixtures) TxServiceWithdrawTax(withdrawAmt string, destAddr, from sdk.AccAddress, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service withdraw-tax --keyring-backend=test --dest-address %s --withdraw-amount %s --from=%s %v", f.IriscliBinary, destAddr, withdrawAmt, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// QueryServiceDefinition is iriscli query service definition
func (f *Fixtures) QueryServiceDefinition(serviceName string) service.ServiceDefinition {
	cmd := fmt.Sprintf("%s query service definition %s %v", f.IriscliBinary, serviceName, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var svcDef service.ServiceDefinition
	err := cdc.UnmarshalJSON([]byte(res), &svcDef)
	require.NoError(f.T, err)
	return svcDef
}

// QueryServiceBinding is iriscli query service binding
func (f *Fixtures) QueryServiceBinding(defChainID, serviceName, bindChainID, provider string) service.SvcBinding {
	cmd := fmt.Sprintf("%s query service binding %s %s %s %s %v", f.IriscliBinary, defChainID, serviceName, bindChainID, provider, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var binding service.SvcBinding
	err := cdc.UnmarshalJSON([]byte(res), &binding)
	require.NoError(f.T, err)
	return binding
}

// QueryServiceBindings is iriscli query service bindings
func (f *Fixtures) QueryServiceBindings(defChainID, serviceName string) []service.SvcBinding {
	cmd := fmt.Sprintf("%s query service bindings %s %s %v", f.IriscliBinary, defChainID, serviceName, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var bindings []service.SvcBinding
	err := cdc.UnmarshalJSON([]byte(res), &bindings)
	require.NoError(f.T, err)
	return bindings
}

// QueryServiceRequests is iriscli query service requests
func (f *Fixtures) QueryServiceRequests(defChainID, serviceName, bindChainID, provider string) []service.SvcRequest {
	cmd := fmt.Sprintf("%s query service requests %s %s %s %s %v", f.IriscliBinary, defChainID, serviceName, bindChainID, provider, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var reqs []service.SvcRequest
	err := cdc.UnmarshalJSON([]byte(res), &reqs)
	require.NoError(f.T, err)
	return reqs
}

// QueryServiceFees is iriscli query service fees
func (f *Fixtures) QueryServiceFees(address string) service.FeesOutput {
	cmd := fmt.Sprintf("%s query service fees %s %v", f.IriscliBinary, address, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var fees service.FeesOutput
	err := cdc.UnmarshalJSON([]byte(res), &fees)
	require.NoError(f.T, err)
	return fees
}
