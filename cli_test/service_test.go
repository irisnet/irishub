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

	guardianAddr := ""
	taxWithdrawAddr := ""
	taxWithdrawAmt := int64(100)

	// define service

	ioutil.WriteFile(serviceFileName, []byte(serviceIDLContent), 0644)
	defer tests.ExecuteT(t, fmt.Sprintf("rm -f %s", serviceFileName), "")

	success, _, _ := f.TxServiceDefine(serviceName, serviceDesc, serviceTags, serviceAuthorDesc, serviceFileName, provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	svcDef := f.QueryServiceDefinition(serviceName, chainID)
	require.Equal(t, serviceName, svcDef.Name)
	require.Equal(t, chainID, svcDef.ChainId)

	// bind service

	success, _, _ = f.TxServiceBind(serviceName, chainID, bindingType, serviceDeposit, servicePrices, avgRspTime, usableTime, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	binding := f.QueryServiceBinding(serviceName, chainID, chainID, provider)
	require.Equal(t, serviceName, binding.DefName)
	require.Equal(t, chainID, binding.DefChainID)
	require.Equal(t, chainID, binding.BindChainID)
	require.Equal(t, provider, binding.Provider)
	require.Equal(t, serviceDeposit, binding.Deposit)

	bindings := f.QueryServiceBindings(serviceName, chainID)
	require.Equal(t, 1, len(bindings))

	// TODO: provider and module account balance check
	// So are the following tests

	// update service binding

	success, _, _ = f.TxServiceBindingUpdate(serviceName, chainID, bindingType, serviceDeposit, servicePrices, avgRspTime, usableTime, provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	binding = f.QueryServiceBinding(serviceName, chainID, chainID, provider)
	require.Equal(t, serviceDeposit, binding.Deposit.String())

	// disable service binding
	success, _, _ = f.TxServiceBindingDisable(serviceName, chainID, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	binding = f.QueryServiceBinding(serviceName, chainID, chainID, provider)
	require.Equal(t, false, binding.Available)

	// refund deposit
	success, _, _ = f.TxServiceBindingRefund(serviceName, chainID, provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	binding = f.QueryServiceBinding(serviceName, chainID, chainID, provider)
	require.Equal(t, serviceName, binding.Deposit.String())
	require.Equal(t, false, binding.Available)

	// service call
	success, _, _ = f.TxServiceCall(serviceName, chainID, chainID, provider, reqMethodID, reqInput, reqServiceFees, consumer, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	// TODO
	reqID := ""

	requests := f.QueryServiceRequests(serviceName, chainID, chainID, consumer)
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
	success, _, _ = f.TxServiceFeesRefund(consumer, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	fees = f.QueryServiceFees(consumer)
	require.Nil(t, fees.IncomingFee)
	require.Equal(t, "", fees.ReturnedFee)

	// withdraw fees
	success, _, _ = f.TxServiceFeesWithdraw(provider, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)

	fees = f.QueryServiceFees(provider)
	require.Nil(t, fees.ReturnedFee)
	require.Equal(t, "", fees.IncomingFee)

	// withdraw tax
	success, _, _ = f.TxServiceTaxWithdraw(taxWithdrawAmt, taxWithdrawAddr, guardianAddr, "-y")
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, f.Port)
}

// TxServiceDefinition is iriscli tx service define
func (f *Fixtures) TxServiceDefine(serviceName, serviceDesc, tags, serviceAuthorDesc, serviceFileName, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service define --service-name %s --service-description %s --tags %s --author-description %s --file %s --from=%s %v", f.IriscliBinary, serviceName, serviceDesc, tags, serviceAuthorDesc, serviceFileName, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceDefinition is iriscli tx service bind
func (f *Fixtures) TxServiceBind(serviceName, defChainID, bindType, deposit, prices string, avgRspTime int64, usableTime int64, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service bind --service-name %s --def-chain-id %s --bind-type %s --deposit %s --prices %s --avg-rsp-time %d --usable-time %d --from=%s %v", f.IriscliBinary, serviceName, defChainID, bindType, deposit, prices, avgRspTime, usableTime, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceBindingRefund is iriscli tx service binding refund
func (f *Fixtures) TxServiceBindingRefund(serviceName, defChainID, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service binding refund --service-name %s --def-chain-id %s --from=%s %v", f.IriscliBinary, serviceName, defChainID, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceBindingUpdate is iriscli tx service binding update
func (f *Fixtures) TxServiceBindingUpdate(serviceName, defChainID, bindType, deposit, prices string, avgRspTime int64, usableTime int64, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service binding update --service-name %s --def-chain-id %s --bind-type %s --deposit %s --prices %s --avg-rsp-time %d --usable-time %d --from=%s %v", f.IriscliBinary, serviceName, defChainID, bindType, deposit, prices, avgRspTime, usableTime, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceBindingDisable is iriscli tx service binding disable
func (f *Fixtures) TxServiceBindingDisable(serviceName, defChainID, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service binding disable --service-name %s --def-chain-id %s --from=%s %v", f.IriscliBinary, serviceName, defChainID, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceCall is iriscli tx service call
func (f *Fixtures) TxServiceCall(serviceName, defChainID, bindChainID, provider string, methodID int16, requestData, serviceFees, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service call --service-name %s --def-chain-id %s --bind-chain-id %s --provider %s --method-id %d --request-data %s --service-fees %s --from=%s %v", f.IriscliBinary, serviceName, defChainID, bindChainID, provider, methodID, requestData, serviceFees, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceRespond is iriscli tx service respond
func (f *Fixtures) TxServiceRespond(reqChainID, requestID, responseData, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service respond --request-chain-id %s --request-id %s --response-data %s --from=%s %v", f.IriscliBinary, reqChainID, requestID, responseData, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceFeesRefund is iriscli tx service fees refund
func (f *Fixtures) TxServiceFeesRefund(from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service fees refund --from=%s %v", f.IriscliBinary, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceFeesWithdraw is iriscli tx service fees withdraw
func (f *Fixtures) TxServiceFeesWithdraw(from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service fees withdraw --from=%s %v", f.IriscliBinary, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxServiceFeesWithdraw is iriscli tx service tax withdraw
func (f *Fixtures) TxServiceTaxWithdraw(withdrawAmt int64, destAddr string, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx service tax withdraw --withdraw-amount %d --dest-address %s --from=%s %v", f.IriscliBinary, withdrawAmt, destAddr, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// QueryServiceDefinition is iriscli query service definition
func (f *Fixtures) QueryServiceDefinition(serviceName, defChainID string) service.SvcDef {
	cmd := fmt.Sprintf("%s query service definition %s %s", f.IriscliBinary, serviceName, defChainID, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var svcDef service.SvcDef
	err := cdc.UnmarshalJSON([]byte(res), &svcDef)
	require.NoError(f.T, err)
	return svcDef
}

// QueryServiceBinding is iriscli query service binding
func (f *Fixtures) QueryServiceBinding(serviceName, defChainID, bindChainID, provider string) service.SvcBinding {
	cmd := fmt.Sprintf("%s query service binding %s %s %s %s", f.IriscliBinary, serviceName, defChainID, bindChainID, provider, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var binding service.SvcBinding
	err := cdc.UnmarshalJSON([]byte(res), &binding)
	require.NoError(f.T, err)
	return binding
}

// QueryServiceBindings is iriscli query service bindings
func (f *Fixtures) QueryServiceBindings(serviceName, defChainID string) []service.SvcBinding {
	cmd := fmt.Sprintf("%s query service bindings %s %s", f.IriscliBinary, serviceName, defChainID, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var bindings []service.SvcBinding
	err := cdc.UnmarshalJSON([]byte(res), &bindings)
	require.NoError(f.T, err)
	return bindings
}

// QueryServiceRequests is iriscli query service requests
func (f *Fixtures) QueryServiceRequests(serviceName, defChainID, bindChainID, consumer string) []service.SvcRequest {
	cmd := fmt.Sprintf("%s query service requests %s %s %s %s", f.IriscliBinary, serviceName, defChainID, bindChainID, consumer, f.Flags())
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

// // executeGetProfilers is iriscli query service bindings
// func (f *Fixtures) QueryServiceFees(t *testing.T, cmdStr string) []guardian.Guardian {
// 	cmd := fmt.Sprintf("%s query service definition %s %s", f.IriscliBinary, serviceName, defChainID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var guardians []guardian.Guardian
// 	err := cdc.UnmarshalJSON([]byte(res), &guardians)
// 	require.NoError(f.T, err)
// 	return guardians
// }
