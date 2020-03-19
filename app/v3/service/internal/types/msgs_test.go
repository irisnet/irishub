package types

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/irisnet/irishub/types"
)

var (
	testCoin1, _ = sdk.IrisCoinType.ConvertToCoin("1000iris")
	testCoin2, _ = sdk.IrisCoinType.ConvertToCoin("100iris")

	testServiceName = "test-service"
	testServiceDesc = "test-service-desc"
	testServiceTags = []string{"tag1", "tag2"}
	testAuthor      = sdk.AccAddress([]byte("test-author"))
	testAuthorDesc  = "test-author-desc"
	testSchemas     = `{"input":{"type":"object"},"output":{"type":"object"}}`

	testProvider     = sdk.AccAddress([]byte("test-provider"))
	testDeposit      = sdk.NewCoins(testCoin1)
	testPricing      = `{"price":[{"denom":"iris-atto","amount":"1000000"}]}`
	testWithdrawAddr = sdk.AccAddress([]byte("test-withdrawal-address"))
	testAddedDeposit = sdk.NewCoins(testCoin2)

	testConsumer      = sdk.AccAddress([]byte("test-consumer"))
	testProviders     = []sdk.AccAddress{testProvider}
	testInput         = `{"pair":"iris-usdt"}`
	testServiceFeeCap = sdk.NewCoins(testCoin2)
	testTimeout       = int64(100)
	testRepeatedFreq  = uint64(120)
	testRepeatedTotal = int64(100)

	testResult           = `{"code":200,"message":""}`
	testOutput           = `{"last":"100"}`
	testTrustee          = sdk.AccAddress([]byte("test-trustee"))
	testTaxWithdrawalAmt = sdk.NewCoins(testCoin1)

	testRequestContextID = GenerateRequestContextID(tmhash.Sum([]byte("test-request-context-id")), 0)
	testRequestID        = GenerateRequestID(testRequestContextID, 1, 1, 1).String()
)

// TestMsgDefineServiceRoute tests Route for MsgDefineService
func TestMsgDefineServiceRoute(t *testing.T) {
	msg := NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgDefineServiceType tests Type for MsgDefineService
func TestMsgDefineServiceType(t *testing.T) {
	msg := NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)

	require.Equal(t, "define_service", msg.Type())
}

// TestMsgDefineServiceValidation tests ValidateBasic for MsgDefineService
func TestMsgDefineServiceValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidLongDesc := strings.Repeat("d", MaxDescriptionLength+1)
	invalidMoreTags := strings.Split("t1,t2,t3,t4,t5,t6,t7,t8,t9,t10,t11", ",")
	invalidLongTags := []string{strings.Repeat("t", MaxTagLength+1)}
	invalidEmptyTags := []string{"t1", ""}
	invalidDuplicateTags := []string{"t1", "t1"}

	invalidSchemas := `{"input":"nonobject","output":"nonobject"}`
	invalidSchemasNoInput := `{"output":{"type":"object"}}`
	invalidSchemasNoOutput := `{"input":{"type":"object"}}`

	testMsgs := []MsgDefineService{
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas),            // valid msg
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, emptyAddress, testAuthorDesc, testSchemas),          // missing author address
		NewMsgDefineService(invalidName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas),                // service name contains illegal characters
		NewMsgDefineService(invalidLongName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas),            // too long service name
		NewMsgDefineService(testServiceName, invalidLongDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas),            // too long service description
		NewMsgDefineService(testServiceName, testServiceDesc, invalidMoreTags, testAuthor, testAuthorDesc, testSchemas),            // too many tags
		NewMsgDefineService(testServiceName, testServiceDesc, invalidLongTags, testAuthor, testAuthorDesc, testSchemas),            // too long tag
		NewMsgDefineService(testServiceName, testServiceDesc, invalidEmptyTags, testAuthor, testAuthorDesc, testSchemas),           // empty tag
		NewMsgDefineService(testServiceName, testServiceDesc, invalidDuplicateTags, testAuthor, testAuthorDesc, testSchemas),       // duplicate tags
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, invalidLongDesc, testSchemas),           // too long author description
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, invalidSchemas),         // invalid schemas
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, invalidSchemasNoInput),  // missing input schema
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, invalidSchemasNoOutput), // missing output schema
	}

	testCases := []struct {
		msg     MsgDefineService
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing author address"},
		{testMsgs[2], false, "service name contains illegal characters"},
		{testMsgs[3], false, "too long service name"},
		{testMsgs[4], false, "too long service description"},
		{testMsgs[5], false, "too many tags"},
		{testMsgs[6], false, "too long tag"},
		{testMsgs[7], false, "empty tag"},
		{testMsgs[8], false, "duplicate tags"},
		{testMsgs[9], false, "too long author description"},
		{testMsgs[10], false, "invalid schemas"},
		{testMsgs[11], false, "missing input schema"},
		{testMsgs[12], false, "missing output schema"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgDefineServiceGetSignBytes tests GetSignBytes for MsgDefineService
func TestMsgDefineServiceGetSignBytes(t *testing.T) {
	msg := NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgDefineService","value":{"author":"faa1w3jhxapdv96hg6r0wgux6kuw","author_description":"test-author-desc","description":"test-service-desc","name":"test-service","schemas":"{\"input\":{\"type\":\"object\"},\"output\":{\"type\":\"object\"}}","tags":["tag1","tag2"]}}`
	require.Equal(t, expected, string(res))
}

// TestMsgDefineServiceGetSigners tests GetSigners for MsgDefineService
func TestMsgDefineServiceGetSigners(t *testing.T) {
	msg := NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)
	res := msg.GetSigners()

	expected := "[746573742D617574686F72]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgBindServiceRoute tests Route for MsgBindService
func TestMsgBindServiceRoute(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgBindServiceType tests Type for MsgBindService
func TestMsgBindServiceType(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing)

	require.Equal(t, "bind_service", msg.Type())
}

// TestMsgBindServiceValidation tests ValidateBasic for MsgBindService
func TestMsgBindServiceValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidDeposit := sdk.Coins{}
	invalidDenomDeposit := sdk.NewCoins(sdk.NewCoin("eth-min", sdk.NewInt(1000)))
	invalidPricing := `{"price":[{"denom":"iris-atto","amount":"1000000"}],"other":"notallowedfield"}`
	invalidDenomPricing := `{"price":[{"denom":"iris","amount":"1000000"}]}`

	testMsgs := []MsgBindService{
		NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing),         // valid msg
		NewMsgBindService(testServiceName, emptyAddress, testDeposit, testPricing),         // missing provider address
		NewMsgBindService(invalidName, testProvider, testDeposit, testPricing),             // service name contains illegal characters
		NewMsgBindService(invalidLongName, testProvider, testDeposit, testPricing),         // too long service name
		NewMsgBindService(testServiceName, testProvider, invalidDeposit, testPricing),      // invalid deposit
		NewMsgBindService(testServiceName, testProvider, invalidDenomDeposit, testPricing), // invalid deposit denom
		NewMsgBindService(testServiceName, testProvider, testDeposit, ""),                  // missing pricing
		NewMsgBindService(testServiceName, testProvider, testDeposit, invalidPricing),      // invalid Pricing JSON Schema instance
		NewMsgBindService(testServiceName, testProvider, testDeposit, invalidDenomPricing), // invalid pricing denom
	}

	testCases := []struct {
		msg     MsgBindService
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
		{testMsgs[2], false, "service name contains illegal characters"},
		{testMsgs[3], false, "too long service name"},
		{testMsgs[4], false, "invalid deposit"},
		{testMsgs[5], false, "invalid deposit denom"},
		{testMsgs[6], false, "missing pricing"},
		{testMsgs[7], false, "invalid Pricing JSON Schema instance"},
		{testMsgs[8], false, "invalid pricing denom"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgBindServiceGetSignBytes tests GetSignBytes for MsgBindService
func TestMsgBindServiceGetSignBytes(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgBindService","value":{"deposit":[{"amount":"1000000000000000000000","denom":"iris-atto"}],"pricing":"{\"price\":[{\"denom\":\"iris-atto\",\"amount\":\"1000000\"}]}","provider":"faa1w3jhxapdwpex7anfv3jhynrxe9z","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgBindServiceGetSigners tests GetSigners for MsgBindService
func TestMsgBindServiceGetSigners(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgUpdateServiceBindingRoute tests Route for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingRoute(t *testing.T) {
	msg := NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "")

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgUpdateServiceBindingType tests Type for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingType(t *testing.T) {
	msg := NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "")

	require.Equal(t, "update_service_binding", msg.Type())
}

// TestMsgUpdateServiceBindingValidation tests ValidateBasic for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}
	emptyAddedDeposit := sdk.Coins{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidDenomDeposit := sdk.NewCoins(sdk.NewCoin("eth-min", sdk.NewInt(1000)))
	invalidPricing := `{"price":[{"denom":"iris-atto","amount":"1000000"}],"other":"notallowedfield"}`
	invalidDenomPricing := `{"price":[{"denom":"iris","amount":"1000000"}]}`

	testMsgs := []MsgUpdateServiceBinding{
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, testPricing),         // valid msg
		NewMsgUpdateServiceBinding(testServiceName, testProvider, emptyAddedDeposit, testPricing),        // empty deposit is allowed
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, ""),                  // empty pricing is allowed
		NewMsgUpdateServiceBinding(testServiceName, testProvider, emptyAddedDeposit, ""),                 // deposit and pricing can be empty at the same time
		NewMsgUpdateServiceBinding(testServiceName, emptyAddress, testAddedDeposit, testPricing),         // missing provider address
		NewMsgUpdateServiceBinding(invalidName, testProvider, testAddedDeposit, testPricing),             // service name contains illegal characters
		NewMsgUpdateServiceBinding(invalidLongName, testProvider, testAddedDeposit, testPricing),         // too long service name
		NewMsgUpdateServiceBinding(testServiceName, testProvider, invalidDenomDeposit, testPricing),      // invalid deposit denom
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, invalidPricing),      // invalid Pricing JSON Schema instance
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, invalidDenomPricing), // invalid pricing denom
	}

	testCases := []struct {
		msg     MsgUpdateServiceBinding
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], true, ""},
		{testMsgs[2], true, ""},
		{testMsgs[3], true, ""},
		{testMsgs[4], false, "missing provider address"},
		{testMsgs[5], false, "service name contains illegal characters"},
		{testMsgs[6], false, "too long service name"},
		{testMsgs[7], false, "invalid deposit denom"},
		{testMsgs[8], false, "invalid Pricing JSON Schema instance"},
		{testMsgs[9], false, "invalid pricing denom"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgUpdateServiceBindingGetSignBytes tests GetSignBytes for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingGetSignBytes(t *testing.T) {
	msg := NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "")
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgUpdateServiceBinding","value":{"deposit":[{"amount":"100000000000000000000","denom":"iris-atto"}],"pricing":"","provider":"faa1w3jhxapdwpex7anfv3jhynrxe9z","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgUpdateServiceBindingGetSigners tests GetSigners for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingGetSigners(t *testing.T) {
	msg := NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "")
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgSetWithdrawAddressRoute tests Route for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressRoute(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testProvider, testWithdrawAddr)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgSetWithdrawAddressType tests Type for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressType(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testProvider, testWithdrawAddr)

	require.Equal(t, "set_withdraw_address", msg.Type())
}

// TestMsgSetWithdrawAddressValidation tests ValidateBasic for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	testMsgs := []MsgSetWithdrawAddress{
		NewMsgSetWithdrawAddress(testProvider, testWithdrawAddr), // valid msg
		NewMsgSetWithdrawAddress(emptyAddress, testWithdrawAddr), // missing provider address
		NewMsgSetWithdrawAddress(testProvider, emptyAddress),     // missing withdrawal address
	}

	testCases := []struct {
		msg     MsgSetWithdrawAddress
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
		{testMsgs[2], false, "missing withdrawal address"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgSetWithdrawAddressGetSignBytes tests GetSignBytes for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressGetSignBytes(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testProvider, testWithdrawAddr)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgSetWithdrawAddress","value":{"provider":"faa1w3jhxapdwpex7anfv3jhynrxe9z","withdraw_address":"faa1w3jhxapdwa5hg6rywfshwctv94skgerjv4ehxdueap6"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgSetWithdrawAddressGetSigners tests GetSigners for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressGetSigners(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testProvider, testWithdrawAddr)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgDisableServiceRoute tests Route for MsgDisableService
func TestMsgDisableServiceRoute(t *testing.T) {
	msg := NewMsgDisableService(testServiceName, testProvider)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgDisableServiceType tests Type for MsgDisableService
func TestMsgDisableServiceType(t *testing.T) {
	msg := NewMsgDisableService(testServiceName, testProvider)

	require.Equal(t, "disable_service", msg.Type())
}

// TestMsgDisableServiceValidation tests ValidateBasic for MsgDisableService
func TestMsgDisableServiceValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)

	testMsgs := []MsgDisableService{
		NewMsgDisableService(testServiceName, testProvider), // valid msg
		NewMsgDisableService(testServiceName, emptyAddress), // missing provider address
		NewMsgDisableService(invalidName, testProvider),     // service name contains illegal characters
		NewMsgDisableService(invalidLongName, testProvider), // too long service name
	}

	testCases := []struct {
		msg     MsgDisableService
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
		{testMsgs[2], false, "service name contains illegal characters"},
		{testMsgs[3], false, "too long service name"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgDisableServiceGetSignBytes tests GetSignBytes for MsgDisableService
func TestMsgDisableServiceGetSignBytes(t *testing.T) {
	msg := NewMsgDisableService(testServiceName, testProvider)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgDisableService","value":{"provider":"faa1w3jhxapdwpex7anfv3jhynrxe9z","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgDisableServiceGetSigners tests GetSigners for MsgDisableService
func TestMsgDisableServiceGetSigners(t *testing.T) {
	msg := NewMsgDisableService(testServiceName, testProvider)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgEnableServiceRoute tests Route for MsgEnableService
func TestMsgEnableServiceRoute(t *testing.T) {
	msg := NewMsgEnableService(testServiceName, testProvider, testAddedDeposit)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgEnableServiceType tests Type for MsgEnableService
func TestMsgEnableServiceType(t *testing.T) {
	msg := NewMsgEnableService(testServiceName, testProvider, testAddedDeposit)

	require.Equal(t, "enable_service", msg.Type())
}

// TestMsgEnableServiceValidation tests ValidateBasic for MsgEnableService
func TestMsgEnableServiceValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}
	emptyAddedDeposit := sdk.Coins{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidDenomDeposit := sdk.NewCoins(sdk.NewCoin("eth-min", sdk.NewInt(1000)))

	testMsgs := []MsgEnableService{
		NewMsgEnableService(testServiceName, testProvider, testAddedDeposit),    // valid msg
		NewMsgEnableService(testServiceName, testProvider, emptyAddedDeposit),   // empty deposit is allowed
		NewMsgEnableService(testServiceName, emptyAddress, testAddedDeposit),    // missing provider address
		NewMsgEnableService(invalidName, testProvider, testAddedDeposit),        // service name contains illegal characters
		NewMsgEnableService(invalidLongName, testProvider, testAddedDeposit),    // too long service name
		NewMsgEnableService(testServiceName, testProvider, invalidDenomDeposit), // invalid deposit denom
	}

	testCases := []struct {
		msg     MsgEnableService
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], true, ""},
		{testMsgs[2], false, "missing provider address"},
		{testMsgs[3], false, "service name contains illegal characters"},
		{testMsgs[4], false, "too long service name"},
		{testMsgs[5], false, "invalid deposit denom"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgEnableServiceGetSignBytes tests GetSignBytes for MsgEnableService
func TestMsgEnableServiceGetSignBytes(t *testing.T) {
	msg := NewMsgEnableService(testServiceName, testProvider, testAddedDeposit)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgEnableService","value":{"deposit":[{"amount":"100000000000000000000","denom":"iris-atto"}],"provider":"faa1w3jhxapdwpex7anfv3jhynrxe9z","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgEnableServiceGetSigners tests GetSigners for MsgEnableService
func TestMsgEnableServiceGetSigners(t *testing.T) {
	msg := NewMsgEnableService(testServiceName, testProvider, testAddedDeposit)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgRefundServiceDepositRoute tests Route for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositRoute(t *testing.T) {
	msg := NewMsgRefundServiceDeposit(testServiceName, testProvider)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgRefundServiceDepositType tests Type for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositType(t *testing.T) {
	msg := NewMsgRefundServiceDeposit(testServiceName, testProvider)

	require.Equal(t, "refund_service_deposit", msg.Type())
}

// TestMsgRefundServiceDepositValidation tests ValidateBasic for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)

	testMsgs := []MsgRefundServiceDeposit{
		NewMsgRefundServiceDeposit(testServiceName, testProvider), // valid msg
		NewMsgRefundServiceDeposit(testServiceName, emptyAddress), // missing provider address
		NewMsgRefundServiceDeposit(invalidName, testProvider),     // service name contains illegal characters
		NewMsgRefundServiceDeposit(invalidLongName, testProvider), // too long service name
	}

	testCases := []struct {
		msg     MsgRefundServiceDeposit
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
		{testMsgs[2], false, "service name contains illegal characters"},
		{testMsgs[3], false, "too long service name"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgRefundServiceDepositGetSignBytes tests GetSignBytes for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositGetSignBytes(t *testing.T) {
	msg := NewMsgRefundServiceDeposit(testServiceName, testProvider)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgRefundServiceDeposit","value":{"provider":"faa1w3jhxapdwpex7anfv3jhynrxe9z","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgRefundServiceDepositGetSigners tests GetSigners for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositGetSigners(t *testing.T) {
	msg := NewMsgRefundServiceDeposit(testServiceName, testProvider)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgRequestServiceRoute tests Route for MsgRequestService
func TestMsgRequestServiceRoute(t *testing.T) {
	msg := NewMsgRequestService(
		testServiceName, testProviders, testConsumer,
		testInput, testServiceFeeCap, testTimeout, false,
		true, testRepeatedFreq, testRepeatedTotal,
	)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgRequestServiceType tests Type for MsgRequestService
func TestMsgRequestServiceType(t *testing.T) {
	msg := NewMsgRequestService(
		testServiceName, testProviders, testConsumer,
		testInput, testServiceFeeCap, testTimeout, false,
		true, testRepeatedFreq, testRepeatedTotal,
	)

	require.Equal(t, "request_service", msg.Type())
}

// TestMsgRequestServiceValidation tests ValidateBasic for MsgRequestService
func TestMsgRequestServiceValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidDenomCoins := sdk.NewCoins(sdk.NewCoin("eth-min", sdk.NewInt(1000)))

	invalidDuplicateProviders := []sdk.AccAddress{testProvider, testProvider}
	invalidInput := "iris-usdt"
	invalidTimeout := int64(0)
	invalidLessRepeatedFreq := uint64(testTimeout) - 10
	invalidRepeatedTotal1 := int64(-2)
	invalidRepeatedTotal2 := int64(0)

	testMsgs := []MsgRequestService{
		NewMsgRequestService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, false, true, testRepeatedFreq, testRepeatedTotal,
		), // valid msg
		NewMsgRequestService(
			testServiceName, testProviders, emptyAddress, testInput, testServiceFeeCap,
			testTimeout, false, true, testRepeatedFreq, testRepeatedTotal,
		), // missing consumer address
		NewMsgRequestService(
			invalidName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, false, true, testRepeatedFreq, testRepeatedTotal,
		), // service name contains illegal characters
		NewMsgRequestService(
			invalidLongName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, false, true, testRepeatedFreq, testRepeatedTotal,
		), // too long service name
		NewMsgRequestService(
			testServiceName, testProviders, testConsumer, testInput, invalidDenomCoins,
			testTimeout, false, true, testRepeatedFreq, testRepeatedTotal,
		), // invalid service fee denom
		NewMsgRequestService(
			testServiceName, nil, testConsumer, testInput, testServiceFeeCap,
			testTimeout, false, true, testRepeatedFreq, testRepeatedTotal,
		), // missing providers
		NewMsgRequestService(
			testServiceName, invalidDuplicateProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, false, true, testRepeatedFreq, testRepeatedTotal,
		), // duplicate providers
		NewMsgRequestService(
			testServiceName, testProviders, testConsumer, "", testServiceFeeCap,
			testTimeout, false, true, testRepeatedFreq, testRepeatedTotal,
		), // missing input
		NewMsgRequestService(
			testServiceName, testProviders, testConsumer, invalidInput, testServiceFeeCap,
			testTimeout, false, true, testRepeatedFreq, testRepeatedTotal,
		), // invalid input
		NewMsgRequestService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			invalidTimeout, false, true, testRepeatedFreq, testRepeatedTotal,
		), // invalid timeout
		NewMsgRequestService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, false, true, invalidLessRepeatedFreq, testRepeatedTotal,
		), // invalid repeated frequency
		NewMsgRequestService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, false, true, testRepeatedFreq, invalidRepeatedTotal1,
		), // repeated total can not be less than -1
		NewMsgRequestService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, false, true, testRepeatedFreq, invalidRepeatedTotal2,
		), // repeated total can not be zero
		NewMsgRequestService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, false, true, uint64(0), testRepeatedTotal,
		), // frequency can be zero
		NewMsgRequestService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, false, false, invalidLessRepeatedFreq, invalidRepeatedTotal1,
		), // do not check the repeated frequency and total when not repeated
	}

	testCases := []struct {
		msg     MsgRequestService
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing consumer address"},
		{testMsgs[2], false, "service name contains illegal characters"},
		{testMsgs[3], false, "too long service name"},
		{testMsgs[4], false, "invalid service fee denom"},
		{testMsgs[5], false, "missing providers"},
		{testMsgs[6], false, "duplicate providers"},
		{testMsgs[7], false, "missing input"},
		{testMsgs[8], false, "invalid input"},
		{testMsgs[9], false, "invalid timeout"},
		{testMsgs[10], false, "invalid repeated frequency"},
		{testMsgs[11], false, "repeated total can not be less than -1"},
		{testMsgs[12], false, "repeated total can not be zero"},
		{testMsgs[13], true, "frequency can be zero"},
		{testMsgs[14], true, "do not check the repeated frequency and total when not repeated"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgRequestServiceGetSignBytes tests GetSignBytes for MsgRequestService
func TestMsgRequestServiceGetSignBytes(t *testing.T) {
	msg := NewMsgRequestService(
		testServiceName, testProviders, testConsumer,
		testInput, testServiceFeeCap, testTimeout, false,
		true, testRepeatedFreq, testRepeatedTotal,
	)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgRequestService","value":{"consumer":"faa1w3jhxapdvdhkuum4d4jhyl0qvse","input":"{\"pair\":\"iris-usdt\"}","providers":["faa1w3jhxapdwpex7anfv3jhynrxe9z"],"repeated":true,"repeated_frequency":"120","repeated_total":"100","service_fee_cap":[{"amount":"100000000000000000000","denom":"iris-atto"}],"service_name":"test-service","super_mode":false,"timeout":"100"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgRequestServiceGetSigners tests GetSigners for MsgRequestService
func TestMsgRequestServiceGetSigners(t *testing.T) {
	msg := NewMsgRequestService(
		testServiceName, testProviders, testConsumer,
		testInput, testServiceFeeCap, testTimeout,
		false, true, testRepeatedFreq, testRepeatedTotal,
	)
	res := msg.GetSigners()

	expected := "[746573742D636F6E73756D6572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgRespondServiceRoute tests Route for MsgRespondService
func TestMsgRespondServiceRoute(t *testing.T) {
	msg := NewMsgRespondService(testRequestID, testProvider, testResult, testOutput)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgRespondServiceType tests Type for MsgRespondService
func TestMsgRespondServiceType(t *testing.T) {
	msg := NewMsgRespondService(testRequestID, testProvider, testResult, testOutput)

	require.Equal(t, "respond_service", msg.Type())
}

// TestMsgRespondServiceValidation tests ValidateBasic for MsgRespondService
func TestMsgRespondServiceValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidRequestID := "invalidRequestID"
	invalidOutput := "invalidOutput"

	validResult400 := `{"code":400,"message":"invalid parameters"}`
	invalidResult := "invalidResult"
	invalidResultCode := `{"code":100,"message":""}`
	invalidResultNoCode := `{"message":""}`
	invalidResultNoMsg := `{"code":200}`

	testMsgs := []MsgRespondService{
		NewMsgRespondService(testRequestID, testProvider, testResult, testOutput),     // valid msg
		NewMsgRespondService(testRequestID, testProvider, validResult400, ""),         // valid msg
		NewMsgRespondService(testRequestID, emptyAddress, testResult, testOutput),     // missing provider address
		NewMsgRespondService(invalidRequestID, testProvider, testResult, testOutput),  // invalid request ID
		NewMsgRespondService(testRequestID, testProvider, "", testOutput),             // missing result
		NewMsgRespondService(testRequestID, testProvider, invalidResult, ""),          // invalid result
		NewMsgRespondService(testRequestID, testProvider, invalidResultCode, ""),      // invalid result code
		NewMsgRespondService(testRequestID, testProvider, invalidResultNoCode, ""),    // missing result code
		NewMsgRespondService(testRequestID, testProvider, invalidResultNoMsg, ""),     // missing result message
		NewMsgRespondService(testRequestID, testProvider, testResult, ""),             // output should be provided when the result code is 200
		NewMsgRespondService(testRequestID, testProvider, testResult, invalidOutput),  // invalid output
		NewMsgRespondService(testRequestID, testProvider, validResult400, testOutput), // output should not be provided when the result code is not 200
	}

	testCases := []struct {
		msg     MsgRespondService
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], true, ""},
		{testMsgs[2], false, "missing provider address"},
		{testMsgs[3], false, "invalid request ID"},
		{testMsgs[4], false, "missing result"},
		{testMsgs[5], false, "invalid result"},
		{testMsgs[6], false, "invalid result code"},
		{testMsgs[7], false, "missing result code"},
		{testMsgs[8], false, "missing result message"},
		{testMsgs[9], false, "output should be provided when the result code is 200"},
		{testMsgs[10], false, "invalid output"},
		{testMsgs[11], false, "output should not be provided when the result code is not 200"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgRespondServiceGetSignBytes tests GetSignBytes for MsgRespondService
func TestMsgRespondServiceGetSignBytes(t *testing.T) {
	msg := NewMsgRespondService(testRequestID, testProvider, testResult, testOutput)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgRespondService","value":{"output":"{\"last\":\"100\"}","provider":"faa1w3jhxapdwpex7anfv3jhynrxe9z","request_id":"3DB0FA99DCB058BC86041BADBD614D6839F8FA20E17CF8AD3BA14C3F1BF613BD0000000000000000000000000000000100000000000000010001","result":"{\"code\":200,\"message\":\"\"}"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgRespondServiceGetSigners tests GetSigners for MsgRespondService
func TestMsgRespondServiceGetSigners(t *testing.T) {
	msg := NewMsgRespondService(testRequestID, testProvider, testResult, testOutput)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgPauseRequestContextRoute tests Route for MsgPauseRequestContext
func TestMsgPauseRequestContextRoute(t *testing.T) {
	msg := NewMsgPauseRequestContext(testRequestContextID, testConsumer)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgPauseRequestContextType tests Type for MsgPauseRequestContext
func TestMsgPauseRequestContextType(t *testing.T) {
	msg := NewMsgPauseRequestContext(testRequestContextID, testConsumer)

	require.Equal(t, "pause_request_context", msg.Type())
}

// TestMsgPauseRequestContextValidation tests ValidateBasic for MsgPauseRequestContext
func TestMsgPauseRequestContextValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}
	invalidRequestContextID := []byte("invalid-request-context-id")

	testMsgs := []MsgPauseRequestContext{
		NewMsgPauseRequestContext(testRequestContextID, testConsumer),    // valid msg
		NewMsgPauseRequestContext(testRequestContextID, emptyAddress),    // missing consumer address
		NewMsgPauseRequestContext(invalidRequestContextID, testConsumer), // invalid request context ID
	}

	testCases := []struct {
		msg     MsgPauseRequestContext
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing consumer address"},
		{testMsgs[2], false, "invalid request context ID"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgPauseRequestContextGetSignBytes tests GetSignBytes for MsgPauseRequestContext
func TestMsgPauseRequestContextGetSignBytes(t *testing.T) {
	msg := NewMsgPauseRequestContext(testRequestContextID, testConsumer)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgPauseRequestContext","value":{"consumer":"faa1w3jhxapdvdhkuum4d4jhyl0qvse","request_context_id":"PbD6mdywWLyGBButvWFNaDn4+iDhfPitO6FMPxv2E70AAAAAAAAAAA=="}}`
	require.Equal(t, expected, string(res))
}

// TestMsgPauseRequestContextGetSigners tests GetSigners for MsgPauseRequestContext
func TestMsgPauseRequestContextGetSigners(t *testing.T) {
	msg := NewMsgPauseRequestContext(testRequestContextID, testConsumer)
	res := msg.GetSigners()

	expected := "[746573742D636F6E73756D6572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgStartRequestContextRoute tests Route for MsgStartRequestContext
func TestMsgStartRequestContextRoute(t *testing.T) {
	msg := NewMsgStartRequestContext(testRequestContextID, testConsumer)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgStartRequestContextType tests Type for MsgStartRequestContext
func TestMsgStartRequestContextType(t *testing.T) {
	msg := NewMsgStartRequestContext(testRequestContextID, testConsumer)

	require.Equal(t, "start_request_context", msg.Type())
}

// TestMsgStartRequestContextValidation tests ValidateBasic for MsgStartRequestContext
func TestMsgStartRequestContextValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}
	invalidRequestContextID := []byte("invalid-request-context-id")

	testMsgs := []MsgStartRequestContext{
		NewMsgStartRequestContext(testRequestContextID, testConsumer),    // valid msg
		NewMsgStartRequestContext(testRequestContextID, emptyAddress),    // missing consumer address
		NewMsgStartRequestContext(invalidRequestContextID, testConsumer), // invalid request context ID
	}

	testCases := []struct {
		msg     MsgStartRequestContext
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing consumer address"},
		{testMsgs[2], false, "invalid request context ID"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgStartRequestContextGetSignBytes tests GetSignBytes for MsgStartRequestContext
func TestMsgStartRequestContextGetSignBytes(t *testing.T) {
	msg := NewMsgStartRequestContext(testRequestContextID, testConsumer)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgStartRequestContext","value":{"consumer":"faa1w3jhxapdvdhkuum4d4jhyl0qvse","request_context_id":"PbD6mdywWLyGBButvWFNaDn4+iDhfPitO6FMPxv2E70AAAAAAAAAAA=="}}`
	require.Equal(t, expected, string(res))
}

// TestMsgStartRequestContextGetSigners tests GetSigners for MsgStartRequestContext
func TestMsgStartRequestContextGetSigners(t *testing.T) {
	msg := NewMsgStartRequestContext(testRequestContextID, testConsumer)
	res := msg.GetSigners()

	expected := "[746573742D636F6E73756D6572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgKillRequestContextRoute tests Route for MsgKillRequestContext
func TestMsgKillRequestContextRoute(t *testing.T) {
	msg := NewMsgKillRequestContext(testRequestContextID, testConsumer)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgKillRequestContextType tests Type for MsgKillRequestContext
func TestMsgKillRequestContextType(t *testing.T) {
	msg := NewMsgKillRequestContext(testRequestContextID, testConsumer)

	require.Equal(t, "kill_request_context", msg.Type())
}

// TestMsgKillRequestContextValidation tests ValidateBasic for MsgKillRequestContext
func TestMsgKillRequestContextValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}
	invalidRequestContextID := []byte("invalid-request-context-id")

	testMsgs := []MsgKillRequestContext{
		NewMsgKillRequestContext(testRequestContextID, testConsumer),    // valid msg
		NewMsgKillRequestContext(testRequestContextID, emptyAddress),    // missing consumer address
		NewMsgKillRequestContext(invalidRequestContextID, testConsumer), // invalid request context ID
	}

	testCases := []struct {
		msg     MsgKillRequestContext
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing consumer address"},
		{testMsgs[2], false, "invalid request context ID"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgKillRequestContextGetSignBytes tests GetSignBytes for MsgKillRequestContext
func TestMsgKillRequestContextGetSignBytes(t *testing.T) {
	msg := NewMsgKillRequestContext(testRequestContextID, testConsumer)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgKillRequestContext","value":{"consumer":"faa1w3jhxapdvdhkuum4d4jhyl0qvse","request_context_id":"PbD6mdywWLyGBButvWFNaDn4+iDhfPitO6FMPxv2E70AAAAAAAAAAA=="}}`
	require.Equal(t, expected, string(res))
}

// TestMsgKillRequestContextGetSigners tests GetSigners for MsgKillRequestContext
func TestMsgKillRequestContextGetSigners(t *testing.T) {
	msg := NewMsgKillRequestContext(testRequestContextID, testConsumer)
	res := msg.GetSigners()

	expected := "[746573742D636F6E73756D6572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgUpdateRequestContextRoute tests Route for MsgUpdateRequestContext
func TestMsgUpdateRequestContextRoute(t *testing.T) {
	msg := NewMsgUpdateRequestContext(testRequestContextID, nil, nil, 0, 0, 0, testConsumer)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgUpdateRequestContextType tests Type for MsgUpdateRequestContext
func TestMsgUpdateRequestContextType(t *testing.T) {
	msg := NewMsgUpdateRequestContext(testRequestContextID, nil, nil, 0, 0, 0, testConsumer)

	require.Equal(t, "update_request_context", msg.Type())
}

// TestMsgUpdateRequestContextValidation tests ValidateBasic for MsgUpdateRequestContext
func TestMsgUpdateRequestContextValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidRequestContextID := []byte("invalid-request-context-id")
	invalidDuplicateProviders := []sdk.AccAddress{testProvider, testProvider}
	invalidTimeout := int64(-1)
	invalidLessRepeatedFreq := uint64(testTimeout) - 10
	invalidRepeatedTotal := int64(-2)
	invalidDenomCoins := sdk.NewCoins(sdk.NewCoin("eth-min", sdk.NewInt(1000)))

	testMsgs := []MsgUpdateRequestContext{
		NewMsgUpdateRequestContext(testRequestContextID, testProviders, testServiceFeeCap, testTimeout, testRepeatedFreq, testRepeatedTotal, testConsumer), // valid msg
		NewMsgUpdateRequestContext(testRequestContextID, nil, nil, 0, 0, 0, testConsumer),                                                                  // allow all not to be updated
		NewMsgUpdateRequestContext(testRequestContextID, nil, nil, 0, 0, 0, emptyAddress),                                                                  // missing consumer address
		NewMsgUpdateRequestContext(invalidRequestContextID, nil, nil, 0, 0, 0, testConsumer),                                                               // invalid request context ID
		NewMsgUpdateRequestContext(testRequestContextID, invalidDuplicateProviders, nil, 0, 0, 0, testConsumer),                                            // duplicate providers
		NewMsgUpdateRequestContext(testRequestContextID, nil, nil, invalidTimeout, 0, 0, testConsumer),                                                     // invalid timeout
		NewMsgUpdateRequestContext(invalidRequestContextID, nil, nil, testTimeout, invalidLessRepeatedFreq, 0, testConsumer),                               // invalid repeated frequency
		NewMsgUpdateRequestContext(testRequestContextID, nil, nil, 0, 0, invalidRepeatedTotal, testConsumer),                                               // invalid repeated total
		NewMsgUpdateRequestContext(testRequestContextID, nil, invalidDenomCoins, 0, 0, 0, testConsumer),                                                    // invalid service fee denom
	}

	testCases := []struct {
		msg     MsgUpdateRequestContext
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], true, ""},
		{testMsgs[2], false, "missing consumer address"},
		{testMsgs[3], false, "invalid request context ID"},
		{testMsgs[4], false, "duplicate providers"},
		{testMsgs[5], false, "invalid timeout"},
		{testMsgs[6], false, "invalid repeated frequency"},
		{testMsgs[7], false, "invalid repeated total"},
		{testMsgs[8], false, "invalid service fee denom"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgUpdateRequestContextGetSignBytes tests GetSignBytes for MsgUpdateRequestContext
func TestMsgUpdateRequestContextGetSignBytes(t *testing.T) {
	msg := NewMsgUpdateRequestContext(testRequestContextID, testProviders, testServiceFeeCap, testTimeout, testRepeatedFreq, testRepeatedTotal, testConsumer)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgUpdateRequestContext","value":{"consumer":"faa1w3jhxapdvdhkuum4d4jhyl0qvse","providers":["faa1w3jhxapdwpex7anfv3jhynrxe9z"],"repeated_frequency":"120","repeated_total":"100","request_context_id":"PbD6mdywWLyGBButvWFNaDn4+iDhfPitO6FMPxv2E70AAAAAAAAAAA==","service_fee_cap":[{"amount":"100000000000000000000","denom":"iris-atto"}],"timeout":"100"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgUpdateRequestContextGetSigners tests GetSigners for MsgUpdateRequestContext
func TestMsgUpdateRequestContextGetSigners(t *testing.T) {
	msg := NewMsgUpdateRequestContext(testRequestContextID, testProviders, testServiceFeeCap, testTimeout, testRepeatedFreq, testRepeatedTotal, testConsumer)
	res := msg.GetSigners()

	expected := "[746573742D636F6E73756D6572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgWithdrawEarnedFeesRoute tests Route for MsgWithdrawEarnedFees
func TestMsgWithdrawEarnedFeesRoute(t *testing.T) {
	msg := NewMsgWithdrawEarnedFees(testProvider)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgWithdrawEarnedFeesType tests Type for MsgWithdrawEarnedFees
func TestMsgWithdrawEarnedFeesType(t *testing.T) {
	msg := NewMsgWithdrawEarnedFees(testProvider)

	require.Equal(t, "withdraw_earned_fees", msg.Type())
}

// TestMsgWithdrawEarnedFeesValidation tests ValidateBasic for MsgWithdrawEarnedFees
func TestMsgWithdrawEarnedFeesValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	testMsgs := []MsgWithdrawEarnedFees{
		NewMsgWithdrawEarnedFees(testProvider), // valid msg
		NewMsgWithdrawEarnedFees(emptyAddress), // missing provider address
	}

	testCases := []struct {
		msg     MsgWithdrawEarnedFees
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgWithdrawEarnedFeesGetSignBytes tests GetSignBytes for MsgWithdrawEarnedFees
func TestMsgWithdrawEarnedFeesGetSignBytes(t *testing.T) {
	msg := NewMsgWithdrawEarnedFees(testProvider)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgWithdrawEarnedFees","value":{"provider":"faa1w3jhxapdwpex7anfv3jhynrxe9z"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgWithdrawEarnedFeesGetSigners tests GetSigners for MsgWithdrawEarnedFees
func TestMsgWithdrawEarnedFeesGetSigners(t *testing.T) {
	msg := NewMsgWithdrawEarnedFees(testProvider)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgWithdrawTaxRoute tests Route for MsgWithdrawTax
func TestMsgWithdrawTaxRoute(t *testing.T) {
	msg := NewMsgWithdrawTax(testTrustee, testTrustee, testTaxWithdrawalAmt)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgWithdrawTaxType tests Type for MsgWithdrawTax
func TestMsgWithdrawTaxType(t *testing.T) {
	msg := NewMsgWithdrawTax(testTrustee, testTrustee, testTaxWithdrawalAmt)

	require.Equal(t, "withdraw_tax", msg.Type())
}

// TestMsgWithdrawTaxValidation tests ValidateBasic for MsgWithdrawTax
func TestMsgWithdrawTaxValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}
	invalidDenomCoins := sdk.NewCoins(sdk.NewCoin("eth-min", sdk.NewInt(1000)))

	testMsgs := []MsgWithdrawTax{
		NewMsgWithdrawTax(testTrustee, testTrustee, testTaxWithdrawalAmt),  // valid msg
		NewMsgWithdrawTax(emptyAddress, testTrustee, testTaxWithdrawalAmt), // missing trustee address
		NewMsgWithdrawTax(testTrustee, emptyAddress, testTaxWithdrawalAmt), // missing destination address
		NewMsgWithdrawTax(testTrustee, testTrustee, invalidDenomCoins),     // invalid withdrawal denom
	}

	testCases := []struct {
		msg     MsgWithdrawTax
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing trustee address"},
		{testMsgs[2], false, "missing destination address"},
		{testMsgs[3], false, "invalid withdrawal denom"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgWithdrawTaxGetSignBytes tests GetSignBytes for MsgWithdrawTax
func TestMsgWithdrawTaxGetSignBytes(t *testing.T) {
	msg := NewMsgWithdrawTax(testTrustee, testTrustee, testTaxWithdrawalAmt)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgWithdrawTax","value":{"amount":[{"amount":"1000000000000000000000","denom":"iris-atto"}],"dest_address":"faa1w3jhxapdw3e82um5v4jscm70uh","trustee":"faa1w3jhxapdw3e82um5v4jscm70uh"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgWithdrawTaxGetSigners tests GetSigners for MsgWithdrawTax
func TestMsgWithdrawTaxGetSigners(t *testing.T) {
	msg := NewMsgWithdrawTax(testTrustee, testTrustee, testTaxWithdrawalAmt)
	res := msg.GetSigners()

	expected := "[746573742D74727573746565]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
