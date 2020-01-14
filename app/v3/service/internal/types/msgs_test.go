package types

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/irisnet/irishub/types"
)

var (
	testCoin1, _ = sdk.IrisCoinType.ConvertToMinDenomCoin("1000iris")
	testCoin2, _ = sdk.IrisCoinType.ConvertToMinDenomCoin("100iris")

	testServiceName = "test-service"
	testServiceDesc = "test-service-desc"
	testServiceTags = []string{"tag1", "tag2"}
	testAuthor      = sdk.AccAddress([]byte("test-author"))
	testAuthorDesc  = "test-author-desc"
	testSchemas     = `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`

	testProvider     = sdk.AccAddress([]byte("test-provider"))
	testDeposit      = sdk.NewCoins(testCoin1)
	testPricing      = `{"price":"100iris-atto"}`
	testWithdrawAddr = sdk.AccAddress([]byte("test-withdrawal-address"))
	testAddedDeposit = sdk.NewCoins(testCoin2)
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

	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidLongDesc := strings.Repeat("d", MaxDescriptionLength+1)
	invalidMoreTags := strings.Split("t1,t2,t3,t4,t5,t6,t7,t8,t9,t10,t11", ",")
	invalidLongTags := []string{strings.Repeat("t", MaxTagLength+1)}

	invalidSchemas := `{"input":"nonobject","output":"nonobject","error":"nonobject"}`
	invalidSchemasNoInput := `{"output":{"type":"object"},"error":{"type":"object"}}`
	invalidSchemasNoOutput := `{"input":{"type":"object"},"error":{"type":"object"}}`
	invalidSchemasNoError := `{"input":{"type":"object"},"output":{"type":"object"}}`

	testMsgs := []MsgDefineService{
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas),            // valid msg
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, emptyAddress, testAuthorDesc, testSchemas),          // missing author address
		NewMsgDefineService("service/name", testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas),             // service name contains illegal characters
		NewMsgDefineService(invalidLongName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas),            // too long service name
		NewMsgDefineService(testServiceName, invalidLongDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas),            // too long service description
		NewMsgDefineService(testServiceName, testServiceDesc, invalidMoreTags, testAuthor, testAuthorDesc, testSchemas),            // too many tags
		NewMsgDefineService(testServiceName, testServiceDesc, invalidLongTags, testAuthor, testAuthorDesc, testSchemas),            // too long tag
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, invalidLongDesc, testSchemas),           // too long author description
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, invalidSchemas),         // invalid schemas
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, invalidSchemasNoInput),  // missing input schema
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, invalidSchemasNoOutput), // missing output schema
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, invalidSchemasNoError),  // missing error schema                              // not possitive coin
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
		{testMsgs[7], false, "too long author description"},
		{testMsgs[8], false, "invalid schemas"},
		{testMsgs[9], false, "missing input schema"},
		{testMsgs[10], false, "missing output schema"},
		{testMsgs[11], false, "missing error schema"},
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

	expected := `{"type":"irishub/service/MsgDefineService","value":{"author":"faa1w3jhxapdv96hg6r0wgux6kuw","author_description":"test-author-desc","description":"test-service-desc","name":"test-service","schemas":"{\"input\":{\"type\":\"object\"},\"output\":{\"type\":\"object\"},\"error\":{\"type\":\"object\"}}","tags":["tag1","tag2"]}}`
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
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testWithdrawAddr)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgBindServiceType tests Type for MsgBindService
func TestMsgBindServiceType(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testWithdrawAddr)

	require.Equal(t, "bind_service", msg.Type())
}

// TestMsgBindServiceValidation tests ValidateBasic for MsgBindService
func TestMsgBindServiceValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidDeposit := sdk.Coins{}
	invalidDenomDeposit := sdk.NewCoins(sdk.NewCoin("iris", sdk.NewInt(1000)))
	invalidPricing := `{"price":"100iris-atto","other":"notallowedfield"}`
	invalidDenomPricing := `{"price":"100iris"}`

	testMsgs := []MsgBindService{
		NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testWithdrawAddr),         // valid msg
		NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, emptyAddress),             // empty withdrawal address is allowed
		NewMsgBindService(testServiceName, emptyAddress, testDeposit, testPricing, testWithdrawAddr),         // missing provider address
		NewMsgBindService(invalidName, testProvider, testDeposit, testPricing, testWithdrawAddr),             // service name contains illegal characters
		NewMsgBindService(invalidLongName, testProvider, testDeposit, testPricing, testWithdrawAddr),         // too long service name
		NewMsgBindService(testServiceName, testProvider, invalidDeposit, testPricing, testWithdrawAddr),      // invalid deposit
		NewMsgBindService(testServiceName, testProvider, invalidDenomDeposit, testPricing, testWithdrawAddr), // invalid deposit denom
		NewMsgBindService(testServiceName, testProvider, testDeposit, "", testWithdrawAddr),                  // missing pricing
		NewMsgBindService(testServiceName, testProvider, testDeposit, invalidPricing, testWithdrawAddr),      // invalid Pricing JSON Schema instance
		NewMsgBindService(testServiceName, testProvider, testDeposit, invalidDenomPricing, testWithdrawAddr), // invalid pricing denom
	}

	testCases := []struct {
		msg     MsgBindService
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], true, ""},
		{testMsgs[2], false, "missing provider address"},
		{testMsgs[3], false, "service name contains illegal characters"},
		{testMsgs[4], false, "too long service name"},
		{testMsgs[5], false, "invalid deposit"},
		{testMsgs[6], false, "invalid deposit denom"},
		{testMsgs[7], false, "missing pricing"},
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

// TestMsgBindServiceGetSignBytes tests GetSignBytes for MsgBindService
func TestMsgBindServiceGetSignBytes(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testWithdrawAddr)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgBindService","value":{"deposit":[{"amount":"1000000000000000000000","denom":"iris-atto"}],"pricing":"{\"price\":\"100iris-atto\"}","provider":"faa1w3jhxapdwpex7anfv3jhynrxe9z","service_name":"test-service","withdraw_address":"faa1w3jhxapdwa5hg6rywfshwctv94skgerjv4ehxdueap6"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgBindServiceGetSigners tests GetSigners for MsgBindService
func TestMsgBindServiceGetSigners(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testWithdrawAddr)
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
	invalidDenomDeposit := sdk.NewCoins(sdk.NewCoin("iris", sdk.NewInt(1000)))
	invalidPricing := `{"price":"100iris-atto","other":"notallowedfield"}`
	invalidDenomPricing := `{"price":"100iris"}`

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
	msg := NewMsgSetWithdrawAddress(testServiceName, testProvider, testWithdrawAddr)

	require.Equal(t, MsgRoute, msg.Route())
}

// TestMsgSetWithdrawAddressType tests Type for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressType(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testServiceName, testProvider, testWithdrawAddr)

	require.Equal(t, "set_withdraw_address", msg.Type())
}

// TestMsgSetWithdrawAddressValidation tests ValidateBasic for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)

	testMsgs := []MsgSetWithdrawAddress{
		NewMsgSetWithdrawAddress(testServiceName, testProvider, testWithdrawAddr), // valid msg
		NewMsgSetWithdrawAddress(testServiceName, emptyAddress, testWithdrawAddr), // missing provider address
		NewMsgSetWithdrawAddress(invalidName, testProvider, testWithdrawAddr),     // service name contains illegal characters
		NewMsgSetWithdrawAddress(invalidLongName, testProvider, testWithdrawAddr), // too long service name
		NewMsgSetWithdrawAddress(testServiceName, testProvider, emptyAddress),     // missing withdrawal address
	}

	testCases := []struct {
		msg     MsgSetWithdrawAddress
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
		{testMsgs[2], false, "service name contains illegal characters"},
		{testMsgs[3], false, "too long service name"},
		{testMsgs[4], false, "missing withdrawal address"},
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
	msg := NewMsgSetWithdrawAddress(testServiceName, testProvider, testWithdrawAddr)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/service/MsgSetWithdrawAddress","value":{"provider":"faa1w3jhxapdwpex7anfv3jhynrxe9z","service_name":"test-service","withdraw_address":"faa1w3jhxapdwa5hg6rywfshwctv94skgerjv4ehxdueap6"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgSetWithdrawAddressGetSigners tests GetSigners for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressGetSigners(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testServiceName, testProvider, testWithdrawAddr)
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
	invalidDenomDeposit := sdk.NewCoins(sdk.NewCoin("iris", sdk.NewInt(1000)))

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
