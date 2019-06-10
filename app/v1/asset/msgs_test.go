package asset

import (
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

// test ValidateBasic for MsgIssueAsset
func TestMsgIssueAsset(t *testing.T) {
	addr := sdk.AccAddress("test")
	tests := []struct {
		testCase string
		MsgIssueAsset
		expectPass bool
	}{
		{"basic good", MsgIssueAsset{BaseAsset{0x00, 0x00, "c", "d", "e", 1, "f", 1, 1, true, addr}}, true},
		{"error family", MsgIssueAsset{BaseAsset{0x02, 0x00, "c", "d", "e", 1, "f", 1, 1, true, addr}}, false},
		{"error source", MsgIssueAsset{BaseAsset{0x00, 0x03, "c", "d", "e", 1, "f", 1, 1, true, addr}}, false},
		{"empty symbol", MsgIssueAsset{BaseAsset{0x00, 0x00, "c", "", "e", 1, "f", 1, 1, true, addr}}, false},
		{"error symbol", MsgIssueAsset{BaseAsset{0x00, 0x00, "c", "434,23d", "e", 1, "f", 1, 1, true, addr}}, false},
		{"empty name", MsgIssueAsset{BaseAsset{0x00, 0x00, "c", "d", "", 1, "f", 1, 1, true, addr}}, false},
		{"error name", MsgIssueAsset{BaseAsset{0x00, 0x00, "c", "d", "2123<s", 1, "f", 1, 1, true, addr}}, false},
		{"zero supply bigger", MsgIssueAsset{BaseAsset{0x00, 0x00, "c", "d", "e", 1, "f", 0, 1, true, addr}}, false},
		{"zero max supply", MsgIssueAsset{BaseAsset{0x00, 0x00, "c", "d", "e", 1, "f", 1, 0, true, addr}}, false},
		{"init supply bigger than max supply", MsgIssueAsset{BaseAsset{0x00, 0x00, "c", "d", "e", 1, "f", 10, 9, true, addr}}, false},
		{"error decimal", MsgIssueAsset{BaseAsset{0x00, 0x00, "c", "d", "e", 19, "f", 1, 1, true, addr}}, false},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgIssueAsset.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgIssueAsset.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}
