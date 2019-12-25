package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/config"
)

// nolint: deadcode unused
var (
	sender   sdk.AccAddress
	testAddr sdk.AccAddress

	nilAddr        = sdk.AccAddress{}
	description    = "description"
	nilDescription = ""
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount(config.GetConfig().GetBech32AccountAddrPrefix(), config.GetConfig().GetBech32AccountPubPrefix())
	sdk.GetConfig().SetBech32PrefixForValidator(config.GetConfig().GetBech32ValidatorAddrPrefix(), config.GetConfig().GetBech32ValidatorPubPrefix())
	sdk.GetConfig().SetBech32PrefixForConsensusNode(config.GetConfig().GetBech32ConsensusAddrPrefix(), config.GetConfig().GetBech32ConsensusPubPrefix())

	sender, _ = sdk.AccAddressFromBech32("faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj")
	testAddr, _ = sdk.AccAddressFromBech32("faa1mrehjkgeg75nz2gk7lr7dnxvvtg4497jxss8hq")
}

// ----------------------------------------------
// test MsgAddProfiler
// ----------------------------------------------

func TestNewMsgAddProfiler(t *testing.T) {
	addGuardian := AddGuardian{description, testAddr, sender}
	msg := NewMsgAddProfiler(description, testAddr, sender)
	require.Equal(t, addGuardian, msg.AddGuardian)
}

func TestMsgAddProfilerRoute(t *testing.T) {
	msg := NewMsgAddProfiler(description, testAddr, sender)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgAddProfilerType(t *testing.T) {
	msg := NewMsgAddProfiler(description, testAddr, sender)
	require.Equal(t, TypeMsgAddProfiler, msg.Type())
}

func TestMsgAddProfilerGetSignBytes(t *testing.T) {
	msg := NewMsgAddProfiler(description, testAddr, sender)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/guardian/MsgAddProfiler","value":{"AddGuardian":{"added_by":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj","address":"faa1mrehjkgeg75nz2gk7lr7dnxvvtg4497jxss8hq","description":"description"}}}`
	require.Equal(t, expected, string(res))
}

func TestMsgAddProfilerGetSigners(t *testing.T) {
	msg := NewMsgAddProfiler(description, testAddr, sender)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// test ValidateBasic for MsgAddProfiler
func TestMsgAddProfilerValidation(t *testing.T) {
	tests := []struct {
		name       string
		expectPass bool
		msg        MsgAddProfiler
	}{
		{"pass", true, NewMsgAddProfiler(description, testAddr, sender)},
		{"invalid Description", false, NewMsgAddProfiler(nilDescription, testAddr, sender)},
		{"invalid Address", false, NewMsgAddProfiler(description, nilAddr, sender)},
		{"invalid AddedBy", false, NewMsgAddProfiler(description, testAddr, nilAddr)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectPass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

// ----------------------------------------------
// test MsgDeleteProfiler
// ----------------------------------------------

func TestNewMsgDeleteProfiler(t *testing.T) {
	deleteGuardian := DeleteGuardian{testAddr, sender}
	msg := NewMsgDeleteProfiler(testAddr, sender)
	require.Equal(t, deleteGuardian, msg.DeleteGuardian)
}

func TestMsgDeleteProfilerRoute(t *testing.T) {
	msg := NewMsgDeleteProfiler(testAddr, sender)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgDeleteProfilerType(t *testing.T) {
	msg := NewMsgDeleteProfiler(testAddr, sender)
	require.Equal(t, TypeMsgDeleteProfiler, msg.Type())
}

func TestMsgDeleteProfilerGetSignBytes(t *testing.T) {
	msg := NewMsgDeleteProfiler(testAddr, sender)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/guardian/MsgDeleteProfiler","value":{"DeleteGuardian":{"address":"faa1mrehjkgeg75nz2gk7lr7dnxvvtg4497jxss8hq","deleted_by":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj"}}}`
	require.Equal(t, expected, string(res))
}

func TestMsgDeleteProfilerGetSigners(t *testing.T) {
	msg := NewMsgDeleteProfiler(testAddr, sender)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// test ValidateBasic for MsgDeleteProfiler
func TestMsgDeleteProfilerValidation(t *testing.T) {
	tests := []struct {
		name       string
		expectPass bool
		msg        MsgDeleteProfiler
	}{
		{"pass", true, NewMsgDeleteProfiler(testAddr, sender)},
		{"invalid Address", false, NewMsgDeleteProfiler(nilAddr, sender)},
		{"invalid DeletedBy", false, NewMsgDeleteProfiler(testAddr, nilAddr)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectPass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

// ----------------------------------------------
// test MsgAddTrustee
// ----------------------------------------------

func TestNewMsgAddTrustee(t *testing.T) {
	addGuardian := AddGuardian{description, testAddr, sender}
	msg := NewMsgAddTrustee(description, testAddr, sender)
	require.Equal(t, addGuardian, msg.AddGuardian)
}

func TestMsgAddTrusteeRoute(t *testing.T) {
	msg := NewMsgAddTrustee(description, testAddr, sender)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgAddTrusteeType(t *testing.T) {
	msg := NewMsgAddTrustee(description, testAddr, sender)
	require.Equal(t, TypeMsgAddTrustee, msg.Type())
}

func TestMsgAddTrusteeGetSignBytes(t *testing.T) {
	msg := NewMsgAddTrustee(description, testAddr, sender)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/guardian/MsgAddTrustee","value":{"AddGuardian":{"added_by":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj","address":"faa1mrehjkgeg75nz2gk7lr7dnxvvtg4497jxss8hq","description":"description"}}}`
	require.Equal(t, expected, string(res))
}

func TestMsgAddTrusteeGetSigners(t *testing.T) {
	msg := NewMsgAddTrustee(description, testAddr, sender)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// test ValidateBasic for MsgAddTrustee
func TestMsgAddTrusteeValidation(t *testing.T) {
	tests := []struct {
		name       string
		expectPass bool
		msg        MsgAddTrustee
	}{
		{"pass", true, NewMsgAddTrustee(description, testAddr, sender)},
		{"invalid Description", false, NewMsgAddTrustee(nilDescription, testAddr, sender)},
		{"invalid Address", false, NewMsgAddTrustee(description, nilAddr, sender)},
		{"invalid AddedBy", false, NewMsgAddTrustee(description, testAddr, nilAddr)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectPass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

// ----------------------------------------------
// test MsgDeleteTrustee
// ----------------------------------------------

func TestNewMsgDeleteTrustee(t *testing.T) {
	deleteGuardian := DeleteGuardian{testAddr, sender}
	msg := NewMsgDeleteTrustee(testAddr, sender)
	require.Equal(t, deleteGuardian, msg.DeleteGuardian)
}

func TestMsgDeleteTrusteeRoute(t *testing.T) {
	msg := NewMsgDeleteTrustee(testAddr, sender)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgDeleteTrusteeType(t *testing.T) {
	msg := NewMsgDeleteTrustee(testAddr, sender)
	require.Equal(t, TypeMsgDeleteTrustee, msg.Type())
}

func TestMsgDeleteTrusteeGetSignBytes(t *testing.T) {
	msg := NewMsgDeleteTrustee(testAddr, sender)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/guardian/MsgDeleteTrustee","value":{"DeleteGuardian":{"address":"faa1mrehjkgeg75nz2gk7lr7dnxvvtg4497jxss8hq","deleted_by":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj"}}}`
	require.Equal(t, expected, string(res))
}

func TestMsgDeleteTrusteeGetSigners(t *testing.T) {
	msg := NewMsgDeleteTrustee(testAddr, sender)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// test ValidateBasic for MsgDeleteTrustee
func TestMsgDeleteTrusteeValidation(t *testing.T) {
	tests := []struct {
		name       string
		expectPass bool
		msg        MsgDeleteTrustee
	}{
		{"pass", true, NewMsgDeleteTrustee(testAddr, sender)},
		{"invalid Address", false, NewMsgDeleteTrustee(nilAddr, sender)},
		{"invalid DeletedBy", false, NewMsgDeleteTrustee(testAddr, nilAddr)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectPass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
