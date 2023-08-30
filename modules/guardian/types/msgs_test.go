package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/v2/address"
)

// nolint: deadcode unused
var (
	sender, _   = sdk.AccAddressFromHexUnsafe(crypto.AddressHash([]byte("sender")).String())
	testAddr, _ = sdk.AccAddressFromHexUnsafe(crypto.AddressHash([]byte("test")).String())

	nilAddr        = sdk.AccAddress{}
	description    = "description"
	nilDescription = ""
)

func init() {
	address.ConfigureBech32Prefix()
}

// ----------------------------------------------
// test MsgAddSuper
// ----------------------------------------------

func TestNewMsgAddSuper(t *testing.T) {
	msg := NewMsgAddSuper(description, testAddr, sender)
	require.Equal(t, description, msg.Description)
	require.Equal(t, testAddr.String(), msg.Address)
	require.Equal(t, sender.String(), msg.AddedBy)
}

func TestMsgAddSuperRoute(t *testing.T) {
	msg := NewMsgAddSuper(description, testAddr, sender)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgAddSuperType(t *testing.T) {
	msg := NewMsgAddSuper(description, testAddr, sender)
	require.Equal(t, TypeMsgAddSuper, msg.Type())
}

func TestMsgAddSuperGetSignBytes(t *testing.T) {
	msg := NewMsgAddSuper(description, testAddr, sender)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/guardian/MsgAddSuper","value":{"added_by":"iaa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgwp4wlf","address":"iaa1n7rdpqvgf37ktx30a2sv2kkszk3m7ncmakdj4g","description":"description"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgAddSuperGetSigners(t *testing.T) {
	msg := NewMsgAddSuper(description, testAddr, sender)
	res := msg.GetSigners()
	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// test ValidateBasic for MsgAddSuper
func TestMsgAddSuperValidation(t *testing.T) {
	tests := []struct {
		name       string
		expectPass bool
		msg        *MsgAddSuper
	}{
		{"pass", true, NewMsgAddSuper(description, testAddr, sender)},
		{"invalid Description", false, NewMsgAddSuper(nilDescription, testAddr, sender)},
		{"invalid Address", false, NewMsgAddSuper(description, nilAddr, sender)},
		{"invalid AddedBy", false, NewMsgAddSuper(description, testAddr, nilAddr)},
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
// test MsgDeleteSuper
// ----------------------------------------------

func TestNewMsgDeleteSuper(t *testing.T) {
	msg := NewMsgDeleteSuper(testAddr, sender)
	require.Equal(t, testAddr.String(), msg.Address)
	require.Equal(t, sender.String(), msg.DeletedBy)
}

func TestMsgDeleteSuperRoute(t *testing.T) {
	msg := NewMsgDeleteSuper(testAddr, sender)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgDeleteSuperType(t *testing.T) {
	msg := NewMsgDeleteSuper(testAddr, sender)
	require.Equal(t, TypeMsgDeleteSuper, msg.Type())
}

func TestMsgDeleteSuperGetSignBytes(t *testing.T) {
	msg := NewMsgDeleteSuper(testAddr, sender)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/guardian/MsgDeleteSuper","value":{"address":"iaa1n7rdpqvgf37ktx30a2sv2kkszk3m7ncmakdj4g","deleted_by":"iaa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgwp4wlf"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgDeleteSuperGetSigners(t *testing.T) {
	msg := NewMsgDeleteSuper(testAddr, sender)
	res := msg.GetSigners()
	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// test ValidateBasic for MsgDeleteSuper
func TestMsgDeleteSuperValidation(t *testing.T) {
	tests := []struct {
		name       string
		expectPass bool
		msg        *MsgDeleteSuper
	}{
		{"pass", true, NewMsgDeleteSuper(testAddr, sender)},
		{"invalid Address", false, NewMsgDeleteSuper(nilAddr, sender)},
		{"invalid DeletedBy", false, NewMsgDeleteSuper(testAddr, nilAddr)},
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
