package types

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	emptyAddress = ""

	testCoin1 = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000))
	testCoin2 = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))

	testServiceName = "test-service"
	testServiceDesc = "test-service-desc"
	testServiceTags = []string{"tag1", "tag2"}
	testAuthor      = sdk.AccAddress(tmhash.SumTruncated([]byte("test-author"))).String()
	testAuthorDesc  = "test-author-desc"
	testSchemas     = `{"input":{"type":"object"},"output":{"type":"object"}}`

	testOwner        = sdk.AccAddress(tmhash.SumTruncated([]byte("test-owner"))).String()
	testProvider     = sdk.AccAddress(tmhash.SumTruncated([]byte("test-provider"))).String()
	testDeposit      = sdk.NewCoins(testCoin1)
	testPricing      = `{"price":"1stake"}`
	testQoS          = uint64(50)
	testOptions      = "{}"
	testWithdrawAddr = sdk.AccAddress(tmhash.SumTruncated([]byte("test-withdrawal-address"))).String()
	testAddedDeposit = sdk.NewCoins(testCoin2)

	testConsumer      = sdk.AccAddress(tmhash.SumTruncated([]byte("test-consumer"))).String()
	testProviders     = []string{testProvider}
	testInput         = `{"header":{},"body":{"pair":"iris-usdt"}}`
	testServiceFeeCap = sdk.NewCoins(testCoin2)
	testTimeout       = int64(100)
	testRepeatedFreq  = uint64(120)
	testRepeatedTotal = int64(100)

	testResult = `{"code":200,"message":""}`
	testOutput = `{"header":{},"body":{"last":"100"}}`

	rawTestRequestContextID = GenerateRequestContextID(tmhash.Sum([]byte("test-request-context-id")), 0)
	testRequestContextID    = rawTestRequestContextID.String()
	testRequestID           = GenerateRequestID(rawTestRequestContextID, 1, 1, 1).String()

	// -------------

	invalidName = "invalid/service/name"
)

// TestMsgDefineServiceRoute tests Route for MsgDefineService
func TestMsgDefineServiceRoute(t *testing.T) {
	msg := NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)
	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgDefineServiceType tests Type for MsgDefineService
func TestMsgDefineServiceType(t *testing.T) {
	msg := NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)
	require.Equal(t, "define_service", msg.Type())
}

// TestMsgDefineServiceValidation tests ValidateBasic for MsgDefineService
func TestMsgDefineServiceValidation(t *testing.T) {
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidLongDesc := strings.Repeat("d", MaxDescriptionLength+1)
	invalidMoreTags := strings.Split("t1,t2,t3,t4,t5,t6,t7,t8,t9,t10,t11", ",")
	invalidLongTags := []string{strings.Repeat("t", MaxTagLength+1)}
	invalidEmptyTags := []string{"t1", ""}
	invalidDuplicateTags := []string{"t1", "t1"}

	invalidSchemas := `{"input":"nonobject","output":"nonobject"}`
	invalidSchemasNoInput := `{"output":{"type":"object"}}`
	invalidSchemasNoOutput := `{"input":{"type":"object"}}`

	testMsgs := []*MsgDefineService{
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
		msg     *MsgDefineService
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

	expected := `{"type":"irismod/service/MsgDefineService","value":{"author":"cosmos1tkelht5u9ywfaww3fehvdtvncwwhf539w406rh","author_description":"test-author-desc","description":"test-service-desc","name":"test-service","schemas":"{\"input\":{\"type\":\"object\"},\"output\":{\"type\":\"object\"}}","tags":["tag1","tag2"]}}`
	require.Equal(t, expected, string(res))
}

// TestMsgDefineServiceGetSigners tests GetSigners for MsgDefineService
func TestMsgDefineServiceGetSigners(t *testing.T) {
	msg := NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)
	res := msg.GetSigners()

	expected := "[5DB3FBAE9C291C9EB9D14E6EC6AD93C39D74D225]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgBindServiceRoute tests Route for MsgBindService
func TestMsgBindServiceRoute(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testQoS, testOptions, testOwner)
	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgBindServiceType tests Type for MsgBindService
func TestMsgBindServiceType(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testQoS, testOptions, testOwner)
	require.Equal(t, "bind_service", msg.Type())
}

// TestMsgBindServiceValidation tests ValidateBasic for MsgBindService
func TestMsgBindServiceValidation(t *testing.T) {
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidDeposit := sdk.Coins{sdk.Coin{Denom: "stake", Amount: sdk.NewInt(-1)}}
	invalidQoS := uint64(0)

	invalidPricing := `{"price":"1stake","other":"notallowedfield"}`
	invalidDenomPricing := `{"price":"5invalid#denom"}`
	invalidPromotionTimePricing := `{"price":"1stake","promotions_by_time":[{"start_time":"2018-10-10T13:30:30","end_time":"2019-10-10T13:30:30Z","discount":"0.8"}]}`
	invalidPromotionVolPricing := `{"price":"1stake","promotions_by_volume":[{"volume":0,"discount":"0.7"}]}`
	invalidOptions := ""

	testMsgs := []*MsgBindService{
		NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testQoS, testOptions, testOwner),                 // valid msg
		NewMsgBindService(testServiceName, emptyAddress, testDeposit, testPricing, testQoS, testOptions, testOwner),                 // missing provider address
		NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testQoS, testOptions, emptyAddress),              // missing owner address
		NewMsgBindService(invalidName, testProvider, testDeposit, testPricing, testQoS, testOptions, testOwner),                     // service name contains illegal characters
		NewMsgBindService(invalidLongName, testProvider, testDeposit, testPricing, testQoS, testOptions, testOwner),                 // too long service name
		NewMsgBindService(testServiceName, testProvider, invalidDeposit, testPricing, testQoS, testOptions, testOwner),              // invalid deposit
		NewMsgBindService(testServiceName, testProvider, testDeposit, "", testQoS, testOptions, testOwner),                          // missing pricing
		NewMsgBindService(testServiceName, testProvider, testDeposit, invalidPricing, testQoS, testOptions, testOwner),              // invalid Pricing JSON Schema instance
		NewMsgBindService(testServiceName, testProvider, testDeposit, invalidDenomPricing, testQoS, testOptions, testOwner),         // invalid pricing denom
		NewMsgBindService(testServiceName, testProvider, testDeposit, invalidPromotionTimePricing, testQoS, testOptions, testOwner), // invalid promotion time lack of time zone
		NewMsgBindService(testServiceName, testProvider, testDeposit, invalidPromotionVolPricing, testQoS, testOptions, testOwner),  // invalid promotion volume
		NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, invalidQoS, testOptions, testOwner),              // invalid qos
		NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testQoS, invalidOptions, testOwner),              // invalid options
	}

	testCases := []struct {
		msg     *MsgBindService
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
		{testMsgs[2], false, "missing owner address"},
		{testMsgs[3], false, "service name contains illegal characters"},
		{testMsgs[4], false, "too long service name"},
		{testMsgs[5], false, "invalid deposit"},
		{testMsgs[6], false, "missing pricing"},
		{testMsgs[7], false, "invalid Pricing JSON Schema instance"},
		{testMsgs[8], false, "invalid pricing denom"},
		{testMsgs[9], false, "invalid promotion time lack of time zone"},
		{testMsgs[10], false, "invalid promotion volume"},
		{testMsgs[11], false, "invalid qos"},
		{testMsgs[12], false, "invalid options"},
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
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testQoS, testOptions, testOwner)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgBindService","value":{"deposit":[{"amount":"10000","denom":"stake"}],"options":"{}","owner":"cosmos1dtrajkx72qwf8gesp2z4rjz6p7klycmmh2lnvn","pricing":"{\"price\":\"1stake\"}","provider":"cosmos10yajkmnug9d2nvccgs7ul3xkur0l9as348gmwe","qos":"50","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgBindServiceGetSigners tests GetSigners for MsgBindService
func TestMsgBindServiceGetSigners(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testQoS, "", testOwner)
	res := msg.GetSigners()
	require.Equal(t, "[6AC7D958DE501C93A3300A8551C85A0FADF2637B]", fmt.Sprintf("%v", res))
}

// TestMsgUpdateServiceBindingRoute tests Route for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingRoute(t *testing.T) {
	msg := NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "", 0, testOptions, testOwner)
	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgUpdateServiceBindingType tests Type for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingType(t *testing.T) {
	msg := NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "", 0, testOptions, testOwner)
	require.Equal(t, "update_service_binding", msg.Type())
}

// TestMsgUpdateServiceBindingValidation tests ValidateBasic for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingValidation(t *testing.T) {
	emptyAddedDeposit := sdk.Coins{}

	invalidLongName := strings.Repeat("s", MaxNameLength+1)

	invalidPricing := `{"price":"1stake","other":"notallowedfield"}`
	invalidDenomPricing := `{"price":"1invaliddenom*"}`
	invalidPromotionTimePricing := `{"price":"1stake","promotions_by_time":[{"start_time":"2018-10-10T13:30:30","end_time":"2019-10-10T13:30:30Z","discount":"0.8"}]}`
	invalidPromotionVolPricing := `{"price":"1stake","promotions_by_volume":[{"volume":0,"discount":"0.7"}]}`
	invalidOptions := ""

	testMsgs := []*MsgUpdateServiceBinding{
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, testPricing, testQoS, testOptions, testOwner),                   // valid msg
		NewMsgUpdateServiceBinding(testServiceName, testProvider, emptyAddedDeposit, testPricing, testQoS, testOptions, testOwner),                  // empty deposit is allowed
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "", testQoS, testOptions, testOwner),                            // empty pricing is allowed
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, testPricing, 0, testOptions, testOwner),                         // 0 is allowed for qos
		NewMsgUpdateServiceBinding(testServiceName, testProvider, emptyAddedDeposit, "", 0, testOptions, testOwner),                                 // deposit, pricing and qos can be empty at the same time
		NewMsgUpdateServiceBinding(testServiceName, emptyAddress, testAddedDeposit, testPricing, testQoS, testOptions, testOwner),                   // missing provider address
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, testPricing, testQoS, testOptions, emptyAddress),                // missing owner address
		NewMsgUpdateServiceBinding(invalidName, testProvider, testAddedDeposit, testPricing, testQoS, testOptions, testOwner),                       // service name contains illegal characters
		NewMsgUpdateServiceBinding(invalidLongName, testProvider, testAddedDeposit, testPricing, testQoS, testOptions, testOwner),                   // too long service name
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, invalidPricing, testQoS, testOptions, testOwner),                // invalid Pricing JSON Schema instance
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, invalidDenomPricing, testQoS, testOptions, testOwner),           // invalid pricing denom
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, invalidPromotionTimePricing, testQoS, testOptions, testOwner),   // invalid promotion time lack of time zone
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, invalidPromotionVolPricing, testQoS, testOptions, testOwner),    // invalid promotion volume
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, invalidPromotionVolPricing, testQoS, invalidOptions, testOwner), // invalid options
	}

	testCases := []struct {
		msg     *MsgUpdateServiceBinding
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], true, ""},
		{testMsgs[2], true, ""},
		{testMsgs[3], true, ""},
		{testMsgs[4], true, ""},
		{testMsgs[5], false, "missing provider address"},
		{testMsgs[6], false, "missing owner address"},
		{testMsgs[7], false, "service name contains illegal characters"},
		{testMsgs[8], false, "too long service name"},
		{testMsgs[9], false, "invalid Pricing JSON Schema instance"},
		{testMsgs[10], false, "invalid pricing denom"},
		{testMsgs[11], false, "invalid promotion time lack of time zone"},
		{testMsgs[12], false, "invalid promotion volume"},
		{testMsgs[13], false, "invalid options"},
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
	msg := NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "{\"price\":\"1\"}", 1, "", testOwner)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgUpdateServiceBinding","value":{"deposit":[{"amount":"100","denom":"stake"}],"owner":"cosmos1dtrajkx72qwf8gesp2z4rjz6p7klycmmh2lnvn","pricing":"{\"price\":\"1\"}","provider":"cosmos10yajkmnug9d2nvccgs7ul3xkur0l9as348gmwe","qos":"1","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgUpdateServiceBindingGetSigners tests GetSigners for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingGetSigners(t *testing.T) {
	msg := NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "", 0, testOptions, testOwner)
	res := msg.GetSigners()
	require.Equal(t, "[6AC7D958DE501C93A3300A8551C85A0FADF2637B]", fmt.Sprintf("%v", res))
}

// TestMsgSetWithdrawAddressRoute tests Route for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressRoute(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testOwner, testWithdrawAddr)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgSetWithdrawAddressType tests Type for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressType(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testOwner, testWithdrawAddr)

	require.Equal(t, "set_withdraw_address", msg.Type())
}

// TestMsgSetWithdrawAddressValidation tests ValidateBasic for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressValidation(t *testing.T) {
	testMsgs := []*MsgSetWithdrawAddress{
		NewMsgSetWithdrawAddress(testOwner, testWithdrawAddr),    // valid msg
		NewMsgSetWithdrawAddress(emptyAddress, testWithdrawAddr), // missing owner address
		NewMsgSetWithdrawAddress(testOwner, emptyAddress),        // missing withdrawal address
	}

	testCases := []struct {
		msg     *MsgSetWithdrawAddress
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing owner address"},
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
	msg := NewMsgSetWithdrawAddress(testOwner, testWithdrawAddr)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgSetWithdrawAddress","value":{"owner":"cosmos1dtrajkx72qwf8gesp2z4rjz6p7klycmmh2lnvn","withdraw_address":"cosmos19zz7wpde40ed545384w0xsrgkydq4pvzdwrupn"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgSetWithdrawAddressGetSigners tests GetSigners for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressGetSigners(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testOwner, testWithdrawAddr)
	res := msg.GetSigners()
	require.Equal(t, "[6AC7D958DE501C93A3300A8551C85A0FADF2637B]", fmt.Sprintf("%v", res))
}

// TestMsgDisableServiceBindingRoute tests Route for MsgDisableServiceBinding
func TestMsgDisableServiceBindingRoute(t *testing.T) {
	msg := NewMsgDisableServiceBinding(testServiceName, testProvider, testOwner)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgDisableServiceBindingType tests Type for MsgDisableServiceBinding
func TestMsgDisableServiceBindingType(t *testing.T) {
	msg := NewMsgDisableServiceBinding(testServiceName, testProvider, testOwner)

	require.Equal(t, "disable_service_binding", msg.Type())
}

// TestMsgDisableServiceBindingValidation tests ValidateBasic for MsgDisableServiceBinding
func TestMsgDisableServiceBindingValidation(t *testing.T) {
	invalidLongName := strings.Repeat("s", MaxNameLength+1)

	testMsgs := []*MsgDisableServiceBinding{
		NewMsgDisableServiceBinding(testServiceName, testProvider, testOwner),    // valid msg
		NewMsgDisableServiceBinding(testServiceName, emptyAddress, testOwner),    // missing provider address
		NewMsgDisableServiceBinding(testServiceName, testProvider, emptyAddress), // missing owner address
		NewMsgDisableServiceBinding(invalidName, testProvider, testOwner),        // service name contains illegal characters
		NewMsgDisableServiceBinding(invalidLongName, testProvider, testOwner),    // too long service name
	}

	testCases := []struct {
		msg     *MsgDisableServiceBinding
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
		{testMsgs[2], false, "missing owner address"},
		{testMsgs[3], false, "service name contains illegal characters"},
		{testMsgs[4], false, "too long service name"},
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

// TestMsgDisableServiceBindingGetSignBytes tests GetSignBytes for MsgDisableServiceBinding
func TestMsgDisableServiceBindingGetSignBytes(t *testing.T) {
	msg := NewMsgDisableServiceBinding(testServiceName, testProvider, testOwner)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgDisableServiceBinding","value":{"owner":"cosmos1dtrajkx72qwf8gesp2z4rjz6p7klycmmh2lnvn","provider":"cosmos10yajkmnug9d2nvccgs7ul3xkur0l9as348gmwe","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgDisableServiceBindingGetSigners tests GetSigners for MsgDisableServiceBinding
func TestMsgDisableServiceBindingGetSigners(t *testing.T) {
	msg := NewMsgDisableServiceBinding(testServiceName, testProvider, testOwner)
	res := msg.GetSigners()
	require.Equal(t, "[6AC7D958DE501C93A3300A8551C85A0FADF2637B]", fmt.Sprintf("%v", res))
}

// TestMsgEnableServiceBindingRoute tests Route for MsgEnableServiceBinding
func TestMsgEnableServiceBindingRoute(t *testing.T) {
	msg := NewMsgEnableServiceBinding(testServiceName, testProvider, testAddedDeposit, testOwner)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgEnableServiceBindingType tests Type for MsgEnableServiceBinding
func TestMsgEnableServiceBindingType(t *testing.T) {
	msg := NewMsgEnableServiceBinding(testServiceName, testProvider, testAddedDeposit, testOwner)

	require.Equal(t, "enable_service_binding", msg.Type())
}

// TestMsgEnableServiceBindingValidation tests ValidateBasic for MsgEnableServiceBinding
func TestMsgEnableServiceBindingValidation(t *testing.T) {
	emptyAddedDeposit := sdk.Coins{}

	invalidLongName := strings.Repeat("s", MaxNameLength+1)

	testMsgs := []*MsgEnableServiceBinding{
		NewMsgEnableServiceBinding(testServiceName, testProvider, testAddedDeposit, testOwner),    // valid msg
		NewMsgEnableServiceBinding(testServiceName, testProvider, emptyAddedDeposit, testOwner),   // empty deposit is allowed
		NewMsgEnableServiceBinding(testServiceName, emptyAddress, testAddedDeposit, testOwner),    // missing provider address
		NewMsgEnableServiceBinding(testServiceName, testProvider, testAddedDeposit, emptyAddress), // missing owner address
		NewMsgEnableServiceBinding(invalidName, testProvider, testAddedDeposit, testOwner),        // service name contains illegal characters
		NewMsgEnableServiceBinding(invalidLongName, testProvider, testAddedDeposit, testOwner),    // too long service name
	}

	testCases := []struct {
		msg     *MsgEnableServiceBinding
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], true, ""},
		{testMsgs[2], false, "missing provider address"},
		{testMsgs[3], false, "missing owner address"},
		{testMsgs[4], false, "service name contains illegal characters"},
		{testMsgs[5], false, "too long service name"},
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

// TestMsgEnableServiceBindingGetSignBytes tests GetSignBytes for MsgEnableServiceBinding
func TestMsgEnableServiceBindingGetSignBytes(t *testing.T) {
	msg := NewMsgEnableServiceBinding(testServiceName, testProvider, testAddedDeposit, testOwner)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgEnableServiceBinding","value":{"deposit":[{"amount":"100","denom":"stake"}],"owner":"cosmos1dtrajkx72qwf8gesp2z4rjz6p7klycmmh2lnvn","provider":"cosmos10yajkmnug9d2nvccgs7ul3xkur0l9as348gmwe","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgEnableServiceBindingGetSigners tests GetSigners for MsgEnableServiceBinding
func TestMsgEnableServiceBindingGetSigners(t *testing.T) {
	msg := NewMsgEnableServiceBinding(testServiceName, testProvider, testAddedDeposit, testOwner)
	res := msg.GetSigners()
	require.Equal(t, "[6AC7D958DE501C93A3300A8551C85A0FADF2637B]", fmt.Sprintf("%v", res))
}

// TestMsgRefundServiceDepositRoute tests Route for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositRoute(t *testing.T) {
	msg := NewMsgRefundServiceDeposit(testServiceName, testProvider, testOwner)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgRefundServiceDepositType tests Type for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositType(t *testing.T) {
	msg := NewMsgRefundServiceDeposit(testServiceName, testProvider, testOwner)

	require.Equal(t, "refund_service_deposit", msg.Type())
}

// TestMsgRefundServiceDepositValidation tests ValidateBasic for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositValidation(t *testing.T) {
	invalidLongName := strings.Repeat("s", MaxNameLength+1)

	testMsgs := []*MsgRefundServiceDeposit{
		NewMsgRefundServiceDeposit(testServiceName, testProvider, testOwner),    // valid msg
		NewMsgRefundServiceDeposit(testServiceName, emptyAddress, testOwner),    // missing provider address
		NewMsgRefundServiceDeposit(testServiceName, testProvider, emptyAddress), // missing owner address
		NewMsgRefundServiceDeposit(invalidName, testProvider, testOwner),        // service name contains illegal characters
		NewMsgRefundServiceDeposit(invalidLongName, testProvider, testOwner),    // too long service name
	}

	testCases := []struct {
		msg     *MsgRefundServiceDeposit
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
		{testMsgs[2], false, "missing owner address"},
		{testMsgs[3], false, "service name contains illegal characters"},
		{testMsgs[4], false, "too long service name"},
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
	msg := NewMsgRefundServiceDeposit(testServiceName, testProvider, testOwner)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgRefundServiceDeposit","value":{"owner":"cosmos1dtrajkx72qwf8gesp2z4rjz6p7klycmmh2lnvn","provider":"cosmos10yajkmnug9d2nvccgs7ul3xkur0l9as348gmwe","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgRefundServiceDepositGetSigners tests GetSigners for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositGetSigners(t *testing.T) {
	msg := NewMsgRefundServiceDeposit(testServiceName, testProvider, testOwner)
	res := msg.GetSigners()
	require.Equal(t, "[6AC7D958DE501C93A3300A8551C85A0FADF2637B]", fmt.Sprintf("%v", res))
}

// TestMsgCallServiceRoute tests Route for MsgCallService
func TestMsgCallServiceRoute(t *testing.T) {
	msg := NewMsgCallService(
		testServiceName, testProviders, testConsumer,
		testInput, testServiceFeeCap, testTimeout,
		true, testRepeatedFreq, testRepeatedTotal,
	)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgCallServiceType tests Type for MsgCallService
func TestMsgCallServiceType(t *testing.T) {
	msg := NewMsgCallService(
		testServiceName, testProviders, testConsumer,
		testInput, testServiceFeeCap, testTimeout,
		true, testRepeatedFreq, testRepeatedTotal,
	)

	require.Equal(t, "call_service", msg.Type())
}

// TestMsgCallServiceValidation tests ValidateBasic for MsgCallService
func TestMsgCallServiceValidation(t *testing.T) {
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidDenomCoins := sdk.Coins{sdk.Coin{Denom: "eth-min", Amount: sdk.NewInt(1000)}}

	invalidDuplicateProviders := []string{testProvider, testProvider}
	invalidInput := "iris-usdt"
	invalidTimeout := int64(0)
	invalidLessRepeatedFreq := uint64(testTimeout) - 10
	invalidRepeatedTotal1 := int64(-2)
	invalidRepeatedTotal2 := int64(0)

	testMsgs := []*MsgCallService{
		NewMsgCallService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, true, testRepeatedFreq, testRepeatedTotal,
		), // valid msg
		NewMsgCallService(
			testServiceName, testProviders, emptyAddress, testInput, testServiceFeeCap,
			testTimeout, true, testRepeatedFreq, testRepeatedTotal,
		), // missing consumer address
		NewMsgCallService(
			invalidName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, true, testRepeatedFreq, testRepeatedTotal,
		), // service name contains illegal characters
		NewMsgCallService(
			invalidLongName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, true, testRepeatedFreq, testRepeatedTotal,
		), // too long service name
		NewMsgCallService(
			testServiceName, testProviders, testConsumer, testInput, invalidDenomCoins,
			testTimeout, true, testRepeatedFreq, testRepeatedTotal,
		), // invalid service fee denom
		NewMsgCallService(
			testServiceName, nil, testConsumer, testInput, testServiceFeeCap,
			testTimeout, true, testRepeatedFreq, testRepeatedTotal,
		), // missing providers
		NewMsgCallService(
			testServiceName, invalidDuplicateProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, true, testRepeatedFreq, testRepeatedTotal,
		), // duplicate providers
		NewMsgCallService(
			testServiceName, testProviders, testConsumer, "", testServiceFeeCap,
			testTimeout, true, testRepeatedFreq, testRepeatedTotal,
		), // missing input
		NewMsgCallService(
			testServiceName, testProviders, testConsumer, invalidInput, testServiceFeeCap,
			testTimeout, true, testRepeatedFreq, testRepeatedTotal,
		), // invalid input
		NewMsgCallService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			invalidTimeout, true, testRepeatedFreq, testRepeatedTotal,
		), // invalid timeout
		NewMsgCallService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, true, invalidLessRepeatedFreq, testRepeatedTotal,
		), // invalid repeated frequency
		NewMsgCallService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, true, testRepeatedFreq, invalidRepeatedTotal1,
		), // repeated total can not be less than -1
		NewMsgCallService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, true, testRepeatedFreq, invalidRepeatedTotal2,
		), // repeated total can not be zero
		NewMsgCallService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, true, uint64(0), testRepeatedTotal,
		), // frequency can be zero
		NewMsgCallService(
			testServiceName, testProviders, testConsumer, testInput, testServiceFeeCap,
			testTimeout, false, invalidLessRepeatedFreq, invalidRepeatedTotal1,
		), // do not check the repeated frequency and total when not repeated
	}

	testCases := []struct {
		msg     *MsgCallService
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

// TestMsgCallServiceGetSignBytes tests GetSignBytes for MsgCallService
func TestMsgCallServiceGetSignBytes(t *testing.T) {
	msg := NewMsgCallService(
		testServiceName, testProviders, testConsumer,
		testInput, testServiceFeeCap, testTimeout,
		true, testRepeatedFreq, testRepeatedTotal,
	)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgCallService","value":{"consumer":"cosmos1d8ydkv60gkj3sc98lvnxpddlfwwluvp6jc59a0","input":"{\"header\":{},\"body\":{\"pair\":\"iris-usdt\"}}","providers":["cosmos10yajkmnug9d2nvccgs7ul3xkur0l9as348gmwe"],"repeated":true,"repeated_frequency":"120","repeated_total":"100","service_fee_cap":[{"amount":"100","denom":"stake"}],"service_name":"test-service","timeout":"100"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgCallServiceGetSigners tests GetSigners for MsgCallService
func TestMsgCallServiceGetSigners(t *testing.T) {
	msg := NewMsgCallService(
		testServiceName, testProviders, testConsumer,
		testInput, testServiceFeeCap, testTimeout,
		true, testRepeatedFreq, testRepeatedTotal,
	)
	res := msg.GetSigners()
	require.Equal(t, "[69C8DB334F45A51860A7FB2660B5BF4B9DFE303A]", fmt.Sprintf("%v", res))
}

// TestMsgRespondServiceRoute tests Route for MsgRespondService
func TestMsgRespondServiceRoute(t *testing.T) {
	msg := NewMsgRespondService(testRequestID, testProvider, testResult, testOutput)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgRespondServiceType tests Type for MsgRespondService
func TestMsgRespondServiceType(t *testing.T) {
	msg := NewMsgRespondService(testRequestID, testProvider, testResult, testOutput)

	require.Equal(t, "respond_service", msg.Type())
}

// TestMsgRespondServiceValidation tests ValidateBasic for MsgRespondService
func TestMsgRespondServiceValidation(t *testing.T) {
	invalidRequestID := "invalidRequestID"
	invalidOutput := "invalidOutput"

	validResult400 := `{"code":400,"message":"invalid parameters"}`
	invalidResult := "invalidResult"
	invalidResultCode := `{"code":100,"message":""}`
	invalidResultNoCode := `{"message":""}`
	invalidResultNoMsg := `{"code":200}`

	testMsgs := []*MsgRespondService{
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
		msg     *MsgRespondService
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

	expected := `{"type":"irismod/service/MsgRespondService","value":{"output":"{\"header\":{},\"body\":{\"last\":\"100\"}}","provider":"cosmos10yajkmnug9d2nvccgs7ul3xkur0l9as348gmwe","request_id":"3DB0FA99DCB058BC86041BADBD614D6839F8FA20E17CF8AD3BA14C3F1BF613BD0000000000000000000000000000000100000000000000010001","result":"{\"code\":200,\"message\":\"\"}"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgRespondServiceGetSigners tests GetSigners for MsgRespondService
func TestMsgRespondServiceGetSigners(t *testing.T) {
	msg := NewMsgRespondService(testRequestID, testProvider, testResult, testOutput)
	res := msg.GetSigners()

	expected := "[793B2B6E7C415AA9B318443DCFC4D6E0DFF2F611]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgPauseRequestContextRoute tests Route for MsgPauseRequestContext
func TestMsgPauseRequestContextRoute(t *testing.T) {
	msg := NewMsgPauseRequestContext(testRequestContextID, testConsumer)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgPauseRequestContextType tests Type for MsgPauseRequestContext
func TestMsgPauseRequestContextType(t *testing.T) {
	msg := NewMsgPauseRequestContext(testRequestContextID, testConsumer)

	require.Equal(t, "pause_request_context", msg.Type())
}

// TestMsgPauseRequestContextValidation tests ValidateBasic for MsgPauseRequestContext
func TestMsgPauseRequestContextValidation(t *testing.T) {
	invalidRequestContextID := "invalid-request-context-id"

	testMsgs := []*MsgPauseRequestContext{
		NewMsgPauseRequestContext(testRequestContextID, testConsumer),    // valid msg
		NewMsgPauseRequestContext(testRequestContextID, emptyAddress),    // missing consumer address
		NewMsgPauseRequestContext(invalidRequestContextID, testConsumer), // invalid request context ID
	}

	testCases := []struct {
		msg     *MsgPauseRequestContext
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

	expected := `{"type":"irismod/service/MsgPauseRequestContext","value":{"consumer":"cosmos1d8ydkv60gkj3sc98lvnxpddlfwwluvp6jc59a0","request_context_id":"3DB0FA99DCB058BC86041BADBD614D6839F8FA20E17CF8AD3BA14C3F1BF613BD0000000000000000"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgPauseRequestContextGetSigners tests GetSigners for MsgPauseRequestContext
func TestMsgPauseRequestContextGetSigners(t *testing.T) {
	msg := NewMsgPauseRequestContext(testRequestContextID, testConsumer)
	res := msg.GetSigners()
	require.Equal(t, "[69C8DB334F45A51860A7FB2660B5BF4B9DFE303A]", fmt.Sprintf("%v", res))
}

// TestMsgStartRequestContextRoute tests Route for MsgStartRequestContext
func TestMsgStartRequestContextRoute(t *testing.T) {
	msg := NewMsgStartRequestContext(testRequestContextID, testConsumer)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgStartRequestContextType tests Type for MsgStartRequestContext
func TestMsgStartRequestContextType(t *testing.T) {
	msg := NewMsgStartRequestContext(testRequestContextID, testConsumer)

	require.Equal(t, "start_request_context", msg.Type())
}

// TestMsgStartRequestContextValidation tests ValidateBasic for MsgStartRequestContext
func TestMsgStartRequestContextValidation(t *testing.T) {
	invalidRequestContextID := "invalid-request-context-id"

	testMsgs := []*MsgStartRequestContext{
		NewMsgStartRequestContext(testRequestContextID, testConsumer),    // valid msg
		NewMsgStartRequestContext(testRequestContextID, emptyAddress),    // missing consumer address
		NewMsgStartRequestContext(invalidRequestContextID, testConsumer), // invalid request context ID
	}

	testCases := []struct {
		msg     *MsgStartRequestContext
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

	expected := `{"type":"irismod/service/MsgStartRequestContext","value":{"consumer":"cosmos1d8ydkv60gkj3sc98lvnxpddlfwwluvp6jc59a0","request_context_id":"3DB0FA99DCB058BC86041BADBD614D6839F8FA20E17CF8AD3BA14C3F1BF613BD0000000000000000"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgStartRequestContextGetSigners tests GetSigners for MsgStartRequestContext
func TestMsgStartRequestContextGetSigners(t *testing.T) {
	msg := NewMsgStartRequestContext(testRequestContextID, testConsumer)
	res := msg.GetSigners()
	require.Equal(t, "[69C8DB334F45A51860A7FB2660B5BF4B9DFE303A]", fmt.Sprintf("%v", res))
}

// TestMsgKillRequestContextRoute tests Route for MsgKillRequestContext
func TestMsgKillRequestContextRoute(t *testing.T) {
	msg := NewMsgKillRequestContext(testRequestContextID, testConsumer)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgKillRequestContextType tests Type for MsgKillRequestContext
func TestMsgKillRequestContextType(t *testing.T) {
	msg := NewMsgKillRequestContext(testRequestContextID, testConsumer)

	require.Equal(t, "kill_request_context", msg.Type())
}

// TestMsgKillRequestContextValidation tests ValidateBasic for MsgKillRequestContext
func TestMsgKillRequestContextValidation(t *testing.T) {
	invalidRequestContextID := "invalid-request-context-id"

	testMsgs := []*MsgKillRequestContext{
		NewMsgKillRequestContext(testRequestContextID, testConsumer),    // valid msg
		NewMsgKillRequestContext(testRequestContextID, emptyAddress),    // missing consumer address
		NewMsgKillRequestContext(invalidRequestContextID, testConsumer), // invalid request context ID
	}

	testCases := []struct {
		msg     *MsgKillRequestContext
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

	expected := `{"type":"irismod/service/MsgKillRequestContext","value":{"consumer":"cosmos1d8ydkv60gkj3sc98lvnxpddlfwwluvp6jc59a0","request_context_id":"3DB0FA99DCB058BC86041BADBD614D6839F8FA20E17CF8AD3BA14C3F1BF613BD0000000000000000"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgKillRequestContextGetSigners tests GetSigners for MsgKillRequestContext
func TestMsgKillRequestContextGetSigners(t *testing.T) {
	msg := NewMsgKillRequestContext(testRequestContextID, testConsumer)
	res := msg.GetSigners()
	require.Equal(t, "[69C8DB334F45A51860A7FB2660B5BF4B9DFE303A]", fmt.Sprintf("%v", res))
}

// TestMsgUpdateRequestContextRoute tests Route for MsgUpdateRequestContext
func TestMsgUpdateRequestContextRoute(t *testing.T) {
	msg := NewMsgUpdateRequestContext(testRequestContextID, nil, nil, 0, 0, 0, testConsumer)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgUpdateRequestContextType tests Type for MsgUpdateRequestContext
func TestMsgUpdateRequestContextType(t *testing.T) {
	msg := NewMsgUpdateRequestContext(testRequestContextID, nil, nil, 0, 0, 0, testConsumer)

	require.Equal(t, "update_request_context", msg.Type())
}

// TestMsgUpdateRequestContextValidation tests ValidateBasic for MsgUpdateRequestContext
func TestMsgUpdateRequestContextValidation(t *testing.T) {
	invalidRequestContextID := "invalid-request-context-id"
	invalidDuplicateProviders := []string{testProvider, testProvider}
	invalidTimeout := int64(-1)
	invalidLessRepeatedFreq := uint64(testTimeout) - 10
	invalidRepeatedTotal := int64(-2)
	invalidDenomCoins := sdk.Coins{sdk.Coin{Denom: "eth-min", Amount: sdk.NewInt(1000)}}

	testMsgs := []*MsgUpdateRequestContext{
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
		msg     *MsgUpdateRequestContext
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

	expected := `{"type":"irismod/service/MsgUpdateRequestContext","value":{"consumer":"cosmos1d8ydkv60gkj3sc98lvnxpddlfwwluvp6jc59a0","providers":["cosmos10yajkmnug9d2nvccgs7ul3xkur0l9as348gmwe"],"repeated_frequency":"120","repeated_total":"100","request_context_id":"3DB0FA99DCB058BC86041BADBD614D6839F8FA20E17CF8AD3BA14C3F1BF613BD0000000000000000","service_fee_cap":[{"amount":"100","denom":"stake"}],"timeout":"100"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgUpdateRequestContextGetSigners tests GetSigners for MsgUpdateRequestContext
func TestMsgUpdateRequestContextGetSigners(t *testing.T) {
	msg := NewMsgUpdateRequestContext(testRequestContextID, testProviders, testServiceFeeCap, testTimeout, testRepeatedFreq, testRepeatedTotal, testConsumer)
	res := msg.GetSigners()
	require.Equal(t, "[69C8DB334F45A51860A7FB2660B5BF4B9DFE303A]", fmt.Sprintf("%v", res))
}

// TestMsgWithdrawEarnedFeesRoute tests Route for MsgWithdrawEarnedFees
func TestMsgWithdrawEarnedFeesRoute(t *testing.T) {
	msg := NewMsgWithdrawEarnedFees(testOwner, testProvider)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgWithdrawEarnedFeesType tests Type for MsgWithdrawEarnedFees
func TestMsgWithdrawEarnedFeesType(t *testing.T) {
	msg := NewMsgWithdrawEarnedFees(testOwner, testProvider)

	require.Equal(t, "withdraw_earned_fees", msg.Type())
}

// TestMsgWithdrawEarnedFeesValidation tests ValidateBasic for MsgWithdrawEarnedFees
func TestMsgWithdrawEarnedFeesValidation(t *testing.T) {
	testMsgs := []*MsgWithdrawEarnedFees{
		NewMsgWithdrawEarnedFees(testOwner, testProvider),    // valid msg
		NewMsgWithdrawEarnedFees(testOwner, emptyAddress),    // empty provider is allowed
		NewMsgWithdrawEarnedFees(emptyAddress, testProvider), // missing owner address
	}

	testCases := []struct {
		msg     *MsgWithdrawEarnedFees
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], true, ""},
		{testMsgs[2], false, "missing owner address"},
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
	msg := NewMsgWithdrawEarnedFees(testOwner, testProvider)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgWithdrawEarnedFees","value":{"owner":"cosmos1dtrajkx72qwf8gesp2z4rjz6p7klycmmh2lnvn","provider":"cosmos10yajkmnug9d2nvccgs7ul3xkur0l9as348gmwe"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgWithdrawEarnedFeesGetSigners tests GetSigners for MsgWithdrawEarnedFees
func TestMsgWithdrawEarnedFeesGetSigners(t *testing.T) {
	msg := NewMsgWithdrawEarnedFees(testOwner, testProvider)
	res := msg.GetSigners()
	require.Equal(t, "[6AC7D958DE501C93A3300A8551C85A0FADF2637B]", fmt.Sprintf("%v", res))
}
