package asset

import (
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

// test ValidateBasic for MsgCreateValidator
func TestMsgIssueAsset(t *testing.T) {
	addr := sdk.AccAddress("test")
	tests := []struct {
		testCase string
		MsgIssueAsset
		expectPass bool
	}{
		{"basic good", MsgIssueAsset{"00", "b", "c", "d", 1, 1, 1, true, addr, []sdk.AccAddress{addr}}, true},
		//{"empty family", MsgIssueAsset{"", "b", "c", "d", 1, 1, 1, true, addr, []sdk.AccAddress{addr}}, false},
		{"error family", MsgIssueAsset{"02", "b", "c", "d", 1, 1, 1, true, addr, []sdk.AccAddress{addr}}, false},
		{"empty name", MsgIssueAsset{"00", "", "c", "d", 1, 1, 1, true, addr, []sdk.AccAddress{addr}}, false},
		{"error name", MsgIssueAsset{"00", ".,re323", "c", "d", 1, 1, 1, true, addr, []sdk.AccAddress{addr}}, false},
		{"empty symbol", MsgIssueAsset{"00", "b", "", "d", 1, 1, 1, true, addr, []sdk.AccAddress{addr}}, false},
		{"error symbol", MsgIssueAsset{"00", "b", "232,e4", "d", 1, 1, 1, true, addr, []sdk.AccAddress{addr}}, false},
		{"zero supply bigger", MsgIssueAsset{"00", "b", "c", "d", 0, 1, 1, true, addr, []sdk.AccAddress{addr}}, false},
		{"zero max supply", MsgIssueAsset{"00", "b", "c", "d", 1, 0, 1, true, addr, []sdk.AccAddress{addr}}, false},
		{"too much supply", MsgIssueAsset{"00", "b", "c", "d", 1e+13, 1e+13, 1, true, addr, []sdk.AccAddress{addr}}, false},
		{"too much max supply", MsgIssueAsset{"00", "b", "c", "d", 1, 1e+13, 1, true, addr, []sdk.AccAddress{addr}}, false},
		{"init supply bigger than max supply", MsgIssueAsset{"00", "b", "c", "d", 10, 9, 1, true, addr, []sdk.AccAddress{addr}}, false},
		{"error decimal", MsgIssueAsset{"00", "b", "c", "d", 1, 1, 19, true, addr, []sdk.AccAddress{addr}}, false},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgIssueAsset.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgIssueAsset.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}
