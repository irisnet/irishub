package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		{"basic good", NewMsgIssueToken("btc", "btc", 18, "satoshi", 1, 1, true, addr), true},
		{"symbol empty", NewMsgIssueToken("", "btc", 18, "satoshi", 1, 1, true, addr), false},
		{"symbol error", NewMsgIssueToken("ab,c", "btc", 18, "satoshi", 1, 1, true, addr), false},
		{"symbol first letter is num", NewMsgIssueToken("4iris", "btc", 18, "satoshi", 1, 1, true, addr), false},
		{"symbol too long", NewMsgIssueToken("aaaaaaaaa", "btc", 18, "satoshi", 1, 1, true, addr), false},
		{"symbol too short", NewMsgIssueToken("a", "btc", 18, "satoshi", 1, 1, true, addr), false},
		{"min_unit error", NewMsgIssueToken("btc", "btc", 18, "a1,3d", 1, 1, true, addr), false},
		{"min_unit too long", NewMsgIssueToken("btc", "btc", 18, "aaaaaaaaaaaaa", 1, 1, true, addr), false},
		{"min_unit too short", NewMsgIssueToken("btc", "btc", 18, "a", 1, 1, true, addr), false},
		{"min_unit  first letter is num", NewMsgIssueToken("btc", "btc", 18, "1a", 1, 1, true, addr), false},
		{"name empty", NewMsgIssueToken("btc", "", 18, "satoshi", 1, 1, true, addr), false},
		{"name blank", NewMsgIssueToken("btc", " ", 18, "satoshi", 1, 1, true, addr), false},
		{"name too long", NewMsgIssueToken("btc", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 1, "satoshi", 1, 1, true, addr), false},
		{"initial supply is zero", NewMsgIssueToken("btc", "btc", 18, "satoshi", 0, 1, true, addr), true},
		{"max supply is zero", NewMsgIssueToken("btc", "btc", 18, "satoshi", 1, 0, true, addr), true},
		{"init supply bigger than max supply", NewMsgIssueToken("btc", "btc", 18, "satoshi", 2, 1, true, addr), false},
		{"decimal error", NewMsgIssueToken("btc", "btc", 19, "satoshi", 1, 1, true, addr), false},
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
		{"basic good", NewMsgEditToken("BTC Token", "btc", 10000, mintable, owner), true},
		{"wrong name", NewMsgEditToken("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "HT", 10000, mintable, owner), false},
		{"wrong symbol", NewMsgEditToken("BTC Token", "HT", 10000, mintable, owner), false},
		{"wrong max_supply", NewMsgEditToken("BTC Token", "btc", 10000000000000, mintable, owner), false},
		{"loss owner", NewMsgEditToken("BTC Token", "btc", 10000, mintable, nil), false},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgEditToken.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgEditToken.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}

func TestMsgEditTokenGetSignBytes(t *testing.T) {
	mintable := False
	var msg = MsgEditToken{
		Symbol:    "btc",
		Name:      "BTC TOKEN",
		Owner:     sdk.AccAddress([]byte("owner")),
		MaxSupply: 21000000,
		Mintable:  mintable,
	}

	res := msg.GetSignBytes()

	expected := `{"type":"irishub/asset/token/MsgEditToken","value":{"max_supply":"21000000","mintable":"false","name":"BTC TOKEN","owner":"cosmos1damkuetjzyud4a","symbol":"btc"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgMintTokenValidateBasic(t *testing.T) {
	testData := []struct {
		msg        string
		symbol     string
		owner      sdk.AccAddress
		to         sdk.AccAddress
		amount     uint64
		expectPass bool
	}{
		{"empty symbol", "", addr1, addr2, 1000, false},
		{"wrong symbol", "p.btc", addr1, addr2, 1000, false},
		{"empty owner", "btc", emptyAddr, addr2, 1000, false},
		{"empty to", "btc", addr1, emptyAddr, 1000, true},
		{"not empty to", "btc", addr1, addr2, 1000, true},
		{"invalid amount", "btc", addr1, addr2, 0, false},
		{"exceed max supply", "btc", addr1, addr2, 100000000000000, false},
		{"basic good", "btc", addr1, addr2, 1000, true},
	}

	for _, td := range testData {
		msg := NewMsgMintToken(td.symbol, td.owner, td.to, td.amount)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.msg)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.msg)
		}
	}
}

func TestMsgTransferTokenValidation(t *testing.T) {
	testData := []struct {
		name       string
		srcOwner   sdk.AccAddress
		symbol     string
		dstOwner   sdk.AccAddress
		expectPass bool
	}{
		{"empty srcOwner", emptyAddr, "btc", addr1, false},
		{"empty symbol", addr1, "", addr2, false},
		{"empty dstOwner", addr1, "btc", emptyAddr, false},
		{"invalid symbol", addr1, "btc-min", addr2, false},
		{"basic good", addr1, "btc", addr2, true},
	}

	for _, td := range testData {
		msg := NewMsgTransferToken(td.srcOwner, td.dstOwner, td.symbol)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}

func TestMsgBurnToken(t *testing.T) {
	testData := []struct {
		name       string
		sender     sdk.AccAddress
		amount     sdk.Coins
		expectPass bool
	}{
		{"empty sender", emptyAddr, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000))), false},
		{"empty amount", addr1, sdk.Coins{}, false},
		{"basic good", addr1, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000))), true},
	}

	for _, td := range testData {
		msg := NewMsgBurnToken(td.sender, td.amount)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}
