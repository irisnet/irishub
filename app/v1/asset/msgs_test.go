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

// test ValidateBasic for MsgIssueToken
func TestMsgIssueAsset(t *testing.T) {
	addr := sdk.AccAddress("test")
	tests := []struct {
		testCase string
		MsgIssueToken
		expectPass bool
	}{
		{"native basic good", NewMsgIssueToken(FUNGIBLE, NATIVE, "a", "btc", "btc", "btc", 18, "satoshi", 1, 1, true, addr, sdk.Coins{}), true},
		{"native family error", NewMsgIssueToken(0x02, NATIVE, "b", "btc", "btc", "btc", 1, "satoshi", 1, 1, true, addr, sdk.Coins{}), false},
		{"native source error", NewMsgIssueToken(FUNGIBLE, 0x03, "c", "btc", "btc", "btc", 1, "satoshi", 1, 1, true, addr, sdk.Coins{}), false},
		{"native symbol empty", NewMsgIssueToken(FUNGIBLE, NATIVE, "d", "", "btc", "btc", 1, "g", 1, 1, true, addr, sdk.Coins{}), false},
		{"native symbol error", NewMsgIssueToken(FUNGIBLE, NATIVE, "e", "ab,c", "btc", "btc", 1, "satoshi", 1, 1, true, addr, sdk.Coins{}), false},
		{"native symbol first letter is num", NewMsgIssueToken(FUNGIBLE, NATIVE, "e", "4iris", "btc", "btc", 1, "satoshi", 1, 1, true, addr, sdk.Coins{}), false},
		{"native symbol too long", NewMsgIssueToken(FUNGIBLE, NATIVE, "e", "aaaaaaaaa", "btc", "btc", 1, "satoshi", 1, 1, true, addr, sdk.Coins{}), false},
		{"native symbol too short", NewMsgIssueToken(FUNGIBLE, NATIVE, "e", "a", "btc", "btc", 1, "satoshi", 1, 1, true, addr, sdk.Coins{}), false},
		{"native symbol_at_source ignored", NewMsgIssueToken(FUNGIBLE, NATIVE, "a", "btc", "c", "btc", 18, "satoshi", 1, 1, true, addr, sdk.Coins{}), true},
		{"native symbol_min_alias error", NewMsgIssueToken(FUNGIBLE, NATIVE, "g", "btc", "btc", "btc", 1, "a1,3d", 1, 1, true, addr, sdk.Coins{}), false},
		{"native symbol_min_alias too long", NewMsgIssueToken(FUNGIBLE, NATIVE, "g", "btc", "btc", "btc", 1, "aaaaaaaaaaaaa", 1, 1, true, addr, sdk.Coins{}), false},
		{"native symbol_min_alias too short", NewMsgIssueToken(FUNGIBLE, NATIVE, "g", "btc", "btc", "btc", 1, "a", 1, 1, true, addr, sdk.Coins{}), false},
		{"native symbol_min_alias  first letter is num", NewMsgIssueToken(FUNGIBLE, NATIVE, "g", "btc", "btc", "btc", 1, "1a", 1, 1, true, addr, sdk.Coins{}), false},
		{"native name empty", NewMsgIssueToken(FUNGIBLE, NATIVE, "h", "btc", "btc", "", 1, "btc", 1, 1, true, addr, sdk.Coins{}), false},
		{"native name blank", NewMsgIssueToken(FUNGIBLE, NATIVE, "h", "btc", "btc", "  ", 1, "btc", 1, 1, true, addr, sdk.Coins{}), false},
		{"native name too long", NewMsgIssueToken(FUNGIBLE, NATIVE, "i", "btc", "btc", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 1, "satoshi", 1, 1, true, addr, sdk.Coins{}), false},
		{"native initial supply is zero", NewMsgIssueToken(FUNGIBLE, NATIVE, "j", "btc", "btc", "btc", 1, "satoshi", 0, 1, true, addr, sdk.Coins{}), true},
		{"native max supply is zero", NewMsgIssueToken(FUNGIBLE, NATIVE, "k", "btc", "btc", "btc", 1, "satoshi", 1, 0, true, addr, sdk.Coins{}), true},
		{"native init supply bigger than max supply", NewMsgIssueToken(FUNGIBLE, NATIVE, "l", "btc", "btc", "btc", 1, "satoshi", 2, 1, true, addr, sdk.Coins{}), false},
		{"native decimal error", NewMsgIssueToken(FUNGIBLE, NATIVE, "m", "btc", "btc", "btc", 19, "satoshi", 1, 1, true, addr, sdk.Coins{}), false},

		{"gateway basic good", NewMsgIssueToken(FUNGIBLE, GATEWAY, "abc", "btc", "btc", "btc", 18, "satoshi", 1, 1, true, addr, sdk.Coins{}), true},
		{"gateway symbol_at_source error", NewMsgIssueToken(FUNGIBLE, GATEWAY, "a", "btc", "a1,d", "btc", 18, "satoshi", 1, 1, true, addr, sdk.Coins{}), false},
		{"gateway symbol_at_source too long", NewMsgIssueToken(FUNGIBLE, GATEWAY, "a", "btc", "abcdefghijklmn", "btc", 18, "satoshi", 1, 1, true, addr, sdk.Coins{}), false},
		{"gateway symbol_at_source too short", NewMsgIssueToken(FUNGIBLE, GATEWAY, "a", "btc", "a", "btc", 18, "satoshi", 1, 1, true, addr, sdk.Coins{}), false},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgIssueToken.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgIssueToken.ValidateBasic(), "test: %v", tc.testCase)
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

	expected := "{\"type\":\"irishub/asset/MsgCreateGateway\",\"value\":{\"details\":\"details\",\"identity\":\"identity\",\"moniker\":\"moniker\",\"owner\":\"faa1damkuetjqqah8w\",\"website\":\"website\"}}"
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

// test ValidateBasic for MsgIssueToken
func TestMsgEditAsset(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	mintable := false
	tests := []struct {
		testCase string
		MsgEditToken
		expectPass bool
	}{
		{"native basic good", NewMsgEditToken("BTC Token", "btc", "satoshi", "x.btc", 10000, &mintable, owner), true},
		{"wrong symbol_at_source", NewMsgEditToken("BTC Token", "HT", "satoshi", "x.btc", 10000, &mintable, owner), false},
		{"wrong symbol_min_alias", NewMsgEditToken("BTC Token", "btc", "btc-min", "x.ht", 10000, &mintable, owner), false},
		{"wrong token_id", NewMsgEditToken("BTC Token", "HTC", "HT", "i.ht", 10000, &mintable, owner), false},
		{"wrong max_supply", NewMsgEditToken("BTC Token", "btc", "satoshi", "x.btc", 10000000000000, &mintable, owner), false},
		{"loss owner", NewMsgEditToken("BTC Token", "btc", "satoshi", "x.btc", 10000, &mintable, nil), false},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgEditToken.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgEditToken.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}

func TestMsgEditTokenRoute(t *testing.T) {
	symbolAtSource := "btc"
	symbolMinAlias := "satoshi"
	tokenId := "x.btc"
	mintable := false
	// build a MsgEditToken
	msg := MsgEditToken{
		SymbolAtSource: symbolAtSource,
		SymbolMinAlias: symbolMinAlias,
		MaxSupply:      10000000,
		Mintable:       &mintable,
		TokenId:        tokenId,
	}

	require.Equal(t, "asset", msg.Route())
}

func TestMsgEditTokenGetSignBytes(t *testing.T) {
	mintable := false
	var msg = MsgEditToken{
		Name:           "BTC TOKEN",
		Owner:          sdk.AccAddress([]byte("owner")),
		TokenId:        "x.btc",
		SymbolAtSource: "btc",
		SymbolMinAlias: "satoshi",
		MaxSupply:      21000000,
		Mintable:       &mintable,
	}

	res := msg.GetSignBytes()

	expected := `{"type":"irishub/asset/MsgEditToken","value":{"max_supply":"21000000","mintable":false,"name":"BTC TOKEN","owner":"faa1damkuetjqqah8w","symbol_at_source":"btc","symbol_min_alias":"satoshi","token_id":"x.btc"}}`
	require.Equal(t, expected, string(res))
}

func TestNewMsgTransferGatewayOwner(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	moniker := "mon"
	newOwner := sdk.AccAddress([]byte("newOwner"))

	msg := NewMsgTransferGatewayOwner(owner, moniker, newOwner)

	require.Equal(t, owner, msg.Owner)
	require.Equal(t, moniker, msg.Moniker)
	require.Equal(t, newOwner, msg.To)
}

func TestMsgTransferGatewayOwnerRoute(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	moniker := "mon"
	newOwner := sdk.AccAddress([]byte("newOwner"))

	msg := NewMsgTransferGatewayOwner(owner, moniker, newOwner)
	require.Equal(t, "asset", msg.Route())
}

func TestMsgTransferGatewayOwnerValidation(t *testing.T) {
	testData := []struct {
		name       string
		owner      sdk.AccAddress
		moniker    string
		to         sdk.AccAddress
		expectPass bool
	}{
		{"empty owner", emptyAddr, "mon", addr1, false},
		{"empty moniker", addr1, "", addr2, false},
		{"too short moniker", addr1, "mo", addr2, false},
		{"too long moniker", addr1, "monikermo", addr2, false},
		{"moniker contains illegal characters", addr1, "moni2", addr2, false},
		{"empty to address", addr1, "mon", emptyAddr, false},
		{"the to address is same as the owner", addr1, "mon", addr1, false},
		{"basic good", addr1, "mon", addr2, true},
	}

	for _, td := range testData {
		msg := NewMsgTransferGatewayOwner(td.owner, td.moniker, td.to)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}

func TestMsgTransferGatewayOwnerGetSignBytes(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	moniker := "mon"
	newOwner := sdk.AccAddress([]byte("newOwner"))

	msg := NewMsgTransferGatewayOwner(owner, moniker, newOwner)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/asset/MsgTransferGatewayOwner","value":{"moniker":"mon","owner":"faa1damkuetjqqah8w","to":"faa1dejhwnmhdejhyxslsx8"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgTransferGatewayOwnerGetSigners(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	moniker := "mon"
	newOwner := sdk.AccAddress([]byte("newOwner"))

	msg := NewMsgTransferGatewayOwner(owner, moniker, newOwner)

	res := msg.GetSigners()

	expected := "[6F776E6572 6E65774F776E6572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

func TestMsgMintTokenValidateBasic(t *testing.T) {
	testData := []struct {
		msg        string
		tokeId     string
		owner      sdk.AccAddress
		to         sdk.AccAddress
		amount     uint64
		expectPass bool
	}{
		{"empty tokeId", "", addr1, addr2, 1000, false},
		{"wrong tokeId", "p.btc", addr1, addr2, 1000, false},
		{"empty owner", "btc", emptyAddr, addr2, 1000, false},
		{"empty to", "btc", addr1, emptyAddr, 1000, true},
		{"not empty to", "btc", addr1, addr2, 1000, true},
		{"invalid amount", "btc", addr1, addr2, 0, false},
		{"exceed max supply", "btc", addr1, addr2, 100000000000000, false},
		{"basic good", "btc", addr1, addr2, 1000, true},
	}

	for _, td := range testData {
		msg := NewMsgMintToken(td.tokeId, td.owner, td.to, td.amount)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.msg)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.msg)
		}
	}
}
