package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/config"
)

// nolint: deadcode unused
var (
	sender, _   = sdk.AccAddressFromHex(crypto.AddressHash([]byte("sender")).String())
	testAddr, _ = sdk.AccAddressFromHex(crypto.AddressHash([]byte("test")).String())

	nilAddr        = sdk.AccAddress{}
	description    = "description"
	nilDescription = ""
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount(config.GetConfig().GetBech32AccountAddrPrefix(), config.GetConfig().GetBech32AccountPubPrefix())
	sdk.GetConfig().SetBech32PrefixForValidator(config.GetConfig().GetBech32ValidatorAddrPrefix(), config.GetConfig().GetBech32ValidatorPubPrefix())
	sdk.GetConfig().SetBech32PrefixForConsensusNode(config.GetConfig().GetBech32ConsensusAddrPrefix(), config.GetConfig().GetBech32ConsensusPubPrefix())
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
	expected := `{"type":"irishub/guardian/MsgAddProfiler","value":{"add_guardian":{"added_by":"faa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgkwnkl5","address":"faa1n7rdpqvgf37ktx30a2sv2kkszk3m7ncm9et244","description":"description"}}}`
	require.Equal(t, expected, string(res))
}

func TestMsgAddProfilerGetSigners(t *testing.T) {
	msg := NewMsgAddProfiler(description, testAddr, sender)
	res := msg.GetSigners()
	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
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
				require.NoError(t, err)
			} else {
				require.Error(t, err)
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
	expected := `{"type":"irishub/guardian/MsgDeleteProfiler","value":{"delete_guardian":{"address":"faa1n7rdpqvgf37ktx30a2sv2kkszk3m7ncm9et244","deleted_by":"faa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgkwnkl5"}}}`
	require.Equal(t, expected, string(res))
}

func TestMsgDeleteProfilerGetSigners(t *testing.T) {
	msg := NewMsgDeleteProfiler(testAddr, sender)
	res := msg.GetSigners()
	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
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
				require.NoError(t, err)
			} else {
				require.Error(t, err)
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
	expected := `{"type":"irishub/guardian/MsgAddTrustee","value":{"add_guardian":{"added_by":"faa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgkwnkl5","address":"faa1n7rdpqvgf37ktx30a2sv2kkszk3m7ncm9et244","description":"description"}}}`
	require.Equal(t, expected, string(res))
}

func TestMsgAddTrusteeGetSigners(t *testing.T) {
	msg := NewMsgAddTrustee(description, testAddr, sender)
	res := msg.GetSigners()
	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
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
				require.NoError(t, err)
			} else {
				require.Error(t, err)
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
	expected := `{"type":"irishub/guardian/MsgDeleteTrustee","value":{"delete_guardian":{"address":"faa1n7rdpqvgf37ktx30a2sv2kkszk3m7ncm9et244","deleted_by":"faa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgkwnkl5"}}}`
	require.Equal(t, expected, string(res))
}

func TestMsgDeleteTrusteeGetSigners(t *testing.T) {
	msg := NewMsgDeleteTrustee(testAddr, sender)
	res := msg.GetSigners()
	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
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
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
