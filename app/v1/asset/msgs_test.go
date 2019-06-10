package asset

import (
	"fmt"
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

var (
	emptyAddr sdk.AccAddress

	addr1 = sdk.AccAddress([]byte("addr1"))
	addr2 = sdk.AccAddress([]byte("addr2"))
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

func TestNewMsgCreateGateway(t *testing.T) {}

func TestMsgCreateGatewayRoute(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	moniker := "moniker"
	identity := "identity"
	details := "details"
	website := "website"

	// build a MsgCreateGateway
	msg := MsgCreateGateway{
		Owner:    owner,
		Moniker:  moniker,
		Identity: identity,
		Details:  details,
		Website:  website,
	}

	require.Equal(t, "asset", msg.Route())
}

func TestMsgCreateGatewayValidation(t *testing.T) {
	testData := []struct {
		name                                string
		owner                               sdk.AccAddress
		moniker, identity, details, website string
		expectPass                          bool
	}{
		{"empty owner", emptyAddr, "mon", "i", "d", "w", false},
		{"empty moniker", addr1, "", "i", "d", "w", false},
		{"too short moniker", addr1, "mo", "i", "d", "w", false},
		{"too long moniker", addr1, "monikermo", "i", "d", "w", false},
		{"moniker contains illegal characters", addr2, "moni2", "i", "d", "w", false},
		{"valid msg", addr2, "moniker", "i", "d", "w", true},
	}

	for _, td := range testData {
		msg := NewMsgCreateGateway(td.owner, td.moniker, td.identity, td.details, td.website)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}

func TestMsgCreateGatewayGetSignBytes(t *testing.T) {
	var msg = MsgCreateGateway{
		Owner:    sdk.AccAddress([]byte("owner")),
		Moniker:  "moniker",
		Identity: "identity",
		Details:  "details",
		Website:  "website",
	}

	res := msg.GetSignBytes()

	expected := `{"type":"irishub/asset/MsgCreateGateway","value":{"details":"details","identity":"identity","moniker":"moniker","owner":"faa1damkuetjqqah8w","website":"website"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgCreateGatewayGetSigners(t *testing.T) {
	var msg = MsgCreateGateway{
		Owner:    sdk.AccAddress([]byte("owner")),
		Moniker:  "moniker",
		Identity: "identity",
		Details:  "details",
		Website:  "website",
	}

	res := msg.GetSigners()

	expected := "[6F776E6572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

func TestNewMsgEditGateway(t *testing.T) {}

func TestMsgEditGatewayRoute(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	moniker := "mon"
	identity := "i"
	details := "d"
	website := "w"

	// build a MsgEditGateway
	msg := MsgEditGateway{
		Owner:    owner,
		Moniker:  moniker,
		Identity: &identity,
		Details:  &details,
		Website:  &website,
	}

	require.Equal(t, "asset", msg.Route())
}

func TestMsgEditGatewayValidation(t *testing.T) {
	identity := "i"
	details := "d"
	website := "w"

	testData := []struct {
		name                       string
		owner                      sdk.AccAddress
		moniker                    string
		identity, details, website *string
		expectPass                 bool
	}{
		{"empty owner", emptyAddr, "mon", &identity, &details, &website, false},
		{"empty moniker", addr1, "", &identity, &details, &website, false},
		{"too short moniker", addr1, "mo", &identity, &details, &website, false},
		{"too long moniker", addr1, "monikermo", &identity, &details, &website, false},
		{"moniker contains illegal characters", addr2, "moni2", &identity, &details, &website, false},
		{"empty identity allowed", addr2, "mon", nil, &details, &website, true},
		{"empty details allowed", addr2, "mon", &identity, nil, &website, true},
		{"empty website allowed", addr2, "mon", &identity, &details, nil, true},
		{"no updated fields", addr2, "mon", nil, nil, nil, false},
	}

	for _, td := range testData {
		msg := NewMsgEditGateway(td.owner, td.moniker, td.identity, td.details, td.website)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}

func TestMsgEditGatewayGetSignBytes(t *testing.T) {
	identity := "i"
	details := "d"
	website := "w"

	var msg = MsgEditGateway{
		Owner:    sdk.AccAddress([]byte("owner")),
		Moniker:  "mon",
		Identity: &identity,
		Details:  &details,
		Website:  &website,
	}

	res := msg.GetSignBytes()

	expected := `{"type":"irishub/asset/MsgEditGateway","value":{"details":"d","identity":"i","moniker":"mon","owner":"faa1damkuetjqqah8w","website":"w"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgEditGatewayGetSigners(t *testing.T) {
	identity := "i"
	details := "d"
	website := "w"

	var msg = MsgEditGateway{
		Owner:    sdk.AccAddress([]byte("owner")),
		Moniker:  "mon",
		Identity: &identity,
		Details:  &details,
		Website:  &website,
	}

	res := msg.GetSigners()

	expected := "[6F776E6572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
