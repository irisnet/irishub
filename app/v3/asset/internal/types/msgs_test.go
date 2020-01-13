package types

import (
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
		{"basic good", NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 1, 1, true, addr), true},
		{"symbol empty", NewMsgIssueToken("", "satoshi", "Bitcoin Network", 18, 1, 1, true, addr), false},
		{"symbol error", NewMsgIssueToken("b&tc", "satoshi", "Bitcoin Network", 18, 1, 1, true, addr), false},
		{"symbol first letter is num", NewMsgIssueToken("4btc", "satoshi", "Bitcoin Network", 18, 1, 1, true, addr), false},
		{"symbol too long", NewMsgIssueToken("btc1111111111", "satoshi", "Bitcoin Network", 18, 1, 1, true, addr), false},
		{"symbol too short", NewMsgIssueToken("ht", "satoshi", "Bitcoin Network", 18, 1, 1, true, addr), false},
		{"name empty", NewMsgIssueToken("btc", "satoshi", "", 18, 1, 1, true, addr), false},
		{"name blank", NewMsgIssueToken("btc", "satoshi", " ", 18, 1, 1, true, addr), false},
		{"name too long", NewMsgIssueToken("btc", "satoshi", "Bitcoin Network aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 18, 1, 1, true, addr), false},
		{"initial supply is zero", NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 0, 1, true, addr), true},
		{"max supply is zero", NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 1, 0, true, addr), true},
		{"init supply bigger than max supply", NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 2, 1, true, addr), false},
		{"decimal error", NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 30, 1, 1, true, addr), false},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgIssueToken.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgIssueToken.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}

// test ValidateBasic for MsgIssueToken
func TestMsgEditToken(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	mintable := False
	tests := []struct {
		testCase string
		MsgEditToken
		expectPass bool
	}{
		{"native basic good", NewMsgEditToken("BTC Token", "i.btc", 10000, mintable, owner), true},
		{"wrong token_id", NewMsgEditToken("BTC Token", "HTC", 10000, mintable, owner), false},
		{"wrong max_supply", NewMsgEditToken("BTC Token", "i.btc", 10000000000000, mintable, owner), false},
		{"loss owner", NewMsgEditToken("BTC Token", "i.btc", 10000, mintable, nil), false},
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
	canonicalSymbol := "btc"
	minUnitAlias := "satoshi"
	tokenId := "x.btc"
	mintable := False
	// build a MsgEditToken
	msg := MsgEditToken{
		CanonicalSymbol: canonicalSymbol,
		MinUnitAlias:    minUnitAlias,
		MaxSupply:       10000000,
		Mintable:        mintable,
		TokenId:         tokenId,
	}

	require.Equal(t, "asset", msg.Route())
}

func TestMsgEditTokenGetSignBytes(t *testing.T) {
	mintable := False
	var msg = MsgEditToken{
		Name:            "BTC TOKEN",
		Owner:           sdk.AccAddress([]byte("owner")),
		TokenId:         "x.btc",
		CanonicalSymbol: "btc",
		MinUnitAlias:    "satoshi",
		MaxSupply:       21000000,
		Mintable:        mintable,
	}

	res := msg.GetSignBytes()

	expected := `{"type":"irishub/asset/MsgEditToken","value":{"canonical_symbol":"btc","max_supply":"21000000","min_unit_alias":"satoshi","mintable":"false","name":"BTC TOKEN","owner":"faa1damkuetjqqah8w","token_id":"x.btc"}}`
	require.Equal(t, expected, string(res))
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
		{"empty owner", "i.btc", emptyAddr, addr2, 1000, false},
		{"empty to", "i.btc", addr1, emptyAddr, 1000, true},
		{"not empty to", "i.btc", addr1, addr2, 1000, true},
		{"invalid amount", "i.btc", addr1, addr2, 0, false},
		{"exceed max supply", "i.btc", addr1, addr2, 100000000000000, false},
		{"basic good", "i.btc", addr1, addr2, 1000, true},
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

func TestMsgTransferTokenOwnerValidation(t *testing.T) {
	testData := []struct {
		name       string
		srcOwner   sdk.AccAddress
		tokenId    string
		dstOwner   sdk.AccAddress
		expectPass bool
	}{
		{"empty srcOwner", emptyAddr, "btc", addr1, false},
		{"empty tokenId", addr1, "", addr2, false},
		{"empty dstOwner", addr1, "i.btc", emptyAddr, false},
		{"invalid tokenId", addr1, "btc-min", addr2, false},
		{"basic good", addr1, "i.btc", addr2, true},
	}

	for _, td := range testData {
		msg := NewMsgTransferTokenOwner(td.srcOwner, td.dstOwner, td.tokenId)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}
