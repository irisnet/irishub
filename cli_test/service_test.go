package clitest

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/modules/service"
)

func TestIrisCLIService(t *testing.T) {
	t.Parallel()

	f := InitFixtures(t)

	proc := f.GDStart()
	defer proc.Stop(false)

	tests.WaitForTMStart(f.Port)
	tests.WaitForNextNBlocksTM(2, f.Port)

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	fooAcc := f.QueryAccount(fooAddr)
	fooCoinAmt := fooAcc.Coins.AmountOf(sdk.DefaultBondDenom).String()
	require.Equal(t, "50000000", fooCoinAmt)

	// testing variables

	chainID := f.ChainID
	serviceName := "test-service"
	serviceDesc := "test-description"
	serviceAuthorDesc := "test-author-description"
	serviceTags := "tags1,tags2"
	serviceIDLContent := ""
	serviceFileName := f.IriscliHome + string(os.PathSeparator) + "test.proto"

	serviceDeposit := fmt.Sprintf("100%s", service.ServiceDenom)
	servicePrices := fmt.Sprintf("50%s", service.ServiceDenom)
	bindingType := "Local"
	avgRspTime := int64(10000)
	usableTime := int64(9999)

	reqMethodID := int16(0)
	reqServiceFees := fmt.Sprintf("50%s", service.ServiceDenom)
	reqInput := "0x"
	respOutput := "0x"

	consumer := barAddr.String()
	provider := fooAddr.String()

	// TODO
	guardianAddr := ""
	taxWithdrawAddr := ""
	taxWithdrawAmt := fmt.Sprintf("0.5%s", service.ServiceDenom)

	// define service

	ioutil.WriteFile(serviceFileName, []byte(serviceIDLContent), 0644)
	defer tests.ExecuteT(t, fmt.Sprintf("rm -f %s", serviceFileName), "")

	success, _, _ := f.TxServiceDefine(serviceName, serviceDesc, serviceTags, serviceAuthorDesc, serviceIDLContent, serviceFileName, provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	svcDefOutput := f.QueryServiceDefinition(serviceName, chainID)
	require.Equal(t, serviceName, svcDefOutput.Definition.Name)
	require.Equal(t, chainID, svcDefOutput.Definition.ChainId)

	// bind service

	success, _, _ = f.TxServiceBind(chainID, serviceName, bindingType, serviceDeposit, servicePrices, avgRspTime, usableTime, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	binding := f.QueryServiceBinding(chainID, serviceName, chainID, provider)
	require.Equal(t, serviceName, binding.DefName)
	require.Equal(t, chainID, binding.DefChainID)
	require.Equal(t, chainID, binding.BindChainID)
	require.Equal(t, provider, binding.Provider)
	require.Equal(t, serviceDeposit, binding.Deposit)

	bindings := f.QueryServiceBindings(chainID, serviceName)
	require.Equal(t, 1, len(bindings))

	// TODO: provider and module account balance check
	// So are the following tests

	// update service binding

	success, _, _ = f.TxServiceUpdateBinding(chainID, serviceName, bindingType, serviceDeposit, servicePrices, avgRspTime, usableTime, provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	binding = f.QueryServiceBinding(chainID, serviceName, chainID, provider)
	require.Equal(t, serviceDeposit, binding.Deposit)

	// disable service binding
	success, _, _ = f.TxServiceDisable(chainID, serviceName, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	binding = f.QueryServiceBinding(chainID, serviceName, chainID, provider)
	require.Equal(t, false, binding.Available)

	// refund deposit
	success, _, _ = f.TxServiceRefundDeposit(chainID, serviceName, provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	binding = f.QueryServiceBinding(chainID, serviceName, chainID, provider)
	require.Equal(t, "", binding.Deposit)

	// service call
	success, _, _ = f.TxServiceCall(chainID, serviceName, chainID, provider, reqMethodID, reqInput, reqServiceFees, consumer, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	// TODO
	reqID := ""

	requests := f.QueryServiceRequests(chainID, serviceName, chainID, consumer)
	require.Equal(t, 1, len(requests))
	require.Equal(t, reqID, requests[0].RequestID())
	require.Equal(t, consumer, requests[0].Consumer)
	require.Equal(t, provider, requests[0].Provider)

	// respond service request
	success, _, _ = f.TxServiceRespond(chainID, reqID, respOutput, provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	// query fees
	fees := f.QueryServiceFees(provider)
	require.Nil(t, fees.ReturnedFee)
	require.Equal(t, "", fees.IncomingFee)

	fees = f.QueryServiceFees(consumer)
	require.Nil(t, fees.IncomingFee)
	require.Equal(t, "", fees.ReturnedFee)

	// refund fees
	success, _, _ = f.TxServiceRefundFees(consumer, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	fees = f.QueryServiceFees(consumer)
	require.Nil(t, fees.IncomingFee)
	require.Equal(t, "", fees.ReturnedFee)

	// withdraw fees
	success, _, _ = f.TxServiceWithdrawFees(provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	fees = f.QueryServiceFees(provider)
	require.Nil(t, fees.ReturnedFee)
	require.Equal(t, "", fees.IncomingFee)

	// withdraw tax
	success, _, _ = f.TxServiceWithdrawTax(taxWithdrawAddr, taxWithdrawAmt, guardianAddr, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)
}

// TxServiceDefine is iriscli tx service define
func (f *Fixtures) TxServiceDefine(serviceName, serviceDesc, tags, serviceAuthorDesc, serviceIDLContent, serviceFileName, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service define --service-name %s --service-description %s --tags %s --author-description %s --idl-content %s --file %s --from=%s %v", f.IriscliBinary, serviceName, serviceDesc, tags, serviceAuthorDesc, serviceIDLContent, serviceFileName, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceBind is iriscli tx service bind
func (f *Fixtures) TxServiceBind(serviceName, defChainID, bindType, deposit, prices string, avgRspTime int64, usableTime int64, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service bind --def-chain-id %s --service-name %s --bind-type %s --deposit %s --prices %s --avg-rsp-time %d --usable-time %d --from=%s %v", f.IriscliBinary, defChainID, serviceName, bindType, deposit, prices, avgRspTime, usableTime, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceRefundDeposit is iriscli tx service refund-deposit
func (f *Fixtures) TxServiceRefundDeposit(defChainID, serviceName, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service refund-deposit --def-chain-id %s --service-name %s --from=%s %v", f.IriscliBinary, defChainID, serviceName, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceUpdateBinding is iriscli tx service update-binding
func (f *Fixtures) TxServiceUpdateBinding(defChainID, serviceName, bindType, deposit, prices string, avgRspTime int64, usableTime int64, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service update-binding --def-chain-id %s --service-name %s --bind-type %s --deposit %s --prices %s --avg-rsp-time %d --usable-time %d --from=%s %v", f.IriscliBinary, defChainID, serviceName, bindType, deposit, prices, avgRspTime, usableTime, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceDisable is iriscli tx service disable
func (f *Fixtures) TxServiceDisable(defChainID, serviceName, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service disable --def-chain-id %s --service-name %s --from=%s %v", f.IriscliBinary, defChainID, serviceName, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceCall is iriscli tx service call
func (f *Fixtures) TxServiceCall(defChainID, serviceName, bindChainID, provider string, methodID int16, requestData, serviceFees, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service call --def-chain-id %s --service-name %s --bind-chain-id %s --provider %s --method-id %d --request-data %s --service-fees %s --from=%s %v", f.IriscliBinary, defChainID, serviceName, bindChainID, provider, methodID, requestData, serviceFees, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceRespond is iriscli tx service respond
func (f *Fixtures) TxServiceRespond(reqChainID, requestID, responseData, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service respond --request-chain-id %s --request-id %s --response-data %s --from=%s %v", f.IriscliBinary, reqChainID, requestID, responseData, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceRefundFees is iriscli tx service refund-fees
func (f *Fixtures) TxServiceRefundFees(from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service refund-fees --from=%s %v", f.IriscliBinary, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceWithdrawFees is iriscli tx service withdraw-fees
func (f *Fixtures) TxServiceWithdrawFees(from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service withdraw-fees --from=%s %v", f.IriscliBinary, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceWithdrawTax is iriscli tx service withdraw-tax
func (f *Fixtures) TxServiceWithdrawTax(destAddr string, withdrawAmt string, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service withdraw-tax --dest-address %s --withdraw-amount %s --from=%s %v", f.IriscliBinary, destAddr, withdrawAmt, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// QueryServiceDefinition is iriscli query service definition
func (f *Fixtures) QueryServiceDefinition(defChainID, serviceName string) service.DefinitionOutput {
	cmd := fmt.Sprintf("%s query service definition %s %s", f.IriscliBinary, defChainID, serviceName, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var svcDefOutput service.DefinitionOutput
	err := cdc.UnmarshalJSON([]byte(res), &svcDefOutput)
	require.NoError(f.T, err)
	return svcDefOutput
}

// QueryServiceBinding is iriscli query service binding
func (f *Fixtures) QueryServiceBinding(defChainID, serviceName, bindChainID, provider string) service.SvcBinding {
	cmd := fmt.Sprintf("%s query service binding %s %s %s %s", f.IriscliBinary, defChainID, serviceName, bindChainID, provider, f.Flags())
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
	cmd := fmt.Sprintf("%s query service bindings %s %s", f.IriscliBinary, defChainID, serviceName, f.Flags())
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
	cmd := fmt.Sprintf("%s query service requests %s %s %s %s", f.IriscliBinary, defChainID, serviceName, bindChainID, provider, f.Flags())
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
	cmd := fmt.Sprintf("%s query service fees %s ", f.IriscliBinary, address, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var fees service.FeesOutput
	err := cdc.UnmarshalJSON([]byte(res), &fees)
	require.NoError(f.T, err)
	return fees
}
